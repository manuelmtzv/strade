package main

import (
	"strade/internal/config"
	"strade/internal/env"
)

func getConfig() config.IngestorConfig {
	return config.IngestorConfig{
		SourceURL:            env.GetString("BROWSER_INGESTOR_SOURCE_URL", ""),
		SettlementsBatchSize: env.GetInt("INGESTOR_SETTLEMENTS_BATCH_SIZE", 10000),
		SettlementsWorkers:   env.GetInt("INGESTOR_SETTLEMENTS_WORKERS", 4),
		DB: config.DBConfig{
			Addr:         env.GetString("DB_ADDR", ""),
			MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 25),
			MaxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}
}
