package store

import (
	"database/sql"
	"strade/internal/models"
)

type SettlementStore struct {
	db *sql.DB
}

func NewSettlementStore(db *sql.DB) *SettlementStore {
	return &SettlementStore{db: db}
}

func (s *SettlementStore) Create(settlement *models.Settlement) error {
	return nil
}

func (s *SettlementStore) CreateTx(tx *sql.Tx, settlement *models.Settlement) error {
	return nil
}

func (s *SettlementStore) FindAll() ([]*models.Settlement, error) {
	return nil, nil
}
