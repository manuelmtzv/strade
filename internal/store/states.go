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

const (
	findAllStatesQuery = `
		SELECT id, name, slug 
		FROM states
		ORDER BY name ASC
	`
	findMunicipalitiesByStateIDQuery = `
		SELECT id, name, slug, state_id
		FROM municipalities 
		WHERE state_id = $1
		ORDER BY name ASC
	`
)

func (s *StateStore) FindAll(ctx context.Context) ([]*models.State, error) {
	rows, err := s.db.QueryContext(ctx, findAllStatesQuery)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var states []*models.State
	for rows.Next() {
		var state models.State
		if err := rows.Scan(
			&state.ID,
			&state.Name,
			&state.Slug,
		); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		states = append(states, &state)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return states, nil
}

func (s *StateStore) FindMunicipalitiesByStateID(ctx context.Context, stateID string) ([]*models.Municipality, error) {
	rows, err := s.db.QueryContext(ctx, findMunicipalitiesByStateIDQuery, stateID)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var municipalities []*models.Municipality
	for rows.Next() {
		var m models.Municipality
		if err := rows.Scan(
			&m.ID,
			&m.Name,
			&m.Slug,
			&m.StateID,
		); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		municipalities = append(municipalities, &m)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return municipalities, nil
}
