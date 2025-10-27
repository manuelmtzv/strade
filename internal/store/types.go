package store

import (
	"database/sql"
	"strade/internal/models"
)

type StateStorage interface {
	Create(state *models.State) error
	CreateTx(tx *sql.Tx, state *models.State) error
	FindAll() ([]*models.State, error)
}

type MunicipalityStorage interface {
	Create(municipality *models.Municipality) error
	CreateTx(tx *sql.Tx, municipality *models.Municipality) error
	FindAll() ([]*models.Municipality, error)
}

type CityStorage interface {
	Create(city *models.City) error
	CreateTx(tx *sql.Tx, city *models.City) error
	FindAll() ([]*models.City, error)
}

type SettlementTypeStorage interface {
	Create(settlementType *models.SettlementType) error
	CreateTx(tx *sql.Tx, settlementType *models.SettlementType) error
	FindAll() ([]*models.SettlementType, error)
}

type SettlementStorage interface {
	Create(settlement *models.Settlement) error
	CreateTx(tx *sql.Tx, settlement *models.Settlement) error
	FindAll() ([]*models.Settlement, error)
}
