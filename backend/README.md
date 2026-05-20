# Backend

[![Backend CI](https://github.com/atilladeniz/next-go-pg/actions/workflows/ci.yml/badge.svg)](https://github.com/atilladeniz/next-go-pg/actions/workflows/ci.yml)

Go Backend organised into **bounded contexts** with **Clean Architecture (DDD)** layers per context.

> **Tip:** LLM-friendly documentation for GORM, Gorilla Mux, etc. can be found in `../.docs/`

## Tech Stack

- **Go 1.23** - Programming Language
- **Gorilla Mux** - HTTP Router
- **GORM** - ORM for PostgreSQL
- **zerolog** - Structured JSON Logging
- **Swagger/swag** - API Documentation

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

| Layer | Description | Allowed imports |
|-------|-------------|------------------|
| Domain | Aggregate roots, value objects, domain events | `shared/domain` only |
| Application | Use cases (`Execute(ctx, ...)`) + ports (interfaces) | this context's `domain/` + `shared/domain` |
| Infrastructure | GORM repo, SMTP sender, River worker, SSE adapter, ... | this context's `application/` + `domain/` + external libs |
| Interfaces (HTTP) | HTTP endpoints with Swagger annotations | this context's `application/` only (no GORM) |

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
├── .env                  # Environment variables
├── .env.example          # Configuration example
└── go.mod
```


## Adding a New Aggregate

The fastest template is the existing `internal/stats/` bounded context — it ships every DDD pattern in ~200 lines: aggregate root with `AggregateBase`, value-object constructor (`StatField`), domain event (`StatIncremented`), repository port, GORM twin + mapper, ACL-friendly use case. Copy it as a starting point, rename, and remove the parts you don't need.

Workflow:

1. **Decide on a bounded context.** A new aggregate joins an existing context if it shares vocabulary and consistency rules; otherwise create a new context folder under `internal/<ctx>/` with the four layer subfolders.
2. **Domain** (`internal/<ctx>/domain/`): pure types only. No `gorm.io/gorm` imports, no I/O. Embed `shared.AggregateBase` if it raises events. Define value-object constructors with invariants (`NewMoney`, `NewSKU`, ...). Define events implementing `EventName() string`.
3. **Application** (`internal/<ctx>/application/`): `ports.go` declares interfaces (`Repository`, `JobEnqueuer`, ...). `<aggregate>_usecases.go` holds the use-case structs whose `Execute(ctx, ...)` orchestrates the aggregate. **Pull events with `agg.PullEvents()` BEFORE `repo.Save(...)`** so a buggy repository can't drop them.
4. **Infrastructure** (`internal/<ctx>/infrastructure/persistence/`): unexported GORM-tagged twin (`gorm<Aggregate>`), mapper (`toDomain` / `fromDomain`), repo impl with the port assertion `var _ <ctx>app.Repository = (*Repository)(nil)`, and `Entities() []any` for AutoMigrate. **Save must not replace `*agg` whole** — mutate only DB-owned fields so pending events survive.
5. **Interfaces** (`internal/<ctx>/interfaces/http/handler.go`): imports only this context's `application/` package. Add Swagger annotations on every endpoint.
6. **Wire** in `internal/composition/composition.go`: build repo → use cases → handler, register routes, append `<ctx>persist.Entities()` to `runAutoMigrations`. If the context needs data from another context, add an Anti-Corruption Layer adapter right here (mirror `statsToExportsReader` / `authToNotificationsDirectory`).
7. **Regenerate API**: `cd .. && just api`.

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

### Regenerate API client

```bash
# Swagger + Orval in one command (from repo root)
cd ..
just api

# This automatically runs:
# 1. swag init → backend/docs/swagger.json
# 2. orval     → frontend/src/shared/api/endpoints/
```

## Development Commands

```bash
just dev-backend   # Start backend (db must be up: just db-up)
just test-backend  # Run all backend tests
just build-backend # Build production binary
just lint          # Lint frontend (backend lint runs in CI)
just swagger       # Regenerate Swagger only
just api           # Regenerate Swagger + Orval client
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

For local development, use Mailpit (included in `just dev`):
- SMTP: localhost:1025
- Web UI: http://localhost:8025

## Additional Resources

- [Clean Architecture Principles](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design (Vernon, *Implementing DDD*)](https://www.informit.com/store/implementing-domain-driven-design-9780321834577) — bounded contexts, aggregates, ACL
- [Hexagonal Architecture (Alistair Cockburn)](https://alistair.cockburn.us/hexagonal-architecture/)

## Contributing

1. Pick a bounded context (`internal/<ctx>/`) or create a new one.
2. Add new aggregates following the "Adding a New Aggregate" workflow above. The `internal/stats/` context is the canonical small example.
3. Maintain layer separation: domain pure, handlers depend only on application ports, cross-context wiring through `composition/` with an ACL.
4. Write tests for new functionality (see `internal/stats/application/usecases_test.go` for the use-case test style — fake repository, fake publisher, regression test for the events-survive-Save invariant).
