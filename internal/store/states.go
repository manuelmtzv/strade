package store

import (
	"context"
	"database/sql"
	"fmt"
	"strade/internal/models"
	"strings"
)

type StateStore struct {
	db *sql.DB
}

func NewStateStore(db *sql.DB) *StateStore {
	return &StateStore{db: db}
}

var (
	upsertStatesQuery = `
		INSERT INTO states (id, name, slug)
		VALUES %s
		ON CONFLICT (id) DO UPDATE
		SET name = EXCLUDED.name, slug = EXCLUDED.slug
	`
)

func (s *StateStore) BulkUpsertTx(ctx context.Context, tx *sql.Tx, states []*models.State) error {
	if len(states) == 0 {
		return nil
	}

	valueStrings := make([]string, 0, len(states))
	valueArgs := make([]any, 0, len(states)*3)
	i := 1

	for _, st := range states {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d)", i, i+1, i+2))
		valueArgs = append(valueArgs, st.ID, st.Name, st.Slug)
		i += 3
	}

	stmt := fmt.Sprintf(upsertStatesQuery, strings.Join(valueStrings, ","))
	_, err := tx.ExecContext(ctx, stmt, valueArgs...)
	return err
}

func (s *StateStore) FindAll() ([]*models.State, error) {
	return nil, nil
}
