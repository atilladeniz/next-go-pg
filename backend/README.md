# Backend

[![Backend CI](https://github.com/atilladeniz/next-go-pg/actions/workflows/ci.yml/badge.svg)](https://github.com/atilladeniz/next-go-pg/actions/workflows/ci.yml)

Go Backend with Clean Architecture, generated with [Goca CLI](https://github.com/sazardev/goca).

> **Tip:** LLM-friendly documentation for GORM, Gorilla Mux, Goca, etc. can be found in `../.docs/`

## Tech Stack

- **Go 1.23** - Programming Language
- **Gorilla Mux** - HTTP Router
- **GORM** - ORM for PostgreSQL
- **zerolog** - Structured JSON Logging
- **Swagger/swag** - API Documentation
- **Goca CLI** - Code Generator for Clean Architecture

## Architecture

The backend is split into four **bounded contexts** (DDD-strategic), each owning its own Clean-Architecture stack (DDD-tactical). Cross-cutting infrastructure lives in `platform/`; the dependency graph is built in `composition/`.

```
internal/
├── shared/domain/   # Shared Kernel: UserID, AggregateBase, DomainEvent interface
├── stats/           # Bounded Context: per-user counters
├── auth/            # Bounded Context: identity (Better Auth read-only)
├── notifications/   # Bounded Context: transactional email
├── exports/         # Bounded Context: CSV/JSON data export
├── platform/        # Cross-cutting: middleware, SSE broker
└── composition/     # Composition root + Anti-Corruption Layers
```

Inside each context:

```
<ctx>/
├── domain/                       # Pure entities, value objects, aggregate roots, domain events
├── application/                  # Ports (repositories, publishers, ...) + use-case structs
├── infrastructure/
│   ├── persistence/              # GORM model + mapper + repo impl + Entities()
│   └── ...                       # SSE adapter, River workers, SMTP sender, etc.
└── interfaces/http/              # HTTP handlers — depend ONLY on this context's application/
```

| Layer | Description | Goca Command (then relocate into a context) |
|-------|-------------|--------------|
| Domain | Aggregate roots, value objects, domain events | `goca make entity` |
| Application | Use cases + ports (interfaces) | `goca make usecase` |
| Infrastructure | GORM repo, SMTP sender, River worker, ... | `goca make repository` |
| Interfaces (HTTP) | HTTP endpoints | `goca make handler` |

Cross-context references are forbidden. The composition root is the only place that knows about every context, the database, and River — and is also where Anti-Corruption Layers (e.g. `statsToExportsReader`, `authToNotificationsDirectory`) translate between contexts.

## Quick Start

### 1. Install dependencies:
```bash
go mod tidy
```


### 2. Configure database (PostgreSQL):

#### Option A: Using Docker (Recommended)
```bash
# Run PostgreSQL
docker run --name postgres-dev \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=backend \
  -p 5432:5432 \
  -d postgres:15

# Or using docker-compose
docker-compose up -d postgres
```


#### Option B: Local PostgreSQL
```bash
# Create database
createdb backend
```


### 3. Configure environment variables:
```bash
# Copy example file
cp .env.example .env

# Edit with your credentials
# DB_PASSWORD=password
# DB_NAME=backend
```


### 4. Run the application:
```bash
go run cmd/server/main.go
```


### 5. Test endpoints:
```bash
# Health check
curl http://localhost:8080/health

# Create user (if you have the User feature)
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}'
```


## Project Structure

```
backend/
├── cmd/
│   ├── server/           # Application entry point (composition.Build → ListenAndServe)
│   ├── migrate/          # golang-migrate CLI (prod SQL migrations)
│   └── river-migrate/    # River job-queue migration CLI
├── internal/
│   ├── shared/domain/                    # Shared Kernel (UserID, AggregateBase, DomainEvent)
│   ├── stats/                            # Bounded Context: per-user counters
│   │   ├── domain/                       # UserStats aggregate, StatField VO, events
│   │   ├── application/                  # Ports + use cases (Execute(ctx, ...))
│   │   ├── infrastructure/
│   │   │   ├── persistence/              # GORM model + mapper + repo + Entities()
│   │   │   └── events/                   # Domain-event → SSE publisher
│   │   └── interfaces/http/              # /stats endpoints
│   ├── auth/                             # Bounded Context: identity (Better Auth)
│   │   ├── domain/                       # User projection
│   │   ├── application/                  # UserDirectory port
│   │   ├── infrastructure/betterauth/    # GORM adapter over Better Auth tables
│   │   └── interfaces/http/              # /me, /hello, /protected/hello
│   ├── notifications/                    # Bounded Context: transactional email
│   │   ├── application/                  # EmailSender, JobEnqueuer, UserDirectory ports
│   │   ├── infrastructure/
│   │   │   ├── email/                    # gomail SMTP sender
│   │   │   └── jobs/                     # River email workers + enqueuer
│   │   └── interfaces/http/              # /webhooks/*
│   ├── exports/                          # Bounded Context: data export
│   │   ├── domain/                       # Format, Status VOs
│   │   ├── application/                  # Store, ProgressPublisher, JobEnqueuer, StatsReader
│   │   ├── infrastructure/
│   │   │   └── jobs/                     # River export worker + enqueuer
│   │   └── interfaces/http/              # /export/*
│   ├── platform/                         # Cross-cutting infrastructure
│   │   ├── middleware/                   # Auth, CORS, logging, rate-limit, metrics
│   │   └── sse/                          # SSE broker
│   └── composition/                      # Composition root + Anti-Corruption Layers
├── pkg/
│   ├── config/           # Application configuration
│   ├── logger/           # zerolog Logger (structured JSON)
│   └── river/            # River client wrapper
├── migrations/           # SQL migrations (prod only — empty in dev)
├── docs/                 # Swagger documentation (generated)
├── .goca.yaml            # Goca configuration
├── .env                  # Environment variables
├── .env.example          # Configuration example
└── go.mod
```


## Goca Commands

### Generate New Feature

```bash
# Complete feature with all layers
goca feature Product --fields "name:string,price:float64,stock:int"

# Feature with validation
goca feature Order --fields "userId:string,total:float64" --validation

# Integrate all features (Routes, DI)
goca integrate --all
```

### Generate Individual Layers

```bash
# Entity only (Domain Layer)
goca make entity Product

# Repository only (Data Layer)
goca make repository Product

# UseCase only (Business Logic)
goca make usecase Product

# Handler only (HTTP Layer)
goca make handler Product
```

### Entity Registry (AutoMigrate)

There is **no central registry**. Each bounded context that owns persistence exposes its own `Entities()` function from `internal/<ctx>/infrastructure/persistence/registry.go`. The composition root aggregates them in `runAutoMigrations`:

```go
// internal/<ctx>/infrastructure/persistence/registry.go
func Entities() []any {
    return []any{&gormProduct{}}  // unexported GORM-tagged twin of the domain type
}

// internal/composition/composition.go (runAutoMigrations)
entities := []any{}
entities = append(entities, statspersist.Entities()...)
entities = append(entities, productspersist.Entities()...)  // ← new context
```

`cmd/server/main.go` remains unchanged — it just calls `composition.Build`.

### After Goca/API Changes

```bash
# Swagger + Orval in one command (from root directory)
cd ..
make api

# This automatically runs:
# 1. swag init → backend/docs/swagger.json
# 2. orval → frontend/src/api/endpoints/
```

## Development Commands

```bash
# Run application
make run

# Run tests
make test

# Build for production
make build

# Linting and formatting
make lint
make fmt

# Generate Swagger
make swagger
```


## Troubleshooting

### Error: "dial tcp [::1]:5432: connection refused"
PostgreSQL database is not running.

**Solution:**
```bash
# With Docker
docker run --name postgres-dev \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=backend \
  -p 5432:5432 \
  -d postgres:15

# Verify it's running
docker ps
```


### Error: "database not configured"
Database environment variables are not configured.

**Solution:**
```bash
# Configure in .env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=backend
```


### Error: "command not found: goca"
Goca CLI is not installed or not in PATH.

**Solution:**
```bash
# Reinstall Goca
go install github.com/sazardev/goca@latest

# Verify installation
goca version
```


### Health Check shows "degraded"
Application runs but cannot connect to database.

**Solution:**
1. Verify PostgreSQL is running
2. Verify environment variables in .env
3. Test connection manually: `psql -h localhost -U postgres -d backend`

## Logging

Structured JSON logging with zerolog:

```go
import "github.com/atilladeniz/next-go-pg/backend/pkg/logger"

// Simple logging
logger.Info().Msg("Server started")
logger.Error().Err(err).Msg("Database failed")

// Structured logging
logger.Info().
    Str("user_id", "123").
    Str("action", "login").
    Msg("User logged in")

// With request context (includes request_id)
logger.WithContext(ctx).Info().Msg("Request processed")
```

### Features

- **Request ID Tracing**: Every HTTP request gets `X-Request-ID`
- **Development**: Pretty colored output
- **Production**: JSON for log aggregation (ELK, Datadog, etc.)
- **Log Levels**: debug, info, warn, error, fatal

### Configuration

```bash
LOG_LEVEL=info          # debug, info, warn, error
ENVIRONMENT=production  # development = pretty output
```

## Authentication Webhooks

The backend handles email sending for Better Auth (Magic Link authentication):

```
POST /api/v1/webhooks/send-magic-link       # Send Magic Link email
POST /api/v1/webhooks/send-verification-email  # Send email verification
POST /api/v1/webhooks/session-created       # Login notification (new device)
```

All webhooks are protected by `X-Webhook-Secret` header.

### Webhook Handler

Located at `internal/notifications/interfaces/http/handler.go`:

- **SendMagicLink**: Sends Magic Link emails via SMTP
- **SendVerificationEmail**: Sends email verification links
- **SessionCreated**: Sends login notification only for NEW devices
  - Checks if device/IP combination was seen before
  - Prevents notification spam for known devices

### Email Configuration

```bash
# .env
SMTP_HOST=127.0.0.1
SMTP_PORT=1025
SMTP_FROM=noreply@localhost
WEBHOOK_SECRET=<shared-secret-with-frontend>
NEXT_PUBLIC_APP_URL=http://localhost:3000
```

For local development, use Mailpit (included in `make dev`):
- SMTP: localhost:1025
- Web UI: http://localhost:8025

## Additional Resources

- [Goca Documentation](https://github.com/sazardev/goca)
- [Clean Architecture Principles](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Complete Tutorial](https://github.com/sazardev/goca/wiki/Complete-Tutorial)

## Contributing

This project was generated with Goca. To contribute:

1. Add new features with `goca feature`
2. Maintain layer separation
3. Write tests for new functionality
4. Follow Clean Architecture conventions

---

Generated with [Goca](https://github.com/sazardev/goca)
