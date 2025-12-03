.PHONY: help dev dev-frontend dev-backend build build-frontend build-backend test clean install lint docker-build docker-up docker-down goca-feature

.DEFAULT_GOAL := help

# Colors
CYAN := \033[36m
GREEN := \033[32m
YELLOW := \033[33m
RESET := \033[0m

help: ## Show available commands
	@echo ""
	@echo "$(CYAN)Available commands:$(RESET)"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ {printf "  $(GREEN)%-18s$(RESET) %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""

# Development
dev: db-up ## Start DB, frontend and backend
	@echo "Starting frontend and backend..."
	bun run dev:all

dev-frontend: ## Start frontend only (localhost:3000)
	bun run dev:frontend

dev-backend: db-up ## Start DB and backend only (localhost:8080)
	cd backend && go run ./cmd/server

db-up: ## Start PostgreSQL database
	@docker compose -f docker-compose.dev.yml up -d --wait

db-down: ## Stop PostgreSQL database
	@docker compose -f docker-compose.dev.yml down

db-reset: ## Reset database (delete all data)
	@docker compose -f docker-compose.dev.yml down -v
	@docker compose -f docker-compose.dev.yml up -d --wait

# Build
build: build-frontend build-backend ## Build frontend and backend

build-frontend: ## Build frontend for production
	bun run build:frontend

build-backend: ## Build backend binary
	cd backend && go build -o bin/server ./cmd/server

# Test
test: ## Run all tests
	bun run test

test-backend: ## Run backend tests
	cd backend && go test ./...

test-frontend: ## Run frontend tests
	bun run --cwd frontend test

# Lint & Format
lint: ## Run Biome + ESLint
	cd frontend && bun run lint

lint-fix: ## Fix linting issues automatically
	cd frontend && bun run lint:fix

format: ## Format code with Biome
	cd frontend && bun run format

typecheck: ## Run TypeScript check with JSON report
	cd frontend && bun run typecheck

# Install dependencies
install: ## Install all dependencies
	cd frontend && bun install
	cd backend && go mod tidy

# Clean
clean: ## Remove build artifacts and dependencies
	rm -rf frontend/.next frontend/node_modules
	rm -rf backend/bin
	rm -rf node_modules

# Docker Production
docker-build: ## Build production Docker images
	docker compose build

docker-up: ## Start production containers
	docker compose up -d

docker-down: ## Stop production containers
	docker compose down

# Goca
goca-feature: ## Generate a new feature with Goca
	@read -p "Feature name: " name; \
	read -p "Fields (e.g. name:string,email:string): " fields; \
	cd backend && ~/bin/goca feature $$name --fields "$$fields"
