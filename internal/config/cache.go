package config

type CacheConfig struct {
	Redis RedisConfig
}

type RedisConfig struct {
	Addr    string
	PW      string
	DB      int
	Enabled bool
}
