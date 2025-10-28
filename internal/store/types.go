package store

import (
	"context"
	"database/sql"
	"strade/internal/models"
)

type StateStorage interface {
	BulkUpsertTx(ctx context.Context, tx *sql.Tx, states []*models.State) error
	FindAll() ([]*models.State, error)
}

type MunicipalityStorage interface {
	BulkUpsertTx(ctx context.Context, tx *sql.Tx, municipalities []*models.Municipality) error
	FindAll() ([]*models.Municipality, error)
}

type CityStorage interface {
	BulkUpsertTx(ctx context.Context, tx *sql.Tx, cities []*models.City) error
	FindAll() ([]*models.City, error)
}

type SettlementTypeStorage interface {
	BulkUpsertTx(ctx context.Context, tx *sql.Tx, settlementTypes []*models.SettlementType) error
	FindAll() ([]*models.SettlementType, error)
}

type SettlementStorage interface {
	BulkUpsertTx(ctx context.Context, tx *sql.Tx, settlements []*models.Settlement) error
	FindAll() ([]*models.Settlement, error)
}
