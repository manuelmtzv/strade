#!/bin/sh
set -e

echo "Running database migrations..."
migrate -path cmd/migrate/migrations -database "$DB_ADDR" up

echo "Starting application..."
exec "$@"
