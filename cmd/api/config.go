package main

import (
	"strade/internal/config"
	"strade/internal/env"
)

func getConfig() config.APIConfig {
	return config.APIConfig{
		Addr: env.GetString("ADDR", ":8080"),
		DB: config.DBConfig{
			Addr:         env.GetString("DB_ADDR", ""),
			MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 25),
			MaxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		Cache: config.CacheConfig{
			Redis: config.RedisConfig{
				Addr:    env.GetString("REDIS_ADDR", ""),
				PW:      env.GetString("REDIS_PW", ""),
				DB:      env.GetInt("REDIS_DB", 0),
				Enabled: env.GetBool("REDIS_ENABLED", false),
			},
		},
	}
}
