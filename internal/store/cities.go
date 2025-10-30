package store

import (
	"context"
	"database/sql"
	"fmt"
	"strade/internal/models"
	"strings"
)

type CityStore struct {
	db *sql.DB
}

func NewCityStore(db *sql.DB) *CityStore {
	return &CityStore{db: db}
}

var (
	upsertCitiesQuery = `
		INSERT INTO cities (id, name, slug, state_id)
		VALUES %s
		ON CONFLICT (id) DO UPDATE
		SET name = EXCLUDED.name, slug = EXCLUDED.slug, state_id = EXCLUDED.state_id
	`
)

func (s *CityStore) BulkUpsertTx(ctx context.Context, tx *sql.Tx, cities []*models.City) error {
	if len(cities) == 0 {
		return nil
	}

	valueStrings := make([]string, 0, len(cities))
	valueArgs := make([]any, 0, len(cities)*4)
	i := 1

	for _, c := range cities {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d)", i, i+1, i+2, i+3))
		valueArgs = append(valueArgs, c.ID, c.Name, c.Slug, c.StateID)
		i += 4
	}

	stmt := fmt.Sprintf(upsertCitiesQuery, strings.Join(valueStrings, ","))
	_, err := tx.ExecContext(ctx, stmt, valueArgs...)
	return err
}

func (s *CityStore) FindAll() ([]*models.City, error) {
	return nil, nil
}
