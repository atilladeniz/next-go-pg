.PHONY: help dev dev-full dev-frontend dev-backend build build-frontend build-backend test clean install install-tools api lint docker-build docker-up docker-down goca-feature deploy deploy-staging deploy-production deploy-rollback deploy-logs deploy-console setup-hooks security-scan security-scan-history fetch-docs search-docs search-docs-index logs-up logs-down logs-open logs-query db-migrate

.DEFAULT_GOAL := help

# Colors
CYAN := \033[36m
GREEN := \033[32m
YELLOW := \033[33m
DIM := \033[2m
RESET := \033[0m

help: ## Show available commands
	@echo ""
	@echo "$(CYAN)    _   __          __       _          ____"
	@echo "   / | / /__  _  __/ /_     (_)____    / ___/____"
	@echo "  /  |/ / _ \\| |/_/ __/    / / ___/   / __ \`/ __ \\"
	@echo " / /|  /  __/>  </ /_     / (__  )   / /_/ / /_/ /"
	@echo "/_/ |_/\\___/_/|_|\\__/ + _/ /____/    \\__, /\\____/"
	@echo "                      /___/        /____/$(RESET)"
	@echo ""
	@echo "$(DIM)Full-Stack Monorepo • Next.js 16 + Go + PostgreSQL$(RESET)"
	@echo ""
	@echo "$(CYAN)━━━ Development ━━━$(RESET)"
	@echo ""
	@echo "  $(GREEN)dev$(RESET)               Start DB, frontend and backend"
	@echo "  $(GREEN)dev-frontend$(RESET)      Start frontend only (localhost:3000)"
	@echo "  $(GREEN)dev-backend$(RESET)       Start DB and backend only (localhost:8080)"
	@echo ""
	@echo "$(CYAN)━━━ Database ━━━$(RESET)"
	@echo ""
	@echo "  $(GREEN)db-up$(RESET)             Start PostgreSQL database"
	@echo "  $(GREEN)db-down$(RESET)           Stop PostgreSQL database"
	@echo "  $(GREEN)db-reset$(RESET)          Reset database (delete all data)"
	@echo "  $(GREEN)db-migrate$(RESET)        Run Better Auth + GORM migrations"
	@echo ""
	@echo "$(CYAN)━━━ Build ━━━$(RESET)"
	@echo ""
	@echo "  $(GREEN)build$(RESET)             Build frontend and backend"
	@echo "  $(GREEN)build-frontend$(RESET)    Build frontend for production"
	@echo "  $(GREEN)build-backend$(RESET)     Build backend binary"
	@echo ""
	@echo "$(CYAN)━━━ Quality ━━━$(RESET)"
	@echo ""
	@echo "  $(GREEN)lint$(RESET)              Run Biome linting"
	@echo "  $(GREEN)lint-fix$(RESET)          Fix linting issues automatically"
	@echo "  $(GREEN)format$(RESET)            Format code with Biome"
	@echo "  $(GREEN)typecheck$(RESET)         Run TypeScript check"
	@echo ""
	@echo "$(CYAN)━━━ Testing ━━━$(RESET)"
	@echo ""
	@echo "  $(GREEN)test$(RESET)              Run all tests"
	@echo "  $(GREEN)test-frontend$(RESET)     Run frontend tests"
	@echo "  $(GREEN)test-backend$(RESET)      Run backend tests"
	@echo ""
	@echo "$(CYAN)━━━ Code Generation ━━━$(RESET)"
	@echo ""
	@echo "  $(GREEN)api$(RESET)               Generate TypeScript API client from OpenAPI"
	@echo "  $(GREEN)swagger$(RESET)           Generate Swagger docs only"
	@echo "  $(GREEN)goca-feature$(RESET)      Generate a new feature with Goca $(DIM)(+ registry hint)$(RESET)"
	@echo ""
	@echo "$(CYAN)━━━ Setup ━━━$(RESET)"
	@echo ""
	@echo "  $(GREEN)install$(RESET)           Install all dependencies"
	@echo "  $(GREEN)setup-hooks$(RESET)       Setup git hooks for security scanning"
	@echo "  $(GREEN)clean$(RESET)             Remove build artifacts and dependencies"
	@echo ""
	@echo "$(CYAN)━━━ Security ━━━$(RESET)"
	@echo ""
	@echo "  $(GREEN)security-scan$(RESET)     Scan current files for secrets"
	@echo "  $(GREEN)security-scan-history$(RESET)  Scan entire git history"
	@echo ""
	@echo "$(CYAN)━━━ Documentation ━━━$(RESET)"
	@echo ""
	@echo "  $(GREEN)fetch-docs$(RESET)        Fetch LLM-friendly docs $(DIM)(url=<url> [name=<name>])$(RESET)"
	@echo "  $(GREEN)search-docs$(RESET)       Semantic search $(DIM)(q=\"query\" [n=5] [fast=1])$(RESET)"
	@echo "  $(GREEN)search-docs-index$(RESET) Pre-build search index $(DIM)(one-time, makes search fast)$(RESET)"
	@echo ""
	@echo "$(CYAN)━━━ Docker ━━━$(RESET)"
	@echo ""
	@echo "  $(GREEN)docker-build$(RESET)      Build production Docker images"
	@echo "  $(GREEN)docker-up$(RESET)         Start production containers"
	@echo "  $(GREEN)docker-down$(RESET)       Stop production containers"
	@echo ""
	@echo "$(CYAN)━━━ Deployment (Kamal) ━━━$(RESET)"
	@echo ""
	@echo "  $(GREEN)deploy-staging$(RESET)    Deploy to staging environment"
	@echo "  $(GREEN)deploy-production$(RESET) Deploy to production environment"
	@echo "  $(GREEN)deploy-rollback$(RESET)   Rollback to previous version"
	@echo "  $(GREEN)deploy-logs$(RESET)       Show deployment logs"
	@echo "  $(GREEN)deploy-console$(RESET)    Open console on production server"
	@echo ""
	@echo "$(CYAN)━━━ Logging (Grafana + Loki) ━━━$(RESET)"
	@echo ""
	@echo "  $(GREEN)logs-up$(RESET)           Start logging stack (Grafana, Loki, Promtail)"
	@echo "  $(GREEN)logs-down$(RESET)         Stop logging stack"
	@echo "  $(GREEN)logs-open$(RESET)         Open Grafana in browser $(DIM)(localhost:3001)$(RESET)"
	@echo "  $(GREEN)logs-query$(RESET)        Query logs via CLI $(DIM)(q=\"query\" [limit=100])$(RESET)"
	@echo ""

# ━━━ Development ━━━

dev: db-up api
	@echo "Starting frontend and backend..."
	bun run dev:all

dev-frontend:
	bun run dev:frontend

dev-backend: db-up
	cd backend && go run ./cmd/server

# ━━━ Database ━━━

db-up:
	@docker compose -f docker-compose.dev.yml up -d --wait

db-down:
	@docker compose -f docker-compose.dev.yml down

db-reset:
	@docker compose -f docker-compose.dev.yml down -v
	@docker compose -f docker-compose.dev.yml up -d --wait

db-migrate:
	@echo "$(YELLOW)Running Better Auth migrations...$(RESET)"
	cd frontend && bunx dotenv-cli -e .env.local -- bunx @better-auth/cli@latest migrate --config src/shared/lib/auth-server/auth.ts --yes
	@echo "$(GREEN)✓ Migrations complete$(RESET)"

# ━━━ Build ━━━

build: build-frontend build-backend

build-frontend:
	bun run build:frontend

build-backend:
	cd backend && go build -o bin/server ./cmd/server

# ━━━ Quality ━━━

lint:
	cd frontend && bun run lint

lint-fix:
	cd frontend && bun run lint:fix

format:
	cd frontend && bun run format

typecheck:
	cd frontend && bun run typecheck

# ━━━ Testing ━━━

test:
	bun run test

test-backend:
	cd backend && go test ./...

test-frontend:
	bun run --cwd frontend test

# ━━━ Code Generation ━━━

swagger:
	cd backend && ~/go/bin/swag init -g cmd/server/main.go -o docs --parseDependency

api: swagger
	cd frontend && bunx orval

goca-feature:
	@read -p "Feature name: " name; \
	read -p "Fields (e.g. name:string,email:string): " fields; \
	cd backend && ~/go/bin/goca feature $$name --fields "$$fields"; \
	REGISTRY="$(CURDIR)/backend/internal/domain/registry.go"; \
	if grep -q "&$$name{}" "$$REGISTRY" 2>/dev/null; then \
		echo "$(DIM)Entity already in registry$(RESET)"; \
	elif [ -f "$$REGISTRY" ]; then \
		sed -i '' "s|\(	// AUTO-GENERATED: New entities will be added above this line\)|	\&$$name{},\n\1|" "$$REGISTRY"; \
		echo "$(GREEN)✓ Added &$$name{} to registry.go$(RESET)"; \
	else \
		echo "$(YELLOW)⚠ registry.go not found - add &$$name{} manually$(RESET)"; \
	fi; \
	echo ""; \
	echo "$(YELLOW)━━━ Next Steps ━━━$(RESET)"; \
	echo ""; \
	echo "1. Run: $(GREEN)make api$(RESET)"; \
	echo "2. Run: $(GREEN)make dev-backend$(RESET)"; \
	echo ""

# ━━━ Setup ━━━

install: setup-hooks install-tools
	bun install
	cd frontend && bun install
	cd backend && go mod tidy

install-tools:
	@echo "$(YELLOW)📦 Checking CLI tools...$(RESET)"
	@command -v goca >/dev/null 2>&1 || { echo "$(YELLOW)Installing goca...$(RESET)" && go install github.com/sazardev/goca@latest; }
	@command -v gitleaks >/dev/null 2>&1 || { echo "$(YELLOW)Installing gitleaks...$(RESET)" && brew install gitleaks 2>/dev/null || echo "$(DIM)  Skip: brew not available$(RESET)"; }
	@command -v sitefetch >/dev/null 2>&1 || { echo "$(YELLOW)Installing sitefetch...$(RESET)" && bun install -g sitefetch 2>/dev/null || npm install -g sitefetch; }
	@echo "$(GREEN)✓ CLI tools ready$(RESET)"

setup-hooks:
	@./scripts/setup-hooks.sh

clean:
	rm -rf frontend/.next frontend/node_modules
	rm -rf backend/bin
	rm -rf node_modules

# ━━━ Docker ━━━

docker-build:
	docker compose build

docker-up:
	docker compose up -d

docker-down:
	docker compose down

# ━━━ Deployment (Kamal) ━━━

KAMAL_CONFIG := -c deploy/config/deploy.yml

deploy-staging:
	kamal deploy $(KAMAL_CONFIG) -d staging

deploy-production:
	@echo "$(YELLOW)⚠ Deploying to PRODUCTION!$(RESET)"
	@read -p "Are you sure? [y/N] " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		kamal deploy $(KAMAL_CONFIG) -d production; \
	else \
		echo "Aborted."; \
	fi

deploy-rollback:
	@read -p "Destination (staging/production): " dest; \
	kamal rollback $(KAMAL_CONFIG) -d $$dest

deploy-logs:
	@read -p "Destination (staging/production): " dest; \
	kamal app logs $(KAMAL_CONFIG) -d $$dest -f

deploy-console:
	@read -p "Destination (staging/production): " dest; \
	kamal app exec $(KAMAL_CONFIG) -d $$dest -i bash

deploy-setup:
	@read -p "Destination (staging/production): " dest; \
	kamal setup $(KAMAL_CONFIG) -d $$dest

# ━━━ Security ━━━

security-scan:
	@echo "$(YELLOW)🔒 Scanning for secrets and sensitive data...$(RESET)"
	@if command -v gitleaks >/dev/null 2>&1; then \
		gitleaks detect --config .gitleaks.toml --no-git --verbose; \
	else \
		echo "$(RED)gitleaks not installed. Install with:$(RESET)"; \
		echo "  brew install gitleaks"; \
	fi

security-scan-history:
	@echo "$(YELLOW)🔒 Scanning git history for secrets...$(RESET)"
	@if command -v gitleaks >/dev/null 2>&1; then \
		gitleaks detect --config .gitleaks.toml --verbose; \
	else \
		echo "$(RED)gitleaks not installed. Install with:$(RESET)"; \
		echo "  brew install gitleaks"; \
	fi

# ━━━ Documentation ━━━

fetch-docs:
ifndef url
	@echo "$(RED)Error: url parameter required$(RESET)"
	@echo ""
	@echo "Usage: make fetch-docs url=<url> [name=<name>]"
	@echo ""
	@echo "Examples:"
	@echo "  make fetch-docs url=https://tanstack.com/query/latest"
	@echo "  make fetch-docs url=https://nextjs.org/docs name=nextjs"
	@echo "  make fetch-docs url=https://orm.drizzle.team name=drizzle"
	@exit 1
else
	@./scripts/fetch-docs.sh "$(url)" "$(name)"
endif

search-docs:
ifndef q
	@echo "$(RED)Error: q parameter required$(RESET)"
	@echo ""
	@echo "Usage: make search-docs q=\"your query\" [n=5] [fast=1]"
	@echo ""
	@echo "Examples:"
	@echo "  make search-docs q=\"how to preload data\"        $(DIM)# Semantic (default)$(RESET)"
	@echo "  make search-docs q=\"prefetchQuery\" fast=1       $(DIM)# Fast fuzzy search$(RESET)"
	@echo "  make search-docs q=\"mutations\" n=3"
	@exit 1
else
	@bun scripts/search-docs.js "$(q)" --top $(or $(n),5) --llm $(if $(fast),--fast,)
endif

search-docs-index:
	@echo "$(YELLOW)🔄 Building semantic search index...$(RESET)"
	@bun scripts/search-docs.js --index
	@echo "$(GREEN)✓ Index ready! Searches will now be fast.$(RESET)"

# ━━━ Logging (Grafana + Loki) ━━━

logs-up:
	@echo "$(YELLOW)📊 Starting logging stack...$(RESET)"
	@docker compose -f docker-compose.dev.yml -f docker-compose.logging.yml up -d loki promtail grafana
	@echo "$(GREEN)✓ Logging stack started$(RESET)"
	@echo ""
	@echo "  Grafana:  $(CYAN)http://localhost:3001$(RESET) $(DIM)(admin/admin)$(RESET)"
	@echo "  Loki:     $(CYAN)http://localhost:3100$(RESET)"
	@echo ""

logs-down:
	@echo "$(YELLOW)Stopping logging stack...$(RESET)"
	@docker compose -f docker-compose.dev.yml -f docker-compose.logging.yml stop loki promtail grafana
	@docker compose -f docker-compose.dev.yml -f docker-compose.logging.yml rm -f loki promtail grafana
	@echo "$(GREEN)✓ Logging stack stopped$(RESET)"

logs-open:
	@echo "$(CYAN)Opening Grafana...$(RESET)"
	@open http://localhost:3001 2>/dev/null || xdg-open http://localhost:3001 2>/dev/null || echo "Open http://localhost:3001 in your browser"

logs-query:
ifndef q
	@echo "$(RED)Error: q parameter required$(RESET)"
	@echo ""
	@echo "Usage: make logs-query q=\"your query\" [limit=100]"
	@echo ""
	@echo "Examples:"
	@echo "  make logs-query q='{service=\"next-go-pg-api\"}'"
	@echo "  make logs-query q='{level=\"error\"}' limit=50"
	@echo "  make logs-query q='{category=\"auth\"}'"
	@exit 1
else
	@curl -sG "http://localhost:3100/loki/api/v1/query_range" \
		--data-urlencode "query=$(q)" \
		--data-urlencode "limit=$(or $(limit),100)" | \
		jq -r '.data.result[].values[][1]' 2>/dev/null || \
		echo "$(RED)Error: Loki not running or query failed$(RESET)"
endif
