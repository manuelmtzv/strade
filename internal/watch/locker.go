package watch

import (
	"context"
	"database/sql"
	"hash/fnv"
	"time"
)

type Locker struct {
	db *sql.DB
}

func NewLocker(db *sql.DB) *Locker {
	return &Locker{db: db}
}

func (l *Locker) TryLock(ctx context.Context, key string, ttl time.Duration) (bool, func(), error) {
	lockID := hashKey(key)

	var acquired bool
	err := l.db.QueryRowContext(ctx, "SELECT pg_try_advisory_lock($1)", lockID).Scan(&acquired)
	if err != nil {
		return false, nil, err
	}

	if !acquired {
		return false, nil, nil
	}

	release := func() {
		_, _ = l.db.ExecContext(context.Background(), "SELECT pg_advisory_unlock($1)", lockID)
	}

	return true, release, nil
}

func hashKey(key string) int64 {
	h := fnv.New64a()
	h.Write([]byte(key))
	return int64(h.Sum64())
}
