# Next-Go-PG  •  Full-Stack Monorepo (Next.js 16 + Go + PostgreSQL)
#
# Run `just` to list recipes, `just <recipe>` to run one.
# https://github.com/casey/just

set shell := ["bash", "-cu"]

kamal_config := "-c deploy/config/deploy.yml"

# ─── Default: list recipes grouped ──────────────────────────────

[private]
default:
    @just --list --unsorted

# ─── Development ────────────────────────────────────────────────

# Start DB, frontend and backend
[group('dev')]
dev: db-up _migrate-up-silent api
    bun run dev:all

# Start everything (+ Grafana, Loki, Promtail)
[group('dev')]
dev-full: db-up _migrate-up-silent logs-up api
    @echo ""
    @echo "✓ Full dev environment ready"
    @echo ""
    @echo "  Frontend:  http://localhost:3000"
    @echo "  Backend:   http://localhost:8080"
    @echo "  Grafana:   http://localhost:3001  (admin/admin)"
    @echo ""
    bun run dev:all

# Start frontend only (localhost:3000)
[group('dev')]
dev-frontend:
    bun run dev:frontend

# Start DB and backend only (localhost:8080)
[group('dev')]
dev-backend: db-up
    cd backend && go run ./cmd/server

# ─── Database ───────────────────────────────────────────────────

# Start PostgreSQL database
[group('db')]
db-up:
    @docker compose -f deploy/compose/docker-compose.dev.yml up -d --wait

# Stop PostgreSQL database
[group('db')]
db-down:
    @docker compose -f deploy/compose/docker-compose.dev.yml down

# Reset database (delete all data)
[group('db')]
[confirm("Reset database? All data will be lost. [y/N]")]
db-reset:
    @docker compose -f deploy/compose/docker-compose.dev.yml down -v
    @docker compose -f deploy/compose/docker-compose.dev.yml up -d --wait

# Run Better Auth + GORM migrations
[group('db')]
db-migrate:
    @echo "Running Better Auth migrations..."
    cd frontend && bunx dotenv-cli -e .env.local -- bunx @better-auth/cli@latest migrate --config src/shared/lib/auth-server/auth.ts --yes
    @echo "✓ Migrations complete"

# ─── SQL Migrations (golang-migrate) ────────────────────────────

# Run all pending SQL migrations
[group('migrate')]
migrate-up:
    @echo "Running SQL migrations..."
    cd backend && go run ./cmd/migrate -up
    @echo "✓ Migrations complete"

[private]
_migrate-up-silent:
    @cd backend && go run ./cmd/migrate -up 2>/dev/null || true

# Rollback last SQL migration
[group('migrate')]
migrate-down:
    @echo "Rolling back last migration..."
    cd backend && go run ./cmd/migrate -down
    @echo "✓ Rollback complete"

# Show current migration version
[group('migrate')]
migrate-version:
    @cd backend && go run ./cmd/migrate -version

# Create new migration files (usage: just migrate-create <name>)
[group('migrate')]
migrate-create name:
    #!/usr/bin/env bash
    set -eu
    NEXT=$(ls backend/migrations/*.up.sql 2>/dev/null | wc -l | tr -d ' ')
    NEXT=$((NEXT + 1))
    PADDED=$(printf "%03d" "$NEXT")
    touch "backend/migrations/${PADDED}_{{ name }}.up.sql"
    touch "backend/migrations/${PADDED}_{{ name }}.down.sql"
    echo "✓ Created migrations:"
    echo "  backend/migrations/${PADDED}_{{ name }}.up.sql"
    echo "  backend/migrations/${PADDED}_{{ name }}.down.sql"

# ─── Build ──────────────────────────────────────────────────────

# Build frontend and backend
[group('build')]
build: build-frontend build-backend

# Build frontend for production
[group('build')]
build-frontend:
    bun run build:frontend

# Build backend binary
[group('build')]
build-backend:
    cd backend && go build -o bin/server ./cmd/server

# ─── Quality ────────────────────────────────────────────────────

# Run Biome linting
[group('quality')]
lint:
    cd frontend && bun run lint

# Fix linting issues automatically
[group('quality')]
lint-fix:
    cd frontend && bun run lint:fix

# Format code with Biome
[group('quality')]
format:
    cd frontend && bun run format

# Run TypeScript check
[group('quality')]
typecheck:
    cd frontend && bun run typecheck

# ─── Testing ────────────────────────────────────────────────────

# Run all tests
[group('test')]
test:
    bun run test

# Run frontend tests
[group('test')]
test-frontend:
    bun run --cwd frontend test

# Run backend tests
[group('test')]
test-backend:
    cd backend && go test ./...

# ─── Code Generation ────────────────────────────────────────────

# Generate Swagger docs only
[group('codegen')]
swagger:
    cd backend/cmd/server && ~/go/bin/swag init -g main.go -o ../../docs --parseDependency --dir .,../../internal/handler

# Generate TypeScript API client from OpenAPI
[group('codegen')]
api: swagger
    cd frontend && bunx orval

# Generate a new feature with Goca (usage: just goca-feature <Name> "<fields>")
[group('codegen')]
goca-feature name fields:
    #!/usr/bin/env bash
    set -eu
    cd backend && ~/go/bin/goca feature "{{ name }}" --fields "{{ fields }}"
    REGISTRY="{{ justfile_directory() }}/backend/internal/domain/registry.go"
    if grep -q "&{{ name }}{}" "$REGISTRY" 2>/dev/null; then
        echo "Entity already in registry"
    elif [ -f "$REGISTRY" ]; then
        sed -i '' "s|\(	// AUTO-GENERATED: New entities will be added above this line\)|	\&{{ name }}{},\n\1|" "$REGISTRY"
        echo "✓ Added &{{ name }}{} to registry.go"
    else
        echo "⚠ registry.go not found - add &{{ name }}{} manually"
    fi
    echo ""
    echo "━━━ Next Steps ━━━"
    echo ""
    echo "1. Run: just api"
    echo "2. Run: just dev-backend"
    echo ""

# ─── Setup ──────────────────────────────────────────────────────

# Install all dependencies
[group('setup')]
install: setup-hooks install-tools
    bun install
    cd frontend && bun install
    cd backend && go mod tidy

# Install required CLI tools (goca, gitleaks, sitefetch)
[group('setup')]
install-tools:
    #!/usr/bin/env bash
    echo "📦 Checking CLI tools..."
    command -v goca >/dev/null 2>&1 || { echo "Installing goca..." && go install github.com/sazardev/goca@latest; }
    command -v gitleaks >/dev/null 2>&1 || { echo "Installing gitleaks..." && brew install gitleaks 2>/dev/null || echo "  Skip: brew not available"; }
    command -v sitefetch >/dev/null 2>&1 || { echo "Installing sitefetch..." && bun install -g sitefetch 2>/dev/null || npm install -g sitefetch; }
    echo "✓ CLI tools ready"

# Setup git hooks for security scanning
[group('setup')]
setup-hooks:
    @./scripts/setup-hooks.sh

# Remove build artifacts and dependencies
[group('setup')]
clean:
    rm -rf frontend/.next frontend/node_modules
    rm -rf backend/bin
    rm -rf node_modules

# ─── Docker ─────────────────────────────────────────────────────

# Build production Docker images
[group('docker')]
docker-build:
    docker compose build

# Start production containers
[group('docker')]
docker-up:
    docker compose up -d

# Stop production containers
[group('docker')]
docker-down:
    docker compose down

# ─── Deployment (Kamal) ─────────────────────────────────────────

# Deploy to staging environment
[group('deploy')]
deploy-staging:
    kamal deploy {{ kamal_config }} -d staging

# Deploy to production (usage: just deploy-production)
[group('deploy')]
[confirm("⚠  Deploy to PRODUCTION? [y/N]")]
deploy-production:
    kamal deploy {{ kamal_config }} -d production

# Rollback deployment (usage: just deploy-rollback <staging|production>)
[group('deploy')]
deploy-rollback dest:
    kamal rollback {{ kamal_config }} -d {{ dest }}

# Follow deployment logs (usage: just deploy-logs <staging|production>)
[group('deploy')]
deploy-logs dest:
    kamal app logs {{ kamal_config }} -d {{ dest }} -f

# Open console on server (usage: just deploy-console <staging|production>)
[group('deploy')]
deploy-console dest:
    kamal app exec {{ kamal_config }} -d {{ dest }} -i bash

# Initial Kamal setup on server (usage: just deploy-setup <staging|production>)
[group('deploy')]
deploy-setup dest:
    kamal setup {{ kamal_config }} -d {{ dest }}

# ─── Security ───────────────────────────────────────────────────

# Scan current files for secrets
[group('security')]
security-scan:
    #!/usr/bin/env bash
    echo "🔒 Scanning for secrets and sensitive data..."
    if command -v gitleaks >/dev/null 2>&1; then
        gitleaks detect --config .gitleaks.toml --no-git --verbose
    else
        echo "gitleaks not installed. Install with:"
        echo "  brew install gitleaks"
    fi

# Scan entire git history for secrets
[group('security')]
security-scan-history:
    #!/usr/bin/env bash
    echo "🔒 Scanning git history for secrets..."
    if command -v gitleaks >/dev/null 2>&1; then
        gitleaks detect --config .gitleaks.toml --verbose
    else
        echo "gitleaks not installed. Install with:"
        echo "  brew install gitleaks"
    fi

# ─── Documentation ──────────────────────────────────────────────

# Fetch LLM-friendly docs (usage: just fetch-docs <url> [name])
[group('docs')]
fetch-docs url name='':
    @./scripts/fetch-docs.sh "{{ url }}" "{{ name }}"

# Semantic search (usage: just search-docs "query" [n] [fast])
[group('docs')]
search-docs q n='5' fast='':
    @bun scripts/search-docs.js "{{ q }}" --top {{ n }} --llm {{ if fast == "" { "" } else { "--fast" } }}

# Pre-build search index (one-time, makes search fast)
[group('docs')]
search-docs-index:
    @echo "🔄 Building semantic search index..."
    @bun scripts/search-docs.js --index
    @echo "✓ Index ready! Searches will now be fast."

# ─── Logging (Grafana + Loki) ───────────────────────────────────

# Start logging stack (Grafana, Loki, Promtail)
[group('logs')]
logs-up:
    @echo "📊 Starting logging stack..."
    @docker compose -f deploy/compose/docker-compose.dev.yml -f deploy/compose/docker-compose.logging.yml up -d loki promtail grafana
    @echo "✓ Logging stack started"
    @echo ""
    @echo "  Grafana:  http://localhost:3001  (admin/admin)"
    @echo "  Loki:     http://localhost:3100"
    @echo ""

# Stop logging stack
[group('logs')]
logs-down:
    @echo "Stopping logging stack..."
    @docker compose -f deploy/compose/docker-compose.dev.yml -f deploy/compose/docker-compose.logging.yml stop loki promtail grafana
    @docker compose -f deploy/compose/docker-compose.dev.yml -f deploy/compose/docker-compose.logging.yml rm -f loki promtail grafana
    @echo "✓ Logging stack stopped"

# Open Grafana in browser (localhost:3001)
[group('logs')]
logs-open:
    @echo "Opening Grafana..."
    @open http://localhost:3001 2>/dev/null || xdg-open http://localhost:3001 2>/dev/null || echo "Open http://localhost:3001 in your browser"

# Query logs via CLI (usage: just logs-query '<logql>' [limit])
[group('logs')]
logs-query q limit='100':
    @curl -sG "http://localhost:3100/loki/api/v1/query_range" --data-urlencode "query={{ q }}" --data-urlencode "limit={{ limit }}" | jq -r '.data.result[].values[][1]' 2>/dev/null || echo "Error: Loki not running or query failed"

# ─── Monitoring ─────────────────────────────────────────────────

# View Prometheus metrics (localhost:8080/metrics)
[group('monitoring')]
metrics:
    @echo "Opening Prometheus metrics..."
    @open http://localhost:8080/metrics 2>/dev/null || xdg-open http://localhost:8080/metrics 2>/dev/null || echo "Open http://localhost:8080/metrics in your browser"

# ─── Database Backups ───────────────────────────────────────────

# Start automatic backup (daily to S3)
[group('backup')]
backup-up:
    @echo "Starting automatic backup system..."
    @docker compose -f deploy/compose/docker-compose.dev.yml -f deploy/compose/docker-compose.backup.yml up -d rustfs rustfs-init pg-backup
    @echo ""
    @echo "✓ Backup system running"
    @echo ""
    @echo "  Schedule:   Daily (change with BACKUP_SCHEDULE)"
    @echo "  Retention:  7 days (change with BACKUP_KEEP_DAYS)"
    @echo "  Storage:    RustFS S3 (localhost:9001 for UI)"
    @echo ""

# Stop backup stack
[group('backup')]
backup-down:
    @echo "Stopping backup system..."
    @docker compose -f deploy/compose/docker-compose.dev.yml -f deploy/compose/docker-compose.backup.yml stop pg-backup rustfs
    @docker compose -f deploy/compose/docker-compose.dev.yml -f deploy/compose/docker-compose.backup.yml rm -f pg-backup rustfs rustfs-init
    @echo "✓ Backup system stopped"

# Create manual backup now
[group('backup')]
backup-now:
    @echo "Creating backup..."
    @docker compose -f deploy/compose/docker-compose.dev.yml -f deploy/compose/docker-compose.backup.yml exec pg-backup sh backup.sh
    @echo "✓ Backup complete"

# List all backups
[group('backup')]
backup-list:
    @echo "Backups in S3:"
    @docker run --rm --network next-go-pg_default --entrypoint sh minio/mc:latest -c "mc alias set s3 http://rustfs:9000 rustfsadmin rustfsadmin >/dev/null 2>&1 && mc ls s3/backups/postgres/" 2>/dev/null || echo "No backups found or backup system not running"

# Restore from latest backup
[group('backup')]
[confirm("⚠  This will overwrite the current database. Continue? [y/N]")]
backup-restore:
    @echo "Restoring from latest backup..."
    @docker compose -f deploy/compose/docker-compose.dev.yml -f deploy/compose/docker-compose.backup.yml exec pg-backup sh restore.sh
    @echo "✓ Database restored"
