package ingest

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strade/internal/store"
	"strade/internal/utils"
	"time"

	"github.com/go-rod/rod"
	"go.uber.org/zap"
)

type BrowserIngestor struct {
	logger    *zap.SugaredLogger
	store     store.Storage
	sourceURL string
}

func NewBrowserIngestor(logger *zap.SugaredLogger, store store.Storage, sourceURL string) *BrowserIngestor {
	return &BrowserIngestor{
		logger:    logger,
		store:     store,
		sourceURL: sourceURL,
	}
}

func (s *BrowserIngestor) Ingest(ctx context.Context) error {
	browser := rod.New().MustConnect()
	defer browser.MustClose()

	page := browser.MustPage(s.sourceURL)
	page.MustWaitLoad()

	page.MustElement("#rblTipo_1").MustClick()
	page.MustElement("#btnDescarga").MustClick()

	wait := page.Browser().MustWaitDownload()
	data := wait()

	s.logger.Infof("Downloaded %d bytes (ZIP)", len(data))

	tmpDir := filepath.Join(".", "tmp")
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return fmt.Errorf("failed to create tmp dir: %w", err)
	}

	zipPath := filepath.Join(tmpDir, fmt.Sprintf("sepomex_%d.zip", time.Now().Unix()))
	if err := os.WriteFile(zipPath, data, 0644); err != nil {
		return fmt.Errorf("failed to save zip file: %w", err)
	}

	s.logger.Infof("Saved zip to %s", zipPath)

	if err := utils.UnzipFile(zipPath, tmpDir); err != nil {
		return fmt.Errorf("failed to unzip file: %w", err)
	}

	s.logger.Info("Ingestion completed successfully.")
	return nil
}
