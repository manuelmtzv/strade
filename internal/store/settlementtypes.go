package store

import (
	"context"
	"database/sql"
	"fmt"
	"strade/internal/models"
	"strings"
)

type SettlementTypeStore struct {
	db *sql.DB
}

func NewSettlementTypeStore(db *sql.DB) *SettlementTypeStore {
	return &SettlementTypeStore{db: db}
}

var (
	upsertSettlementTypesQuery = `
		INSERT INTO settlement_types (id, name, slug)
		VALUES %s
		ON CONFLICT (id) DO UPDATE
		SET name = EXCLUDED.name, slug = EXCLUDED.slug
	`
	findAllSettlementTypesQuery = `
		SELECT id, name, slug 
		FROM settlement_types
		ORDER BY id ASC;
	`
)

func (s *SettlementTypeStore) BulkUpsertTx(ctx context.Context, tx *sql.Tx, settlementTypes []*models.SettlementType) error {
	if len(settlementTypes) == 0 {
		return nil
	}

	valueStrings := make([]string, 0, len(settlementTypes))
	valueArgs := make([]any, 0, len(settlementTypes)*3)
	i := 1

	for _, st := range settlementTypes {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d)", i, i+1, i+2))
		valueArgs = append(valueArgs, st.ID, st.Name, st.Slug)
		i += 3
	}

	stmt := fmt.Sprintf(upsertSettlementTypesQuery, strings.Join(valueStrings, ","))
	_, err := tx.ExecContext(ctx, stmt, valueArgs...)
	return err
}

func (s *SettlementTypeStore) FindAll() ([]*models.SettlementType, error) {
	rows, err := s.db.QueryContext(context.Background(), findAllSettlementTypesQuery)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var settlementTypes []*models.SettlementType
	for rows.Next() {
		var st models.SettlementType
		if err := rows.Scan(
			&st.ID,
			&st.Name,
			&st.Slug,
		); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		settlementTypes = append(settlementTypes, &st)
	}

	return settlementTypes, nil
}
