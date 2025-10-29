package main

import (
	"strade/internal/config"
	"strade/internal/env"
	"time"

	"go.uber.org/zap"
)

type Config struct {
	Ingestor config.IngestorConfig
	Watcher  config.WatcherConfig
}

func getConfig(logger *zap.SugaredLogger) Config {
	interval := env.GetDuration("WATCHER_INTERVAL", 6*time.Hour)

	if interval < 6*time.Hour {
		logger.Warn("Watcher interval is less than 6 hours, setting to 6 hours")
		interval = 6 * time.Hour
	}

	return Config{
		Ingestor: config.IngestorConfig{
			SourceURL:            env.GetString("BROWSER_INGESTOR_SOURCE_URL", ""),
			SettlementsBatchSize: env.GetInt("INGESTOR_SETTLEMENTS_BATCH_SIZE", 10000),
			SettlementsWorkers:   env.GetInt("INGESTOR_SETTLEMENTS_WORKERS", 4),
			DB: config.DBConfig{
				Addr:         env.GetString("DB_ADDR", ""),
				MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 25),
				MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 25),
				MaxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
			},
		},
		Watcher: config.WatcherConfig{
			SourceURL: env.GetString("BROWSER_WATCHER_SOURCE_URL", ""),
			Interval:  interval,
			Jitter:    env.GetDuration("WATCHER_JITTER", 30*time.Second),
			LockKey:   env.GetString("WATCHER_LOCK_KEY", "browser:watcher:lock"),
			WMKey:     env.GetString("WATCHER_WM_KEY", "browser:watcher:watermark"),
		},
	}
}
