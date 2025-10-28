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
	return nil, nil
}
