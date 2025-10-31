ifneq ($(ENVIRONMENT),PRODUCTION)
  ifeq ("$(wildcard .env)","")
    $(warning .env file not found, skipping...)
  else
    include .env
    export $(shell sed 's/=.*//' .env)
  endif
endif

MIGRATIONS_PATH = ./cmd/migrate/migrations
MIGRATE_COMMAND = migrate -path "$(MIGRATIONS_PATH)" -database "$(DB_ADDR)"

.PHONY: create-migration
create-migration:
	@echo "Creating new migration"
	@migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate:
	@echo "Migrating database"
	@$(MIGRATE_COMMAND) up

.PHONY: migrate-down
migrate-down:
	@echo "Rolling back database"
	@$(MIGRATE_COMMAND) down

.PHONY: migrate-reset
migrate-reset: 
	@echo "Resetting database"
	@$(MIGRATE_COMMAND) drop
	@$(MIGRATE_COMMAND) up

.PHONY: shift-migrations
shift-migrations:
	@./scripts/shiftmigrations.sh $(filter-out $@,$(MAKECMDGOALS))

.PHONY: docker-dev-up
docker-dev-up:
	@echo "Starting development environment (db + redis only)..."
	@docker-compose -f docker/docker-compose.dev.yml up -d db redis

.PHONY: docker-dev-full-up
docker-dev-full-up:
	@echo "Starting full development environment (db + redis + api + watcher)..."
	@docker-compose -f docker/docker-compose.dev.yml up -d

.PHONY: docker-dev-down
docker-dev-down:
	@echo "Stopping development environment..."
	@docker-compose -f docker/docker-compose.dev.yml down

.PHONY: docker-dev-logs
docker-dev-logs:
	@docker-compose -f docker/docker-compose.dev.yml logs -f

.PHONY: docker-up
docker-up:
	@echo "Starting full environment (db + redis + api + watcher)..."
	@docker-compose -f docker/docker-compose.yml up -d

.PHONY: docker-down
docker-down:
	@echo "Stopping full environment..."
	@docker-compose -f docker/docker-compose.yml down

.PHONY: docker-logs
docker-logs:
	@docker-compose -f docker/docker-compose.yml logs -f

.PHONY: docker-clean
docker-clean:
	@echo "Cleaning up all Docker resources..."
	@docker-compose -f docker/docker-compose.yml down -v
	@docker-compose -f docker/docker-compose.dev.yml down -v
