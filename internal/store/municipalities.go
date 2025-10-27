package store

import (
	"database/sql"
	"strade/internal/models"
)

type MunicipalityStore struct {
	db *sql.DB
}

func NewMunicipalityStore(db *sql.DB) *MunicipalityStore {
	return &MunicipalityStore{db: db}
}

func (s *MunicipalityStore) Create(municipality *models.Municipality) error {
	return nil
}

func (s *MunicipalityStore) CreateTx(tx *sql.Tx, municipality *models.Municipality) error {
	return nil
}

func (s *MunicipalityStore) FindAll() ([]*models.Municipality, error) {
	return nil, nil
}
