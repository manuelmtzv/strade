package watch

import (
	"context"
	"math/rand"
	"strade/internal/config"
	"strade/internal/ingest"
	"strade/internal/store"
	"strade/internal/utils"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"go.uber.org/zap"
)

type BrowserWatcher struct {
	logger    *zap.SugaredLogger
	config    config.WatcherConfig
	ingestor  ingest.Ingestor
	watermark *store.WatermarkStore
	locker    *Locker
}

func NewBrowserWatcher(
	logger *zap.SugaredLogger,
	config config.WatcherConfig,
	ingestor ingest.Ingestor,
	watermark *store.WatermarkStore,
	locker *Locker,
) *BrowserWatcher {
	return &BrowserWatcher{
		logger:    logger,
		config:    config,
		ingestor:  ingestor,
		watermark: watermark,
		locker:    locker,
	}
}

func (w *BrowserWatcher) Run(ctx context.Context) error {
	ticker := time.NewTicker(w.config.Interval)
	defer ticker.Stop()

	w.logger.Info("Browser watcher started")

	if err := w.tick(ctx); err != nil {
		w.logger.Warnw("initial tick failed", "err", err)
	}

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("Browser watcher stopped")
			return ctx.Err()
		case <-ticker.C:
			if w.config.Jitter > 0 {
				jitterDuration := time.Duration(rand.Int63n(int64(w.config.Jitter)))
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(jitterDuration):
				}
			}

			if err := w.tick(ctx); err != nil {
				w.logger.Warnw("watcher tick failed", "err", err)
			}
		}
	}
}

func (w *BrowserWatcher) tick(ctx context.Context) error {
	acquired, release, err := w.locker.TryLock(ctx, w.config.LockKey, 30*time.Second)
	if err != nil {
		return err
	}
	if !acquired {
		w.logger.Debug("lock not acquired, another instance is running")
		return nil
	}
	defer release()

	w.logger.Info("lock acquired, starting ingestion check")

	lastWatermark, err := w.watermark.Get(ctx, w.config.WMKey)
	if err != nil {
		w.logger.Warnw("failed to get watermark", "err", err)
	}

	var lastWatermarkTime time.Time
	if lastWatermark == "" {
		w.logger.Info("no watermark found, will process all available data")
		lastWatermarkTime = time.Time{}
	} else {
		w.logger.Infow("current watermark", "value", lastWatermark)
		lastWatermarkTime, err = time.Parse("2006-01-02", lastWatermark)
		if err != nil {
			w.logger.Warnw("failed to parse watermark, using zero time", "err", err)
			lastWatermarkTime = time.Time{}
		}
	}

	hasNewData, err := w.check(ctx, lastWatermarkTime)
	if err != nil {
		return err
	}

	if !hasNewData {
		w.logger.Info("no new data to ingest")
		return nil
	}

	w.logger.Info("new data detected, running ingestor")
	if err := w.ingestor.Ingest(ctx); err != nil {
		return err
	}

	newWatermark := time.Now().Format("2006-01-02")
	if err := w.watermark.Set(ctx, w.config.WMKey, newWatermark); err != nil {
		w.logger.Warnw("failed to update watermark", "err", err)
	}

	w.logger.Infow("ingestion completed successfully", "newWatermark", newWatermark)
	return nil
}

func (w *BrowserWatcher) check(ctx context.Context, lastWatermark time.Time) (bool, error) {
	browser := rod.New().Context(ctx).MustConnect()
	defer browser.MustClose()

	page := browser.MustPage(w.config.SourceURL)
	page.MustWaitLoad()

	rawText := strings.TrimSpace(page.MustElement("#lblfec").MustText())
	if rawText == "" {
		return false, nil
	}

	rawDate := strings.TrimSpace(strings.Split(rawText, ":")[1])

	lastDate, err := utils.ParseSpanishDate(rawDate)
	if err != nil {
		w.logger.Warnw("failed to parse date", "rawDate", rawDate, "err", err)
		return false, err
	}

	if lastDate.After(lastWatermark) {
		return true, nil
	}

	return false, nil
}
