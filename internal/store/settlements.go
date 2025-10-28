package store

import (
	"context"
	"database/sql"
	"fmt"
	"strade/internal/models"

	"github.com/lib/pq"
)

type SettlementStore struct {
	db *sql.DB
}

func NewSettlementStore(db *sql.DB) *SettlementStore {
	return &SettlementStore{db: db}
}

var (
	upsertSettlementQuery = `
		INSERT INTO settlements (
			postal_code, name, slug, settlement_type_id, 
			municipality_id, city_id, state_id, 
			office_postal_code, zone
		)
		SELECT 
			postal_code, name, slug, settlement_type_id,
			municipality_id, city_id, state_id,
			office_postal_code, zone
		FROM tmp_settlements
		ON CONFLICT (postal_code, name, municipality_id) DO UPDATE
		SET 
			slug = EXCLUDED.slug,
			settlement_type_id = EXCLUDED.settlement_type_id,
			city_id = EXCLUDED.city_id,
			state_id = EXCLUDED.state_id,
			office_postal_code = EXCLUDED.office_postal_code,
			zone = EXCLUDED.zone;
	`
)

func (s *SettlementStore) BulkUpsertTx(ctx context.Context, tx *sql.Tx, settlements []*models.Settlement) error {
	if len(settlements) == 0 {
		return nil
	}

	_, err := tx.ExecContext(ctx, `
		CREATE TEMP TABLE tmp_settlements (
			postal_code CHAR(5),
			name VARCHAR(100),
			slug VARCHAR(100),
			settlement_type_id CHAR(2),
			municipality_id CHAR(3),
			city_id CHAR(3),
			state_id CHAR(2),
			office_postal_code CHAR(5),
			zone VARCHAR(10)
		) ON COMMIT DROP;
	`)
	if err != nil {
		return fmt.Errorf("create temp table: %w", err)
	}

	stmt, err := tx.Prepare(pq.CopyIn("tmp_settlements", 
		"postal_code", "name", "slug", "settlement_type_id",
		"municipality_id", "city_id", "state_id",
		"office_postal_code", "zone"))
	if err != nil {
		return fmt.Errorf("prepare copy: %w", err)
	}

	for _, st := range settlements {
		if _, err := stmt.Exec(
			st.PostalCode, st.Name, st.Slug, st.SettlementTypeID,
			st.MunicipalityID, st.CityID, st.StateID,
			st.OfficePostalCode, st.Zone,
		); err != nil {
			_ = stmt.Close()
			return fmt.Errorf("copy row: %w", err)
		}
	}

	if _, err := stmt.Exec(); err != nil {
		return fmt.Errorf("final exec: %w", err)
	}
	if err := stmt.Close(); err != nil {
		return fmt.Errorf("close stmt: %w", err)
	}

	_, err = tx.ExecContext(ctx, upsertSettlementQuery)
	if err != nil {
		return fmt.Errorf("merge upsert: %w", err)
	}

	return nil
}

func (s *SettlementStore) FindAll() ([]*models.Settlement, error) {
	return nil, nil
}
