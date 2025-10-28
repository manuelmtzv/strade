package store

import (
	"context"
	"database/sql"
	"errors"
)

type WatermarkStore struct {
	db *sql.DB
}

func NewWatermarkStore(db *sql.DB) *WatermarkStore {
	return &WatermarkStore{db: db}
}

func (s *WatermarkStore) Get(ctx context.Context, key string) (string, error) {
	var value string
	query := `SELECT value FROM watermarks WHERE key = $1`
	err := s.db.QueryRowContext(ctx, query, key).Scan(&value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return value, nil
}

func (s *WatermarkStore) Set(ctx context.Context, key, value string) error {
	query := `
		INSERT INTO watermarks (key, value, updated_at)
		VALUES ($1, $2, CURRENT_TIMESTAMP)
		ON CONFLICT (key) 
		DO UPDATE SET value = EXCLUDED.value, updated_at = CURRENT_TIMESTAMP
	`
	_, err := s.db.ExecContext(ctx, query, key, value)
	return err
}
