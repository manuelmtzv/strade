package ingest

import (
	"context"
	"strade/internal/config"
	"strade/internal/store"

	"go.uber.org/zap"
)

type Ingestor interface {
	Ingest(ctx context.Context) error
}

func NewIngestor(config config.IngestorConfig, logger *zap.SugaredLogger, store *store.Storage, ingestorType string) Ingestor {
	switch ingestorType {
	case "browser":
		return NewBrowserIngestor(config, logger, store)
	default:
		return nil
	}
}
