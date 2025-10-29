package config

import "time"

type APIConfig struct {
	Addr            string
	ShutdownTimeout time.Duration
	DB              DBConfig
	Cache           CacheConfig
}
