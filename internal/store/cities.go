package store

import (
	"database/sql"
	"strade/internal/models"
)

type CityStore struct {
	db *sql.DB
}

func NewCityStore(db *sql.DB) *CityStore {
	return &CityStore{db: db}
}

func (s *CityStore) Create(city *models.City) error {
	return nil
}

func (s *CityStore) CreateTx(tx *sql.Tx, city *models.City) error {
	return nil
}

func (s *CityStore) FindAll() ([]*models.City, error) {
	return nil, nil
}
