package watch

import (
	"context"
	"strade/internal/config"
	"strade/internal/ingest"
	"strade/internal/store"

	"go.uber.org/zap"
)

type Watcher interface {
	Run(ctx context.Context) error
}

func NewWatcher(
	cfg config.WatcherConfig,
	logger *zap.SugaredLogger,
	storage *store.Storage,
	ingestor ingest.Ingestor,
	watcherType string,
) Watcher {
	switch watcherType {
	case "browser":
		locker := NewLocker(storage.DB)
		return NewBrowserWatcher(logger, cfg, ingestor, storage.WatermarkStore, locker)
	default:
		return nil
	}
}
