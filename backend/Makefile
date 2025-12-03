# Makefile for backend
.PHONY: help build run test clean docker-build docker-run deps lint migrate-up migrate-down

# Variables
APP_NAME := backend
DOCKER_IMAGE := backend:latest
MIGRATE_PATH := ./migrations
DATABASE_URL := postgres://postgres:@localhost/backend?sslmode=disable

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

deps: ## Install dependencies
	go mod download
	go mod tidy

build: ## Build the application
	go build -o bin/$(APP_NAME) cmd/server/main.go

run: ## Run the application
	go run cmd/server/main.go

test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint: ## Run linter
	golangci-lint run

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html

# Database migrations
migrate-install: ## Install migrate tool
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate-up: ## Apply database migrations
	migrate -path $(MIGRATE_PATH) -database "$(DATABASE_URL)" up

migrate-down: ## Rollback last migration
	migrate -path $(MIGRATE_PATH) -database "$(DATABASE_URL)" down 1

migrate-force: ## Force migration version (use with caution)
	migrate -path $(MIGRATE_PATH) -database "$(DATABASE_URL)" force $(VERSION)

# Docker commands
docker-build: ## Build Docker image
	docker build -t $(DOCKER_IMAGE) .

docker-run: ## Run application in Docker
	docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE)

docker-compose-up: ## Start services with Docker Compose
	docker-compose up -d

docker-compose-down: ## Stop services with Docker Compose
	docker-compose down

# Development helpers
dev-db: ## Start development database
	docker run --name backend-postgres -e POSTGRES_DB=backend -e POSTGRES_PASSWORD=password -p 5432:5432 -d postgres:15

dev-db-stop: ## Stop development database
	docker stop backend-postgres && docker rm backend-postgres

fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

mod-upgrade: ## Upgrade dependencies
	go get -u ./...
	go mod tidy

# Production helpers
build-prod: ## Build for production
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o bin/$(APP_NAME) cmd/server/main.go

# Security
sec-scan: ## Run security scan
	gosec ./...

# API documentation
api-docs: ## Generate API documentation
	swag init -g cmd/server/main.go
