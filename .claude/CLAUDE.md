# Next-Go-PG - Project Context

@../AGENTS.md

## IMPORTANT: Technical Docs

For **Next.js**: read `frontend/node_modules/next/dist/docs/` — version-matched, always current with the installed `next` package. The repo-root `AGENTS.md` enforces this rule.

For everything else: **check `.docs/` first** before searching the internet:

```
.docs/
├── tanstack-query.md   # TanStack Query / React Query
├── better-auth.md      # Better Auth
├── kamal-deploy.md     # Kamal Deployment (Docker)
├── logging.md          # Logging (zerolog + Pino)
├── river.md            # River Job Queue
├── background-jobs.md  # Background Job Integration Guide
├── disaster-recovery.md # Backups & restore (postgres-backup-s3 + RustFS)
├── rustfs.md           # RustFS (S3-compatible storage)
├── openspec.md         # OpenSpec slash-command cheatsheet
└── fsd-liniting.xml    # FSD lint rules (Steiger)
```

LLM-friendly documentation for the rest of the tech stack lives there.

---

## Project Overview

Full-Stack Monorepo with Next.js 16 Frontend and Go Backend, PostgreSQL database and Better Auth for authentication.

## Tech Stack

### Frontend (`/frontend`)

- **Framework**: Next.js 16 with App Router and Turbopack
- **Architecture**: Feature-Sliced Design (FSD)
- **Language**: TypeScript 5.9
- **Styling**: Tailwind CSS 4 + shadcn/ui (neutral theme)
- **State**: TanStack Query (React Query)
- **Auth Client**: Better Auth React Client
- **API Client**: Orval-generated hooks from OpenAPI Spec
- **Linting**: Biome + Steiger (FSD Linting)
- **Logging**: Pino (structured JSON)
- **Package Manager**: Bun

### Backend (`/backend`)

- **Language**: Go
- **Framework**: Gorilla Mux Router
- **Architecture**: Bounded Contexts + Clean Architecture (DDD) — see "Backend Architecture" below
- **ORM**: GORM
- **Auth**: Better Auth Session Validation
- **API Docs**: Swagger/swag
- **Logging**: zerolog (structured JSON)
- **Background Jobs**: River (PostgreSQL-native job queue)
- **Module**: `github.com/atilladeniz/next-go-pg/backend`

### Infrastructure

- **Database**: PostgreSQL 16 (Docker)
- **Dev Environment**: Docker Compose for DB
- **Log Aggregation**: Grafana + Loki + Promtail (self-hosted)

---

## IMPORTANT: Backend Architecture (DDD)

The backend used to be scaffolded with Goca. That tool's flat-layer output (`internal/{domain,usecase,repository,handler}`) is structurally incompatible with the bounded-context layout below — `path:` is ignored by the CLI, there is no flag to override output, and porting the generated files takes longer than writing them by hand. **Goca was removed from the toolchain on `refactor/backend-clean-architecture`. Do not reintroduce it. Do not run `goca` against this repo.**

### Adding a new aggregate — the only workflow

The fastest starting point is to **copy `internal/stats/`** — it is the canonical small example and ships every DDD pattern (aggregate root with `shared.AggregateBase`, value-object constructor `NewStatField`, domain event `StatIncremented`, repository port, GORM twin + mapper, ACL-friendly use case, regression test against the events-survive-Save invariant). Rename, strip what you do not need, and follow the checklist in "Bounded Contexts + Clean Architecture Layers".

### Bounded Contexts + Clean Architecture Layers

The backend is split into four **bounded contexts** (DDD-strategic), each owning its own four-layer Clean Architecture stack (DDD-tactical). Cross-cutting infrastructure lives in `platform/`; the only place that wires everything together is `composition/`.

```
backend/internal/
├── shared/domain/                # Shared Kernel: UserID, AggregateBase, DomainEvent interface
├── stats/                        # Bounded Context: per-user counters
│   ├── domain/                   # UserStats aggregate, StatField VO, StatIncremented event
│   ├── application/              # Repository port, DomainEventPublisher, GetUserStats / IncrementStatField use cases
│   ├── infrastructure/
│   │   ├── persistence/          # GORM model + mapper + repository impl + Entities()
│   │   └── events/               # SSE-backed domain event publisher
│   └── interfaces/http/          # /stats endpoints
├── auth/                         # Bounded Context: identity (read-only, owned by Better Auth)
│   ├── domain/                   # User projection
│   ├── application/              # UserDirectory port
│   ├── infrastructure/betterauth # GORM adapter over Better Auth's user / session tables
│   └── interfaces/http/          # /me, /hello, /protected/hello
├── notifications/                # Bounded Context: transactional email
│   ├── application/              # EmailSender, JobEnqueuer, UserDirectory ports + payloads
│   ├── infrastructure/
│   │   ├── email/                # gomail sender adapter
│   │   └── jobs/                 # River workers + enqueuer
│   └── interfaces/http/          # /webhooks/*
├── exports/                      # Bounded Context: user data export
│   ├── domain/                   # Format, Status VOs
│   ├── application/              # Store, ProgressPublisher, JobEnqueuer, StatsReader (ACL)
│   ├── infrastructure/
│   │   ├── (memory store)        # in-memory artifact store
│   │   └── jobs/                 # DataExportWorker + enqueuer
│   └── interfaces/http/          # /export/*
├── platform/                     # Cross-cutting infrastructure
│   ├── middleware/               # Auth, CORS, logging, rate-limit, metrics
│   └── sse/                      # SSE broker
└── composition/                  # Composition root — Build / Shutdown + Anti-Corruption Layers
```

**Layer dependency direction (inward only):**

```
composition → <ctx>/interfaces/http → <ctx>/application → <ctx>/domain
composition → <ctx>/infrastructure  → <ctx>/application → <ctx>/domain
<all>       → shared/domain
```

- Each context's `domain/` imports nothing internal except the Shared Kernel (`shared/domain`).
- Each context's `application/` imports only its own `domain/` and the Shared Kernel.
- Each context's `infrastructure/...` imports only its own `application/` and `domain/`.
- Each context's `interfaces/http/` consumes only its own `application/`. **Never** another context's package, and **never** GORM.
- Contexts never import each other. Cross-context wiring goes through an **Anti-Corruption Layer** in `composition/composition.go` (e.g. `statsToExportsReader`, `authToNotificationsDirectory`).
- `composition/` is the only package that knows about every context, the database, and River.

### Step-by-step: adding a new aggregate

Walk through this using a hypothetical `Invoice` aggregate. The whole thing is hand-written — no code generator. Five small files plus a wiring step in composition.

```text
# 0. Pick a bounded context.
#    A new aggregate joins an existing context if it shares vocabulary
#    AND consistency rules; otherwise create a new top-level folder
#    under internal/ (e.g. internal/billing/).

# 1. Domain (pure): internal/billing/domain/invoice.go
#    - No gorm.io/gorm import. No I/O. No HTTP.
#    - Embed shared.AggregateBase if the aggregate raises events.
#    - Define value objects with constructor invariants
#      (e.g. NewMoney(amount, currency) (Money, error)).
#    - Define domain events (e.g. InvoicePaid) with EventName() string.

# 2. Application port + use case: internal/billing/application/
#    - ports.go declares the interface:
#        type Repository interface {
#            GetByID(ctx, billing.InvoiceID) (*billing.Invoice, error)
#            Save(ctx, *billing.Invoice) error
#        }
#    - invoice_usecases.go holds the use-case struct:
#        type MarkInvoicePaid struct {
#            Repo   Repository
#            Events shared.DomainEventPublisher
#        }
#        func (uc MarkInvoicePaid) Execute(ctx, id billing.InvoiceID) (*billing.Invoice, error) {
#            agg, err := uc.Repo.GetByID(ctx, id)
#            ...
#            agg.MarkPaid()
#            events := agg.PullEvents()  // BEFORE Save
#            if err := uc.Repo.Save(ctx, agg); err != nil { return nil, err }
#            uc.Events.Publish(ctx, events...)
#            return agg, nil
#        }

# 3. Infrastructure (persistence): internal/billing/infrastructure/persistence/
#    - gorm_models.go    → unexported gormInvoice with GORM tags
#    - invoice_mapper.go → invoiceToDomain / invoiceFromDomain
#    - invoice_repo.go   → Repository impl, asserts the port:
#        var _ billingapp.Repository = (*Repository)(nil)
#      Save MUST NOT replace *agg as a whole — that wipes
#      AggregateBase.pendingEvents. Mutate only DB-owned fields back
#      into agg. (See internal/stats/.../repository.go for the pattern.)
#    - registry.go       → func Entities() []any { return []any{&gormInvoice{}} }

# 4. HTTP adapter: internal/billing/interfaces/http/handler.go
#    - Depends ONLY on billingapp.*. Never imports gorm, persistence,
#      or any other bounded context.
#    - Swagger annotations on every endpoint.

# 5. Wire in internal/composition/composition.go:
#      billingRepo := billingpersist.NewRepository(db)
#      markPaidUC := &billingapp.MarkInvoicePaid{Repo: billingRepo, Events: ...}
#      billingHandler := billinghttp.NewHandler(markPaidUC)
#      entities = append(entities, billingpersist.Entities()...)  // in runAutoMigrations
#    - If billing needs data from another context, write an ACL adapter here
#      (mirror statsToExportsReader / authToNotificationsDirectory).
#    - Register routes in the routerDeps section.

# 6. Regenerate the API client
cd .. && just api

# 7. Restart backend (AutoMigrate runs on startup via composition.Build)
just dev-backend
```

Default starting point: **copy `internal/stats/` and modify**. It is the smallest complete example of every pattern above.

### Entity Registry

Each bounded context that owns persistence exports its own `Entities()` function
from `internal/<context>/infrastructure/persistence/`. The composition root
aggregates them in `runAutoMigrations`:

```go
// internal/<context>/infrastructure/persistence/registry.go
func Entities() []any {
    return []any{&gormInvoice{}}
}

// internal/composition/composition.go (runAutoMigrations)
entities := []any{}
entities = append(entities, statspersist.Entities()...)
entities = append(entities, billingpersist.Entities()...)  // ← new context
```

There is **no central registry** — each context owns its own. The composition root is the only place that knows about every context's models.

### API Generation Workflow

`just api` automatically runs:

1. **swag init** → Generates `backend/docs/swagger.json` from Go comments
2. **orval** → Generates TypeScript Hooks in `frontend/src/shared/api/`

```bash
# Run after every API change:
just api

# Or separately:
just swagger     # Only generate Swagger
cd frontend && bunx orval  # Only run Orval
```

### Important for Claude

When you modify backend endpoints:

1. Add Swagger comments to Handler (`// @Summary`, `// @Router`, etc.)
2. Run `just api`
3. Frontend can use the new hooks (`useGetX`, `usePostX`, etc.)

---

## IMPORTANT: Feature-Sliced Design (FSD) in Frontend

### Layer Hierarchy

```
app/        → widgets, features, entities, shared
widgets/    → features, entities, shared
features/   → entities, shared
entities/   → shared
shared/     → (only external libs)
```

### FSD Structure

```
frontend/src/
├── app/                        # Next.js App Router (outside FSD)
│   ├── (auth)/                 # Auth pages (login, register, verify)
│   │   ├── login/
│   │   ├── magic-link/verify/  # Magic Link verification UI
│   │   └── verify-email/       # Email verification UI
│   └── (protected)/            # Protected pages
│       ├── dashboard/
│       └── settings/           # Session management
├── widgets/
│   └── header/                 # Header + UserMenu + ModeToggle
│       ├── ui/
│       └── index.ts
├── features/
│   ├── auth/                   # Magic Link Login, AuthSync
│   │   ├── ui/
│   │   │   ├── login-form.tsx
│   │   │   ├── login-card.tsx
│   │   │   └── email-sent-card.tsx
│   │   ├── model/
│   │   │   ├── use-login.ts
│   │   │   └── use-auth-sync.ts
│   │   └── index.ts
│   ├── user-settings/          # Session Management
│   │   ├── ui/
│   │   ├── model/
│   │   └── index.ts
│   ├── security-settings/      # 2FA Setup, Backup Codes
│   │   ├── ui/
│   │   ├── model/
│   │   └── index.ts
│   ├── data-export/            # CSV/JSON Export with SSE Progress
│   │   ├── ui/
│   │   ├── model/
│   │   └── index.ts
│   └── stats/                  # Stats Grid, SSE
│       ├── ui/
│       ├── model/
│       └── index.ts
├── entities/
│   └── user/                   # User Types, UserInfo
│       ├── ui/
│       ├── model/
│       └── index.ts
└── shared/
    ├── ui/                     # shadcn/ui
    ├── api/                    # Orval-generated
    ├── lib/
    │   ├── auth-client/        # Client-safe Auth (Magic Link)
    │   ├── auth-server/        # Server-only Auth Config
    │   ├── geo/                # User Agent Parsing
    │   ├── go-jwt/             # Go-backend JWT verification
    │   ├── logger/             # Pino Logger
    │   └── ...
    └── config/                 # Providers, Theme
```

### Import Rules (STRICT)

```tsx
// ✅ ALLOWED: Import from lower layer
import { Button } from "@shared/ui/button"
import { User } from "@entities/user"
import { LoginForm } from "@features/auth"

// ❌ FORBIDDEN: Import from higher layer
// In shared/ NEVER import features/!
// In entities/ NEVER import features/!
// In features/ NEVER import other features/!
```

### TypeScript Path Aliases

```tsx
// Always use layer aliases:
import { Button } from "@shared/ui/button"
import { useAuthSync } from "@features/auth"
import { SessionUser } from "@entities/user"
import { Header } from "@widgets/header"
```

### Public API Pattern

Every slice MUST have an `index.ts`:

```tsx
// features/auth/index.ts
export { LoginForm } from "./ui/login-form"
export { useAuthSync, broadcastSignOut, broadcastSignIn } from "./model/use-auth-sync"

// features/user-settings/index.ts
export { SessionsList } from "./ui/sessions-list"
export { SessionCard } from "./ui/session-card"
export { useSessions } from "./model/use-sessions"
```

### Linting

```bash
# FSD rules check (integrated in lint)
bun run lint

# Only Steiger
bunx steiger src
```

### Adding a New Feature

1. Create feature folder: `src/features/<name>/`
2. Add segments: `ui/`, `model/`, `api/` (as needed)
3. Public API: `index.ts` with all exports
4. Import only from `shared/` or `entities/`

---

## IMPORTANT: HydrationBoundary Pattern (TanStack Recommended)

### Architecture Principle

```
Server Component (prefetchQuery) → HydrationBoundary → Client Component (useQuery)
```

### Protected Page Pattern

```tsx
// app/(protected)/dashboard/page.tsx - SERVER COMPONENT
import { dehydrate, HydrationBoundary } from "@tanstack/react-query"
import { cookies } from "next/headers"
import { redirect } from "next/navigation"
import { getStats, getGetStatsQueryKey } from "@shared/api/endpoints/users/users"
import { getQueryClient } from "@shared/lib/query-client"
import { getSession } from "@shared/lib/auth-server"

export default async function DashboardPage() {
  // 1. Check session
  const session = await getSession()
  if (!session) redirect("/login")

  // 2. Get cookies for auth
  const cookieStore = await cookies()
  const cookieHeader = cookieStore.getAll().map((c) => `${c.name}=${c.value}`).join("; ")

  // 3. Prefetch with Orval function
  const queryClient = getQueryClient()
  await queryClient.prefetchQuery({
    queryKey: getGetStatsQueryKey(),
    queryFn: () => getStats({ headers: { Cookie: cookieHeader }, cache: "no-store" }),
  })

  // 4. Wrap with HydrationBoundary
  return (
    <HydrationBoundary state={dehydrate(queryClient)}>
      <Header user={session.user} />
      <Content />
    </HydrationBoundary>
  )
}
```

### Client Component (no initialData needed!)

```tsx
// content.tsx - CLIENT COMPONENT
"use client"

import { useGetStats } from "@shared/api/endpoints/users/users"
import { useSSE } from "@features/stats"

export function Content() {
  useSSE() // Real-time Updates

  // Data is already hydrated!
  const { data } = useGetStats()
  // No loading state needed!
}
```

### What NOT to do

❌ `useSession()` in Protected Pages → Flicker
❌ Skeleton for Session Loading
❌ Client-side Redirect
❌ Manual `fetch()` calls → ALWAYS use Orval!

### What to do INSTEAD

✅ Server Component checks Session with `getSession()`
✅ `redirect()` if no session
✅ `prefetchQuery` with Orval function
✅ `HydrationBoundary` for cache hydration

---

## SSE + React Query Pattern

Real-time updates without polling:

1. **Backend** sends SSE events on changes
2. **Frontend** `useSSE()` hook listens for events
3. **React Query** is automatically invalidated

```tsx
// Backend: SSE Broadcast
h.sseBroker.Broadcast("stats-updated", `{"field":"projects"}`)

// Frontend: Hook
export function useSSE() {
  const queryClient = useQueryClient()

  useEffect(() => {
    const eventSource = new EventSource(`${API_BASE}/api/v1/events`)
    eventSource.addEventListener("stats-updated", () => {
      queryClient.invalidateQueries({ queryKey: getGetStatsQueryKey() })
    })
  }, [])
}
```

---

## Important Files

### API Definition

- `backend/docs/swagger.json` - Generated from Go comments
- `frontend/orval.config.ts` - Orval config for API client generation

### Data Fetching (FSD Paths)

- `frontend/src/shared/lib/query-client.ts` - Shared QueryClient for Server + Client
- `frontend/src/shared/api/custom-fetch.ts` - Fetch wrapper for Orval
- `frontend/src/shared/api/endpoints/` - Orval-generated hooks
- **IMPORTANT**: ALWAYS use Orval functions with `prefetchQuery` + `HydrationBoundary`!

### Authentication (FSD Paths)

- `frontend/src/shared/lib/auth-server/` - Better Auth Server Config + Session Helper
- `frontend/src/shared/lib/auth-client/` - Better Auth Client (Magic Link + Actions)
- `frontend/src/features/auth/` - Magic Link Login + Cross-Tab Sync
- `frontend/src/features/user-settings/` - Session Management (List, Revoke)
- `frontend/src/features/security-settings/` - 2FA / TOTP setup + Backup Codes
- `frontend/src/features/data-export/` - CSV/JSON export with SSE progress
- `backend/internal/notifications/interfaces/http/handler.go` - Email Webhooks (Magic Link, Verification, Login Notification)

### Real-time

- `backend/internal/platform/sse/broker.go` - SSE Broker (platform-wide)
- `backend/internal/stats/infrastructure/events/publisher.go` - Stats domain-event → SSE adapter
- `frontend/src/features/stats/model/use-sse.ts` - SSE Client Hook

### Logging

- `backend/pkg/logger/logger.go` - zerolog Logger with helpers
- `backend/internal/platform/middleware/logging.go` - HTTP request logging + Request-ID
- `frontend/src/shared/lib/logger/index.ts` - Pino Logger

### UI Components (FSD Paths)

- `frontend/src/shared/ui/` - shadcn/ui components
- `frontend/src/widgets/header/` - App Header with Mode Toggle
- `frontend/src/shared/config/` - Providers, Theme Config

## Commands

```bash
# Development
just dev              # Start DB + Frontend + Backend
just dev-frontend     # Frontend only (localhost:3000)
just dev-backend      # Backend only (localhost:8080)

# Database
just db-up            # Start PostgreSQL
just db-down          # Stop PostgreSQL
just db-reset         # Reset database

# API Generation
just api              # Generate TypeScript client from OpenAPI

# Quality
just lint             # Biome + Steiger Linting
just lint-fix         # Auto-fix Lint Errors
just typecheck        # TypeScript Check

# Security
just security-scan    # Scan for secrets/sensitive data
just setup-hooks      # Setup Git Hooks

# Build
just build            # Build Frontend + Backend
just build-frontend   # Next.js Production Build
just build-backend    # Go Binary

# Logging Stack (Grafana + Loki)
just logs-up          # Start logging stack
just logs-down        # Stop logging stack
just logs-open        # Open Grafana (localhost:3001)
just logs-query '...'    # Query logs via CLI

# Production migrations (not used in dev — see Migration Strategy section)
just prod-migrate-up        # Apply pending SQL migrations
just prod-river-migrate-up  # Apply River job-queue migrations

# Spec-driven development (OpenSpec)
just spec-list             # List active changes
just spec-validate         # Validate all changes and specs
```

---

## Migration Strategy: Dev (AutoMigrate) vs Prod (SQL)

**Dev path** (`just dev` → `cmd/server` → `composition.Build`): GORM `AutoMigrate` runs for every entity contributed by each bounded context's `infrastructure/persistence/Entities()` function (currently `statspersist.Entities()`). River's job-queue schema is migrated on startup via `riverPkg.RunMigrations`. **No CLI tool is invoked.**

**Prod path** (`just prod-migrate-up` / `just prod-river-migrate-up`): For clustered deployments where multiple replicas can't safely race on AutoMigrate, use `backend/cmd/migrate` (golang-migrate against SQL files in `backend/migrations/`) and `backend/cmd/river-migrate` as pre-deploy hooks. Currently optional — the single-replica Kamal setup is fine with AutoMigrate. The `backend/migrations/` folder ships empty by default; add SQL files only when AutoMigrate stops being sufficient.

**Don't mix.** If you add a SQL migration, also keep the Go entity in the relevant context's `Entities()` in sync — otherwise AutoMigrate and the SQL files describe two different schemas.

---

## Spec-driven Development (OpenSpec)

Non-trivial features go through OpenSpec before implementation. Slash commands and full workflow are documented in `AGENTS.md` (loaded automatically) and `.docs/openspec.md` (cheatsheet).

**Workflow:**

1. `/opsx:propose "<idea>"` — create change + artifacts (proposal, design, tasks)
2. Review/edit artifacts under `openspec/changes/<name>/`
3. `/opsx:apply` — implement the tasks
4. `/opsx:verify` then `/opsx:archive`

**Rule:** Before touching code for a feature listed under `openspec/changes/`, read its `proposal.md` (and `design.md` if present) first. The spec is the contract.

Prerequisite (one-time): `npm install -g @fission-ai/openspec@latest`

---

## Background Jobs (River)

Email sending is processed asynchronously via River background jobs:

### Architecture

```
Webhook Request → Job Enqueue → PostgreSQL → River Worker → Email Send
```

### Job Types

| Job | Description |
|-----|-------------|
| `send_magic_link` | Magic link login emails |
| `send_verification_email` | Email verification |
| `send_2fa_otp` | 2FA one-time passwords |
| `send_login_notification` | New device login alerts |

### Adding New Jobs

Job workers live inside the bounded context that owns the work.

1. Define job args + worker in `backend/internal/<ctx>/infrastructure/jobs/`
2. Implement worker with `river.WorkerDefaults`
3. Register in that context's `jobs/registry.go` (called from composition)
4. Create enqueue helper in that context's `jobs/enqueue.go`
5. Expose a `JobEnqueuer` port in `<ctx>/application/ports.go` so handlers depend on the interface, not the River client.

### Files

```
backend/internal/notifications/infrastructure/jobs/
├── workers.go    # Email workers (Magic Link, Verification, 2FA, ...)
├── registry.go   # Worker registration
└── enqueue.go    # Enqueuer adapter — implements notifapp.JobEnqueuer

backend/internal/exports/infrastructure/jobs/
├── worker.go     # DataExportWorker
├── registry.go   # Worker registration
└── enqueue.go    # Enqueuer adapter — implements exportsapp.JobEnqueuer

backend/pkg/river/
└── client.go     # River client wrapper

backend/cmd/river-migrate/
└── main.go       # Migration CLI (prod)
```

### Fallback

If River is unavailable, emails are sent synchronously (fallback mode).

See `.docs/background-jobs.md` for full documentation.

---

## Logging

Structured JSON logging is configured for both frontend and backend. See `.docs/logging.md` for full documentation.

### Backend (zerolog)

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

// Business events (audit trails)
logger.BusinessEvent(ctx, "user.created", "user", userID, nil)
```

### Frontend (Pino)

```typescript
import { log, createLogger } from "@shared/lib/logger"

// Simple logging
log.info("Page loaded")
log.error("API failed", { endpoint: "/api/users" })

// Component-specific logger
const authLogger = createLogger("auth")
authLogger.info("Login submitted")
```

### Configuration

```bash
# Backend
LOG_LEVEL=info          # debug, info, warn, error
ENVIRONMENT=production  # development = pretty output

# Frontend
LOG_LEVEL=info
NODE_ENV=production
```

### Request ID Tracing

Every HTTP request gets `X-Request-ID` header for distributed tracing.

### Log Aggregation (Grafana + Loki)

Self-hosted log aggregation stack included:

```bash
just logs-up      # Start Grafana + Loki + Promtail
just logs-open    # Open Grafana at localhost:3001
```

Query logs in Grafana with LogQL:
```logql
{service="next-go-pg-api"}                    # All backend logs
{level="error"}                               # Errors only
{category="auth"}                             # Auth events
{service="next-go-pg-api"} |= "user_id"       # Filter by content
```

See `.docs/logging.md` for full documentation.

---

## Security: Gitleaks Pre-Commit Hook

This project uses **gitleaks** to detect secrets before commit.

### What is blocked?

- API Keys, Tokens, Passwords
- Absolute paths with usernames
- Database URLs with real credentials
- Private Keys

### Important for Claude

- **NEVER** write absolute paths with usernames in code/docs
- **NEVER** hardcode real secrets
- Example URLs always with `localhost` or placeholders
- Config: `.gitleaks.toml`

### On Commit Block

```bash
# Show what was blocked
just security-scan

# Emergency bypass (NOT RECOMMENDED)
git commit --no-verify
```

## Project Structure

```
next-go-pg/
├── backend/
│   ├── cmd/
│   │   ├── server/           # Main API server (loads config → composition.Build → ListenAndServe)
│   │   ├── migrate/          # golang-migrate CLI (prod SQL migrations)
│   │   └── river-migrate/    # River job-queue migration CLI
│   ├── internal/
│   │   ├── shared/domain/             # Shared Kernel (UserID, AggregateBase, DomainEvent)
│   │   ├── stats/                     # Bounded Context (per-user counters)
│   │   │   ├── domain/                # Pure: aggregate, value objects, events
│   │   │   ├── application/           # Ports + use cases
│   │   │   ├── infrastructure/
│   │   │   │   ├── persistence/       # GORM model + mapper + repo + Entities()
│   │   │   │   └── events/            # Domain-event → SSE adapter
│   │   │   └── interfaces/http/       # /stats endpoints
│   │   ├── auth/                      # Bounded Context (Better Auth integration)
│   │   │   ├── domain/                # User projection
│   │   │   ├── application/           # UserDirectory port
│   │   │   ├── infrastructure/betterauth/  # GORM adapter over Better Auth tables
│   │   │   └── interfaces/http/       # /me, /hello, /protected/hello
│   │   ├── notifications/             # Bounded Context (transactional email)
│   │   │   ├── application/           # EmailSender, JobEnqueuer, UserDirectory ports
│   │   │   ├── infrastructure/
│   │   │   │   ├── email/             # gomail SMTP sender
│   │   │   │   └── jobs/              # River email workers + enqueuer
│   │   │   └── interfaces/http/       # /webhooks/*
│   │   ├── exports/                   # Bounded Context (CSV/JSON data export)
│   │   │   ├── domain/                # Format, Status VOs
│   │   │   ├── application/           # Store, ProgressPublisher, JobEnqueuer, StatsReader
│   │   │   ├── infrastructure/
│   │   │   │   ├── (memory store)     # in-memory artifact store
│   │   │   │   └── jobs/              # River export worker + enqueuer
│   │   │   └── interfaces/http/       # /export/*
│   │   ├── platform/                  # Cross-cutting infrastructure
│   │   │   ├── middleware/            # Auth, CORS, Logging, Rate limit, Metrics
│   │   │   └── sse/                   # Server-Sent Events broker
│   │   └── composition/               # Composition root + Anti-Corruption Layers
│   ├── migrations/           # SQL migrations (prod only — empty in dev)
│   └── docs/                 # Swagger JSON (generated)
├── frontend/
│   ├── src/
│   │   ├── app/              # Next.js App Router
│   │   ├── widgets/          # Composite UI (Header)
│   │   ├── features/         # auth, user-settings, security-settings,
│   │   │                     # data-export, stats
│   │   ├── entities/         # Business Objects (User)
│   │   └── shared/           # Reusable (UI, API, Lib)
│   └── orval.config.ts
├── openspec/
│   ├── changes/             # Active proposals (archive/ holds shipped ones)
│   └── specs/               # Canonical capability specs
└── infra/
    ├── docker/                         # Dockerfile + supervisord.conf
    ├── kamal/                          # deploy.yml + hooks + secrets.example
    ├── loki/                           # Loki + Promtail configs
    ├── grafana/                        # Grafana provisioning
    └── compose/
        ├── docker-compose.yml          # Production
        ├── docker-compose.dev.yml      # Dev DB + Mailpit
        ├── docker-compose.logging.yml  # Logging stack
        └── docker-compose.backup.yml   # Backup stack
```

> **Layer ownership at a glance.** Each bounded context owns four layers.
> `<ctx>/domain/` holds pure types — no I/O, no GORM — plus value objects with
> constructor invariants, aggregate roots embedding `shared.AggregateBase`, and
> domain events implementing `EventName() string`. `<ctx>/application/` defines
> repository / publisher / enqueuer ports and use-case structs that orchestrate
> them via `Execute(ctx, ...)`. `<ctx>/infrastructure/persistence/` holds the
> GORM-tagged twin types and mappers; each repository impl asserts the port
> with `var _ <ctx>app.<Port> = (*<Impl>)(nil)`. `<ctx>/interfaces/http/`
> imports only its own `application/` package. `composition/` is the single
> place that builds the dependency graph and wires Anti-Corruption Layers
> between contexts (e.g. `statsToExportsReader`, `authToNotificationsDirectory`).
> The old flat `internal/{domain,repository,handler,usecase,jobs,sse,templates}/`
> packages were removed by the `backend-clean-architecture` change.

## Conventions

### Code Style

- Tabs for indentation (Biome Config)
- Double quotes for strings
- No semicolons (unless required)
- German UI texts

### Git Commits (Conventional Commits)

Commits are validated by **commitlint**. Use this format:

```
<type>(<scope>): <description>
```

| Type | Description | Version Bump |
|------|-------------|--------------|
| `feat` | New feature | Minor |
| `fix` | Bug fix | Patch |
| `docs` | Documentation | - |
| `style` | Formatting | - |
| `refactor` | Code restructuring | - |
| `perf` | Performance | Patch |
| `test` | Tests | - |
| `build` | Build system | - |
| `ci` | CI config | - |
| `chore` | Maintenance | - |

Examples:
- `feat: add user authentication`
- `fix(api): resolve timeout issue`
- `feat!: redesign API` (breaking change → Major)

Rules:
- English commit messages
- **NEVER** add "Generated with Claude Code" or similar tags
- **NEVER** add "Co-Authored-By: Claude" or any AI co-author lines

### API

- OpenAPI 3.0 as source of truth
- Run `just api` after spec changes
- Never manually edit generated files

## Linting Ignores

- `src/shared/api/` - Orval generated
- `src/shared/ui/` - shadcn generated

## Environment Variables

### Frontend (.env.local)

```
DATABASE_URL=postgres://postgres:postgres@localhost:5432/nextgopg
NEXT_PUBLIC_APP_URL=http://localhost:3000
NEXT_PUBLIC_API_URL=http://localhost:8080
BETTER_AUTH_SECRET=<secret>
WEBHOOK_SECRET=<webhook-secret>
```

### Backend

```
DATABASE_URL=postgres://postgres:postgres@localhost:5432/nextgopg
PORT=8080
SMTP_HOST=127.0.0.1
SMTP_PORT=1025
SMTP_FROM=noreply@localhost
WEBHOOK_SECRET=<webhook-secret>
NEXT_PUBLIC_APP_URL=http://localhost:3000
```

---

## Background Jobs: River

This project uses **River** for PostgreSQL-native background job processing (~66k jobs/sec).

### Job Types

| Job | Purpose | Trigger |
|-----|---------|---------|
| `send_magic_link` | Magic Link Login Email | Auth Webhook |
| `send_verification_email` | Email Verification | Auth Webhook |
| `send_2fa_otp` | 2FA Code Email | Auth Webhook |
| `send_login_notification` | Login Alert (new device) | Auth Webhook |
| `data_export` | CSV/JSON Export with Progress | User Request |

### Architecture

```
Handler → Enqueue Job → PostgreSQL → River Worker → Process → SSE Events
```

### Key Files

```
backend/internal/notifications/infrastructure/jobs/
├── registry.go     # Worker registration
├── enqueue.go      # Adapter implementing notifapp.JobEnqueuer
└── workers.go      # Email workers (Magic Link, Verification, 2FA, ...)

backend/internal/exports/infrastructure/jobs/
├── registry.go     # Worker registration
├── enqueue.go      # Adapter implementing exportsapp.JobEnqueuer
└── worker.go       # Data export worker with SSE progress

frontend/src/features/data-export/
├── model/use-export.ts   # SSE listener for export-progress
└── ui/export-card.tsx    # Export UI with progress bar
```

### Adding a New Job

```go
// internal/<ctx>/infrastructure/jobs/...
// 1. Define args struct
type MyJobArgs struct {
    UserID string `json:"userId"`
}

func (MyJobArgs) Kind() string { return "my_job" }

// 2. Create worker
type MyJobWorker struct {
    river.WorkerDefaults[MyJobArgs]
}

func (w *MyJobWorker) Work(ctx context.Context, job *river.Job[MyJobArgs]) error {
    // Process job (depend only on this context's application ports)
    return nil
}

// 3. Register in this context's registry.go (called from composition)
river.AddWorker(workers, &MyJobWorker{})

// 4. Expose <ctx>app.JobEnqueuer port and have enqueue.go implement it.
//    HTTP handlers depend on the port, never on the River client.
```

### Fallback

Email jobs fall back to synchronous sending if River is unavailable.

---

## Authentication: Magic Link

This project uses **Magic Link authentication** (passwordless) via Better Auth.

### Flow

1. User enters email on `/login`
2. Backend sends Magic Link email via webhook
3. User clicks link → `/magic-link/verify?token=...`
4. Token is verified → User is logged in → Redirect to `/dashboard`

### Features

- **Rate Limiting**: 3 magic link requests per minute
- **Email Verification**: New users must verify email first
- **Session Management**: View/revoke sessions at `/settings`
- **Login Notifications**: Email on new device/IP login
- **Cross-Tab Sync**: Logout syncs across all tabs

### Backend Webhooks

All email sending is handled by the Go backend:

```
POST /api/v1/webhooks/send-magic-link      # Magic Link email
POST /api/v1/webhooks/send-verification-email  # Email verification
POST /api/v1/webhooks/session-created      # Login notification (new device only)
```

Protected by `X-Webhook-Secret` header.

### Auth Files

```
frontend/src/shared/lib/auth-server/auth.ts                       # Better Auth config
frontend/src/shared/lib/auth-client/                              # Client (magicLinkClient plugin)
frontend/src/features/auth/                                       # Login UI + hooks
frontend/src/features/user-settings/                              # Session management
backend/internal/auth/                                            # Auth bounded context (UserDirectory port + Better Auth adapter)
backend/internal/notifications/interfaces/http/handler.go         # Email webhooks (consumes notifapp.EmailSender + notifapp.UserDirectory)
```

### Session Management

```tsx
// List all sessions
const { sessions } = useSessions()

// Revoke single session
await authClient.revokeSession({ token })

// Revoke all other sessions
await authClient.revokeOtherSessions()
```
