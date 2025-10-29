package store

import "database/sql"

type Storage struct {
	DB                  *sql.DB
	StateStore          StateStorage
	MunicipalityStore   MunicipalityStorage
	CityStore           CityStorage
	SettlementTypeStore SettlementTypeStorage
	SettlementStore     SettlementStorage
	WatermarkStore      *WatermarkStore
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		DB:                  db,
		StateStore:          NewStateStore(db),
		MunicipalityStore:   NewMunicipalityStore(db),
		CityStore:           NewCityStore(db),
		SettlementTypeStore: NewSettlementTypeStore(db),
		SettlementStore:     NewSettlementStore(db),
		WatermarkStore:      NewWatermarkStore(db),
	}
}
