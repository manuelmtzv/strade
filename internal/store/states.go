package store

import (
	"database/sql"
	"strade/internal/models"
)

type StateStore struct {
	db *sql.DB
}

func NewStateStore(db *sql.DB) *StateStore {
	return &StateStore{db: db}
}

func (s *StateStore) Create(state *models.State) error {
	return nil
}

func (s *StateStore) CreateTx(tx *sql.Tx, state *models.State) error {
	return nil
}

func (s *StateStore) FindAll() ([]*models.State, error) {
	return nil, nil
}
