package main

import (
	"context"
	"database/sql"
	"strade/internal/db"
	"strade/internal/env"
	"strade/internal/ingest"
	"strade/internal/store"

	"go.uber.org/zap"
)

func main() {
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer func() {
		_ = logger.Sync()
	}()

	if err := env.Load(); err != nil {
		logger.Panicf("Failed to load env: %v", err)
	}

	cfg := getConfig()

	dbConn, err := db.New(
		cfg.DB.Addr,
		cfg.DB.MaxOpenConns,
		cfg.DB.MaxIdleConns,
		cfg.DB.MaxIdleTime,
	)
	if err != nil {
		logger.Panicf("Database connection failed: %v", err)
	}
	defer func(dbConn *sql.DB) {
		_ = dbConn.Close()
	}(dbConn)
	logger.Info("Database connection established")

	storage := store.NewStorage(dbConn)

	browserIngestor := ingest.NewBrowserIngestor(logger, storage, cfg.SourceURL)

	ctx := context.Background()

	if err := browserIngestor.Ingest(ctx); err != nil {
		logger.Panicf("Failed to ingest: %v", err)
	}
}
