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
			id, postal_code, name, slug, settlement_type_id, 
			municipality_id, city_id, state_id, 
			office_postal_code, zone
		)
		SELECT 
			id, postal_code, name, slug, settlement_type_id,
			municipality_id, city_id, state_id,
			office_postal_code, zone
		FROM tmp_settlements
		ON CONFLICT (id) DO UPDATE
		SET 
			postal_code = EXCLUDED.postal_code,
			name = EXCLUDED.name,
			slug = EXCLUDED.slug,
			settlement_type_id = EXCLUDED.settlement_type_id,
			municipality_id = EXCLUDED.municipality_id,
			city_id = EXCLUDED.city_id,
			state_id = EXCLUDED.state_id,
			office_postal_code = EXCLUDED.office_postal_code,
			zone = EXCLUDED.zone;
	`
	findSettlementByPostalCodeQuery = `
		SELECT
			s.id,
			s.postal_code, s.name, s.slug,
			s.office_postal_code, s.zone,
			stt.id, stt.name, stt.slug,
			m.id, m.name, m.slug,
			c.id, c.name, c.slug,
			st.id, st.name, st.slug
		FROM settlements s
		LEFT JOIN settlement_types stt ON s.settlement_type_id = stt.id
		LEFT JOIN municipalities m ON s.municipality_id = m.id
		LEFT JOIN cities c ON s.city_id = c.id
		LEFT JOIN states st ON s.state_id = st.id
		WHERE postal_code = $1;
	`
	searchSettlementByNameQuery = `
		SELECT 
            s.id, s.postal_code, s.name, s.slug,
            s.office_postal_code, s.zone,
            stt.id, stt.name, stt.slug,
            m.id, m.name, m.slug,
            c.id, c.name, c.slug,
            st.id, st.name, st.slug
        FROM settlements s
        LEFT JOIN settlement_types stt ON s.settlement_type_id = stt.id
        LEFT JOIN municipalities m ON s.municipality_id = m.id
        LEFT JOIN cities c ON s.city_id = c.id
        LEFT JOIN states st ON s.state_id = st.id
        WHERE (s.name ILIKE $1 OR s.slug ILIKE $1) 
        LIMIT $2
	`
)

func (s *SettlementStore) BulkUpsertTx(ctx context.Context, tx *sql.Tx, settlements []*models.Settlement) error {
	if len(settlements) == 0 {
		return nil
	}

	_, err := tx.ExecContext(ctx, `
		CREATE TEMP TABLE tmp_settlements (
			id CHAR(9),
			postal_code CHAR(5),
			name VARCHAR(100),
			slug VARCHAR(100),
			settlement_type_id CHAR(2),
			municipality_id CHAR(5),
			city_id VARCHAR(10),
			state_id CHAR(2),
			office_postal_code CHAR(5),
			zone VARCHAR(10)
		) ON COMMIT DROP;
	`)
	if err != nil {
		return fmt.Errorf("create temp table: %w", err)
	}

	stmt, err := tx.Prepare(pq.CopyIn("tmp_settlements",
		"id", "postal_code", "name", "slug", "settlement_type_id",
		"municipality_id", "city_id", "state_id",
		"office_postal_code", "zone"))
	if err != nil {
		return fmt.Errorf("prepare copy: %w", err)
	}

	for _, st := range settlements {
		if _, err := stmt.Exec(
			st.ID, st.PostalCode, st.Name, st.Slug, st.SettlementTypeID,
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

func (s *SettlementStore) FindByPostalCode(ctx context.Context, postalCode string) ([]*models.Settlement, error) {
	rows, err := s.db.QueryContext(ctx, findSettlementByPostalCodeQuery, postalCode)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var settlements []*models.Settlement
	for rows.Next() {
		var st models.Settlement
		st.SettlementType = &models.SettlementType{}
		st.Municipality = &models.Municipality{}
		st.City = &models.City{}
		st.State = &models.State{}

		if err := rows.Scan(
			&st.ID,
			&st.PostalCode,
			&st.Name,
			&st.Slug,
			&st.OfficePostalCode,
			&st.Zone,
			&st.SettlementType.ID,
			&st.SettlementType.Name,
			&st.SettlementType.Slug,
			&st.Municipality.ID,
			&st.Municipality.Name,
			&st.Municipality.Slug,
			&st.City.ID,
			&st.City.Name,
			&st.City.Slug,
			&st.State.ID,
			&st.State.Name,
			&st.State.Slug,
		); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		settlements = append(settlements, &st)
	}

	return settlements, nil
}

func (s *SettlementStore) SearchByName(ctx context.Context, name string, limit int) ([]*models.Settlement, error) {
	rows, err := s.db.QueryContext(ctx, searchSettlementByNameQuery, name, limit)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var settlements []*models.Settlement
	for rows.Next() {
		var st models.Settlement
		st.SettlementType = &models.SettlementType{}
		st.Municipality = &models.Municipality{}
		st.City = &models.City{}
		st.State = &models.State{}

		if err := rows.Scan(
			&st.ID,
			&st.PostalCode,
			&st.Name,
			&st.Slug,
			&st.OfficePostalCode,
			&st.Zone,
			&st.SettlementType.ID,
			&st.SettlementType.Name,
			&st.SettlementType.Slug,
			&st.Municipality.ID,
			&st.Municipality.Name,
			&st.Municipality.Slug,
			&st.City.ID,
			&st.City.Name,
			&st.City.Slug,
			&st.State.ID,
			&st.State.Name,
			&st.State.Slug,
		); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		settlements = append(settlements, &st)
	}

	return settlements, nil
}
