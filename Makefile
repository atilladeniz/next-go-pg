.PHONY: help dev dev-frontend dev-backend build build-frontend build-backend test clean install api lint docker-build docker-up docker-down goca-feature deploy deploy-staging deploy-production deploy-rollback deploy-logs deploy-console setup-hooks security-scan security-scan-history

.DEFAULT_GOAL := help

# Colors
CYAN := \033[36m
GREEN := \033[32m
YELLOW := \033[33m
DIM := \033[2m
RESET := \033[0m

help: ## Show available commands
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

install: setup-hooks
	bun install
	cd frontend && bun install
	cd backend && go mod tidy

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
