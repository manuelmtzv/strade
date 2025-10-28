package config

type IngestorConfig struct {
	SourceURL            string
	SettlementsBatchSize int
	SettlementsWorkers   int
	DB                   DBConfig
}
