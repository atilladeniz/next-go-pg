.PHONY: help dev dev-frontend dev-backend build build-frontend build-backend test clean install api lint docker-build docker-up docker-down goca-feature

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
	@echo "  $(GREEN)clean$(RESET)             Remove build artifacts and dependencies"
	@echo ""
	@echo "$(CYAN)━━━ Docker ━━━$(RESET)"
	@echo ""
	@echo "  $(GREEN)docker-build$(RESET)      Build production Docker images"
	@echo "  $(GREEN)docker-up$(RESET)         Start production containers"
	@echo "  $(GREEN)docker-down$(RESET)       Stop production containers"
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

install:
	bun install
	cd frontend && bun install
	cd backend && go mod tidy

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
