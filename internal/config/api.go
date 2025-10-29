package config

type APIConfig struct {
	Addr  string
	DB    DBConfig
	Cache CacheConfig
}
