package store

import (
	"context"
	"database/sql"
	"fmt"
	"strade/internal/models"
	"strings"
)

type MunicipalityStore struct {
	db *sql.DB
}

func NewMunicipalityStore(db *sql.DB) *MunicipalityStore {
	return &MunicipalityStore{db: db}
}

var (
	upsertMunicipalitiesQuery = `
		INSERT INTO municipalities (id, name, slug, state_id)
		VALUES %s
		ON CONFLICT (id) DO UPDATE
		SET name = EXCLUDED.name, slug = EXCLUDED.slug, state_id = EXCLUDED.state_id
	`
)

func (s *MunicipalityStore) BulkUpsertTx(ctx context.Context, tx *sql.Tx, municipalities []*models.Municipality) error {
	if len(municipalities) == 0 {
		return nil
	}

	valueStrings := make([]string, 0, len(municipalities))
	valueArgs := make([]any, 0, len(municipalities)*4)
	i := 1

	for _, m := range municipalities {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d)", i, i+1, i+2, i+3))
		valueArgs = append(valueArgs, m.ID, m.Name, m.Slug, m.StateID)
		i += 4
	}

	stmt := fmt.Sprintf(upsertMunicipalitiesQuery, strings.Join(valueStrings, ","))
	_, err := tx.ExecContext(ctx, stmt, valueArgs...)
	return err
}

func (s *MunicipalityStore) FindAll() ([]*models.Municipality, error) {
	return nil, nil
}
