GOOSE_DRIVER       := postgres
GOOSE_DBSTRING     := host=localhost port=5432 user=main password=password dbname=event-management sslmode=disable
GOOSE_MIGRATION_DIR := ./event-management/migrations
CONNECT_URL        := http://localhost:8083
SCHEMA_VERSION     := 20260403185907

export GOOSE_DRIVER GOOSE_DBSTRING GOOSE_MIGRATION_DIR

# ── Docker ──────────────────────────────────
.PHONY: up down reset

up:
	docker compose up -d

up-components:
	docker compose up -d elasticsearch broker event-management-db redis connect

down:
	docker compose down

reset:
	docker compose down -v

# ── Wait helpers ────────────────────────────
.PHONY: wait-db wait-connect

wait-db:
	@printf "Waiting for PostgreSQL..."
	@until docker compose exec -T event-management-db pg_isready -U main > /dev/null 2>&1; do sleep 1; printf "."; done
	@echo " ready"

wait-connect:
	@printf "Waiting for Kafka Connect..."
	@until curl -sf $(CONNECT_URL)/ > /dev/null 2>&1; do sleep 2; printf "."; done
	@echo " ready"

# ── Migrations ──────────────────────────────
.PHONY: migrate-schema migrate-seed migrate-down migrate-status

migrate-schema: wait-db
	goose up-to $(SCHEMA_VERSION)

migrate-seed: wait-db
	goose up

migrate-down: wait-db
	goose reset

migrate-status: wait-db
	goose status

# ── Connector ───────────────────────────────
.PHONY: connector-create connector-delete connector-status

connector-create: wait-connect
	@curl -s -X POST -H "Accept:application/json" -H "Content-Type:application/json" \
			$(CONNECT_URL)/connectors/ -d @connector-config.json
	@echo ""
	@printf "Waiting for snapshot to complete..."
	@sleep 3
	@until curl -sf $(CONNECT_URL)/connectors/postgres-cdc/status 2>/dev/null | grep -q '"RUNNING"'; do sleep 2; printf "."; done
	@echo " done"

connector-delete:
	@curl -s -X DELETE $(CONNECT_URL)/connectors/postgres-cdc || true
	@echo ""

connector-status:
	@curl -s $(CONNECT_URL)/connectors/postgres-cdc/status | python3 -m json.tool 2>/dev/null || echo "Connector not found"

# ── Workflows ───────────────────────────────
.PHONY: setup teardown

setup: migrate-schema connector-create migrate-seed
	@echo "✓ CDC setup complete"

teardown: connector-delete migrate-down
	@echo "✓ CDC teardown complete"