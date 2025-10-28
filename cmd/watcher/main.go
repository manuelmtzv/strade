package main

import (
	"context"
	"database/sql"
	"os"
	"os/signal"
	"strade/internal/db"
	"strade/internal/env"
	"strade/internal/ingest"
	"strade/internal/store"
	"strade/internal/watch"
	"syscall"

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
		cfg.Ingestor.DB.Addr,
		cfg.Ingestor.DB.MaxOpenConns,
		cfg.Ingestor.DB.MaxIdleConns,
		cfg.Ingestor.DB.MaxIdleTime,
	)
	if err != nil {
		logger.Panicf("Database connection failed: %v", err)
	}
	defer func(dbConn *sql.DB) {
		_ = dbConn.Close()
	}(dbConn)
	logger.Info("Database connection established")

	storage := store.NewStorage(dbConn)

	browserIngestor := ingest.NewIngestor(cfg.Ingestor, logger, storage, "browser")

	watcher := watch.NewWatcher(cfg.Watcher, logger, storage, browserIngestor, "browser")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Info("Shutdown signal received")
		cancel()
	}()

	logger.Info("Starting watcher...")
	if err := watcher.Run(ctx); err != nil && err != context.Canceled {
		logger.Panicf("Watcher failed: %v", err)
	}

	logger.Info("Watcher stopped gracefully")
}
