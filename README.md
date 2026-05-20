# Next-Go-PG

[![CI](https://github.com/atilladeniz/next-go-pg/actions/workflows/ci.yml/badge.svg)](https://github.com/atilladeniz/next-go-pg/actions/workflows/ci.yml)

Full-Stack Monorepo with Next.js Frontend and Go Backend.

## Technical Documentation

**IMPORTANT:** LLM-friendly documentation for the entire tech stack can be found in `.docs/`:

```
.docs/
├── tanstack-query.md   # TanStack Query
├── better-auth.md      # Better Auth
├── river.md            # River Job Queue
├── background-jobs.md  # Background Job Integration
├── kamal-deploy.md     # Kamal Deployment
├── logging.md          # Logging (zerolog + Pino)
├── disaster-recovery.md # Database Backups & Recovery
├── rustfs.md           # RustFS (S3-compatible storage)
├── openspec.md         # OpenSpec spec-driven workflow
└── fsd-liniting.xml    # FSD lint rules (Steiger)
```

> **Always check `.docs/` first** before searching the internet!
>
> **Next.js** is the exception — read `frontend/node_modules/next/dist/docs/` instead. Version-matched, always current with the installed `next` package (rule enforced by [AGENTS.md](./AGENTS.md)).

## Tech Stack

| Component | Technology |
|-----------|------------|
| Frontend | Next.js 16, TypeScript, Tailwind CSS, shadcn/ui |
| Frontend Architecture | **Feature-Sliced Design (FSD)** |
| Backend | Go, Gorilla Mux, **Bounded Contexts + Clean Architecture (DDD)**, GORM |
| Database | PostgreSQL 16 |
| Auth | Better Auth (Magic Link, JWT) |
| Background Jobs | **River** (PostgreSQL-native, ~66k jobs/sec) |
| Real-time | Server-Sent Events (SSE) |
| API | Swagger/swag → Orval |
| Logging | zerolog (Go) + Pino (Next.js) |
| Log Aggregation | Grafana + Loki + Promtail |
| Database Backups | postgres-backup-s3 + RustFS (S3-compatible) |
| Linting | Biome + Steiger (FSD) |

## Prerequisites

- [Bun](https://bun.sh/) (Frontend)
- [Go 1.26+](https://go.dev/) (Backend)
- [Docker](https://www.docker.com/) (Database)
- [OpenSpec](https://github.com/Fission-AI/OpenSpec) (Spec-driven Workflow, optional)

### Install CLI Tools

```bash
# Gitleaks (Security Scanning)
brew install gitleaks

# Sitefetch (Documentation Fetching)
bun install -g sitefetch

# OpenSpec (Spec-driven feature workflow, optional but recommended)
npm install -g @fission-ai/openspec@latest
```

## Quick Start

```bash
# Clone repository
git clone <repo-url>
cd next-go-pg

# Install dependencies (also installs gitleaks, sitefetch + sets up git hooks)
just install

# Start database
just db-up

# Create Better Auth tables
just db-migrate

# Start development
just dev
```

Open:
- Frontend: [http://localhost:3000](http://localhost:3000)
- Backend: [http://localhost:8080](http://localhost:8080)

### Propose Your First Feature

This template ships with [OpenSpec](https://github.com/Fission-AI/OpenSpec) — a spec-driven workflow for AI-agent development. Before implementing non-trivial features, write a proposal so Claude (or any other agent) executes against a contract instead of assumptions:

```bash
# In Claude Code:
/opsx:propose "add billing portal"

# Review the generated proposal/design/tasks under openspec/changes/
# Then:
/opsx:apply       # Implement
/opsx:verify      # Sanity check
/opsx:archive     # Mark complete
```

See `.docs/openspec.md` for the full cheatsheet and `AGENTS.md` for the workflow rules.

## Project Structure

```
next-go-pg/
├── backend/                 # Go Backend (Bounded Contexts + Clean Architecture)
│   ├── cmd/
│   │   ├── server/          # Main API server (loads config → composition.Build)
│   │   ├── migrate/         # golang-migrate CLI (prod SQL migrations)
│   │   └── river-migrate/   # River job-queue migration CLI
│   ├── internal/
│   │   ├── shared/domain/   # Shared Kernel (UserID, AggregateBase, DomainEvent)
│   │   ├── stats/           # BC: per-user counters (domain, application, infra, http)
│   │   ├── auth/            # BC: identity (Better Auth read-only adapter)
│   │   ├── notifications/   # BC: transactional email (webhooks + River workers)
│   │   ├── exports/         # BC: CSV/JSON data export with SSE progress
│   │   ├── platform/        # Cross-cutting: middleware (Auth/CORS/Log) + SSE broker
│   │   └── composition/     # Composition root + Anti-Corruption Layers
│   ├── migrations/          # SQL migrations (prod only — empty in dev)
│   ├── pkg/logger/          # zerolog Logger
│   └── docs/                # Swagger (generated)
├── frontend/                # Next.js Frontend (FSD Architecture)
│   ├── src/
│   │   ├── app/             # Next.js App Router
│   │   ├── widgets/         # Composite UI (Header)
│   │   ├── features/        # auth, user-settings, security-settings,
│   │   │                    # data-export, stats
│   │   ├── entities/        # Business Objects (User)
│   │   └── shared/          # Reusable (UI, API, Lib, Logger)
│   └── orval.config.ts      # API Generator Config
├── infra/
│   ├── docker/                 # Dockerfile + supervisord.conf
│   ├── compose/                # All docker-compose files
│   │   ├── docker-compose.yml          # Production
│   │   ├── docker-compose.dev.yml      # Dev DB + Mailpit
│   │   ├── docker-compose.logging.yml  # Logging (Grafana + Loki)
│   │   └── docker-compose.backup.yml   # Backup (postgres-backup-s3 + RustFS)
│   ├── kamal/                  # deploy.yml + hooks + secrets.example
│   ├── loki/                   # Loki & Promtail Config
│   └── grafana/                # Grafana Provisioning
├── openspec/                # Spec-driven change proposals (OpenSpec)
│   ├── changes/             # Active and archived changes
│   └── specs/               # Canonical capability specs
├── AGENTS.md                # Agent rules (Next.js docs, OpenSpec workflow)
├── justfile                 # Build Commands (https://github.com/casey/just)
└── README.md
```

## CI/CD Pipeline

This project uses GitHub Actions for continuous integration and automatic releases.

### Automated Checks (on every PR)

| Check | Tool | Description |
|-------|------|-------------|
| Backend Lint | golangci-lint | Go code quality |
| Backend Security | gosec | Security vulnerabilities |
| Backend Test | go test | Unit tests with race detection |
| Frontend Lint | Biome + Steiger | Code style + FSD architecture |
| Frontend Typecheck | TypeScript | Type safety |
| Dependency Review | GitHub | CVE scanning |
| Secret Scan | Gitleaks | Prevent leaked secrets |

### Automatic Releases (Release Please)

Releases are automated based on [Conventional Commits](https://www.conventionalcommits.org/):

| Commit Type | Version Bump | Example |
|-------------|--------------|---------|
| `fix:` | Patch (0.0.X) | `fix: resolve login bug` |
| `feat:` | Minor (0.X.0) | `feat: add dark mode` |
| `feat!:` | Major (X.0.0) | `feat!: redesign API` |

When you merge to `main`:
1. Release Please creates a "Release PR" with updated CHANGELOG
2. You review and merge the Release PR
3. GitHub Release is created automatically
4. Go binaries + Docker images are built

### Commit Message Format

Commits are validated by **commitlint**:

```bash
# Format
<type>(<scope>): <description>

# Examples
feat: add user authentication
fix(api): resolve timeout issue
docs: update README
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`

## Security

This project uses [gitleaks](https://github.com/gitleaks/gitleaks) to prevent committing secrets and sensitive data.

### Setup (automatic with `just install`)

```bash
# Install gitleaks
brew install gitleaks    # macOS
# or
go install github.com/gitleaks/gitleaks/v8@latest

# Setup git hooks (runs automatically with just install)
just setup-hooks
```

### What it scans for

- API keys, tokens, and passwords
- Absolute paths with usernames
- Database URLs with embedded credentials
- Private keys and certificates
- AWS/GCP/Azure credentials

### Manual scan

```bash
just security-scan    # Scan entire codebase
```

### Bypass (not recommended)

```bash
git commit --no-verify
```

## Just Commands

This project uses [just](https://github.com/casey/just) as command runner. Install with `brew install just`. Run `just` (no args) to see all grouped recipes.

### Setup

```bash
just install          # Install all deps + CLI tools + git hooks
just install-tools    # gitleaks, sitefetch (idempotent)
just clean            # Remove build artifacts and node_modules
```

### Development

```bash
just dev              # Start DB + Frontend + Backend
just dev-full         # Same + Grafana / Loki / Promtail
just dev-frontend     # Frontend only (localhost:3000)
just dev-backend      # Backend only (localhost:8080)
```

### Database

```bash
just db-up            # Start PostgreSQL
just db-down          # Stop PostgreSQL
just db-reset         # Reset database (prompts for confirmation)
just db-migrate       # Run Better Auth migrations
```

### Production Migrations (prod deploy hook — NOT used in dev)

```bash
just prod-migrate-up         # Apply pending SQL migrations
just prod-migrate-down       # Roll back last migration
just prod-migrate-version    # Show current version
just prod-migrate-create <name>  # Scaffold new .up/.down files

just prod-river-migrate-up         # Apply River job-queue migrations
just prod-river-migrate-down       # Roll back last River migration
just prod-river-migrate-version    # Show River migration version
```

> Dev path runs GORM AutoMigrate on `cmd/server` startup. Use these
> `prod-*` recipes only when clustered deployments make AutoMigrate
> racing unsafe. See `backend/migrations/README.md`.

### API

```bash
just api              # Generate TypeScript client from OpenAPI
just swagger          # Only regenerate Swagger JSON
```

### Quality

```bash
just lint             # Biome + Steiger (FSD) Linting
just lint-fix         # Auto-fix
just format           # Format with Biome
just typecheck        # TypeScript Check
just test             # Run all tests
just test-frontend    # Frontend tests only
just test-backend     # Backend tests only
just security-scan    # Scan working tree for secrets
just security-scan-history  # Scan entire git history
```

### Build

```bash
just build            # Frontend + Backend
just build-frontend   # Next.js Production Build
just build-backend    # Go Binary
```

### Spec-driven Development (OpenSpec)

```bash
just spec-list        # List active changes
just spec-specs       # List canonical specs
just spec-validate    # Validate all changes + specs
just spec-status <name>   # Status of a specific change
just spec-view        # Interactive dashboard
```

### Deployment (Kamal)

```bash
just deploy-staging                 # Deploy to staging
just deploy-production              # Deploy to production (confirms)
just deploy-rollback <dest>         # Roll back staging|production
just deploy-logs <dest>             # Follow remote app logs
just deploy-console <dest>          # Open shell on remote
just deploy-setup <dest>            # Initial Kamal setup
```

### Monitoring

```bash
just metrics          # Open Prometheus metrics (localhost:8080/metrics)
```

### Documentation

```bash
just search-docs "query"           # Search docs with semantic search
just search-docs "query" 10        # Custom result count
just search-docs-index             # Pre-build search index (one-time)
just fetch-docs <url>              # Fetch LLM-friendly docs
just fetch-docs <url> <name>       # With custom filename
```

### Logging Stack

```bash
just logs-up      # Start Grafana + Loki + Promtail
just logs-down    # Stop logging stack
just logs-open    # Open Grafana (localhost:3001)
just logs-query '{level="error"}'  # Query logs via CLI
```

### Database Backups

Fully automatic PostgreSQL backups to S3-compatible storage (RustFS).

```bash
just backup-up       # Start automatic backup system
just backup-down     # Stop backup stack
just backup-now      # Create backup immediately
just backup-list     # List all backups in S3
just backup-restore  # Restore from latest backup (prompts for confirmation)
```

- **Schedule**: Daily (configurable via `BACKUP_SCHEDULE`)
- **Retention**: 7 days (configurable via `BACKUP_KEEP_DAYS`)
- **Storage**: RustFS Console at http://localhost:9001

See [Disaster Recovery](.docs/disaster-recovery.md) for details.

## API Workflow

1. Add Swagger comments to Go handler
2. Generate TypeScript client: `just api`
3. Use generated hooks:

```tsx
import { useGetStats } from "@shared/api/endpoints/users/users"

function MyComponent() {
  const { data, isLoading } = useGetStats()
  // ...
}
```

## Background Jobs (River)

PostgreSQL-native job queue with ~66k jobs/sec throughput.

### Job Types

| Job | Description | Trigger |
|-----|-------------|---------|
| `send_magic_link` | Magic Link email | Login request |
| `send_verification_email` | Email verification | New user |
| `send_2fa_otp` | 2FA code | 2FA enabled |
| `send_login_notification` | Login alert | New device/IP |
| `data_export` | CSV/JSON export | User request |

### Data Export Feature

Export user data with real-time progress via SSE:

```tsx
// Frontend: Start export
const { mutate: startExport } = usePostExportStart()
startExport({ data: { format: "csv", dataType: "all" } })

// Listen for progress via SSE
useEffect(() => {
  const es = new EventSource("/api/v1/events")
  es.addEventListener("export-progress", (e) => {
    const progress = JSON.parse(e.data)
    // { jobId, status, progress: 0-100, downloadId }
  })
}, [])
```

### Architecture

```
User Request → Handler → Enqueue Job → PostgreSQL → River Worker → Process
                                                          ↓
                                              SSE Broadcast ← Progress Events
```

## Authentication

Better Auth with **Magic Link** (passwordless) authentication.

### Flow

1. User enters email on `/login`
2. Backend sends Magic Link email
3. User clicks link → Token verified → Logged in

### Pages

- `/login` - Magic Link login
- `/magic-link/verify` - Link verification UI
- `/verify-email` - Email verification for new users
- `/settings` - Session management (view/revoke sessions)
- `/dashboard` - Protected area

### Features

- **Passwordless**: Magic Link authentication
- **Rate Limiting**: 3 requests per minute
- **Session Management**: View and revoke sessions
- **Login Notifications**: Email on new device/IP
- **Cross-Tab Sync**: Logout syncs across tabs

### Client Usage

```tsx
import { signIn, signOut, authClient } from "@shared/lib/auth-client"

// Request Magic Link
await signIn.magicLink({
  email,
  callbackURL: "/dashboard",
})

// Sign out
await signOut()

// List sessions
const { data: sessions } = await authClient.listSessions()

// Revoke session
await authClient.revokeSession({ token })
```

### Local Email Testing

Mailpit is included for local email testing:

```bash
just dev  # Starts Mailpit on port 8025
```

Open [http://localhost:8025](http://localhost:8025) to view emails.

## Environment Variables

### Frontend (`frontend/.env.local`)

```
DATABASE_URL=postgres://postgres:postgres@localhost:5432/nextgopg
NEXT_PUBLIC_APP_URL=http://localhost:3000
NEXT_PUBLIC_API_URL=http://localhost:8080
BETTER_AUTH_SECRET=<at-least-32-characters>
```

### Backend

```
DATABASE_URL=postgres://postgres:postgres@localhost:5432/nextgopg
PORT=8080
```

## Docker

### Development

```bash
just db-up    # PostgreSQL only
```

### Production

```bash
just docker-build   # Build images
just docker-up      # Start containers
just docker-down    # Stop containers
```

## Adding a Backend Feature

The backend follows **bounded contexts + Clean Architecture (DDD)**. Each new aggregate lives inside a context under `backend/internal/<ctx>/`. The easiest template is to copy the `stats/` context — it's small, complete, and uses every DDD pattern (aggregate root with `AggregateBase`, value-object constructors, domain events, repository port, ACL friendliness).

1. Create `backend/internal/<ctx>/` with four sub-folders: `domain/`, `application/`, `infrastructure/persistence/`, `interfaces/http/`. See [`backend/README.md`](./backend/README.md#architecture) for the full layer guide.
2. Domain: pure types only — no GORM, no I/O. Use value-object constructors for invariants and embed `shared.AggregateBase` if the aggregate raises events.
3. Application: declare ports (`Repository`, `JobEnqueuer`, ...) and use-case structs with `Execute(ctx, ...)`. Use-cases pull domain events with `agg.PullEvents()` **before** `repo.Save(...)` so a repository round-trip can't drop them.
4. Infrastructure: GORM-tagged unexported twin types + mappers + repo impl. Assert the port with `var _ <ctx>app.Repository = (*Repository)(nil)`. Expose `Entities() []any` for AutoMigrate.
5. HTTP: handler imports only this context's `application/` package. Swagger annotations on every endpoint.
6. Wire in `backend/internal/composition/composition.go`. If cross-context data is needed, add an Anti-Corruption Layer right there (mirror `statsToExportsReader` / `authToNotificationsDirectory`).
7. Run `just api` to regenerate Swagger + the Orval TypeScript client.

### API Workflow

`just api` automatically runs:
1. **swag init** → Generates Swagger from Go comments
2. **orval** → Generates TypeScript React Query Hooks

## Further Documentation

- [Technical Docs (.docs)](./.docs/README.md) - LLM-friendly Tech Stack Docs
- [Frontend README](./frontend/README.md)
- [Backend README](./backend/README.md)
