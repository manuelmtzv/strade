package store

import (
	"database/sql"
	"strade/internal/models"
)

type SettlementTypeStore struct {
	db *sql.DB
}

func NewSettlementTypeStore(db *sql.DB) *SettlementTypeStore {
	return &SettlementTypeStore{db: db}
}

func (s *SettlementTypeStore) Create(settlementType *models.SettlementType) error {
	return nil
}

func (s *SettlementTypeStore) CreateTx(tx *sql.Tx, settlementType *models.SettlementType) error {
	return nil
}

func (s *SettlementTypeStore) FindAll() ([]*models.SettlementType, error) {
	return nil, nil
}
