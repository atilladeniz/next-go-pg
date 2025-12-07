# Next-Go-PG

[![CI](https://github.com/atilladeniz/next-go-pg/actions/workflows/ci.yml/badge.svg)](https://github.com/atilladeniz/next-go-pg/actions/workflows/ci.yml)

Full-Stack Monorepo with Next.js Frontend and Go Backend.

## Technical Documentation

**IMPORTANT:** LLM-friendly documentation for the entire tech stack can be found in `.docs/`:

```
.docs/
├── nextjs.md           # Next.js 16 App Router
├── tanstack-query.md   # TanStack Query
├── better-auth.md      # Better Auth
├── gorm.md             # GORM ORM
├── goca.md             # Goca CLI
├── orval.md            # Orval API Generator
├── shadcn.md           # shadcn/ui
├── tailwind.md         # Tailwind CSS 4
├── kamal-deploy.md     # Kamal Deployment
├── logging.md          # Logging (zerolog + Pino)
└── disaster-recovery.md # Database Backups & Recovery
```

> **Always check `.docs/` first** before searching the internet!

## Tech Stack

| Component | Technology |
|-----------|------------|
| Frontend | Next.js 16, TypeScript, Tailwind CSS, shadcn/ui |
| Frontend Architecture | **Feature-Sliced Design (FSD)** |
| Backend | Go, Gorilla Mux, Clean Architecture, GORM |
| Code Generator | **Goca CLI** (Go Clean Architecture) |
| Database | PostgreSQL 16 |
| Auth | Better Auth |
| API | Swagger/swag → Orval |
| Logging | zerolog (Go) + Pino (Next.js) |
| Log Aggregation | Grafana + Loki + Promtail |
| Database Backups | postgres-backup-s3 + RustFS (S3-compatible) |
| Linting | Biome + Steiger (FSD) |

## Prerequisites

- [Bun](https://bun.sh/) (Frontend)
- [Go 1.21+](https://go.dev/) (Backend)
- [Goca CLI](https://github.com/sazardev/goca) (Backend Code Generation)
- [Docker](https://www.docker.com/) (Database)

### Install CLI Tools

```bash
# Goca (Backend Code Generation)
go install github.com/sazardev/goca@latest

# Gitleaks (Security Scanning)
brew install gitleaks

# Sitefetch (Documentation Fetching)
bun install -g sitefetch
```

## Quick Start

```bash
# Clone repository
git clone <repo-url>
cd next-go-pg

# Install dependencies
make install

# Start database
make db-up

# Create Better Auth tables
cd frontend && DATABASE_URL="postgres://postgres:postgres@localhost:5432/nextgopg" bunx @better-auth/cli migrate -y

# Start development
make dev
```

Open:
- Frontend: [http://localhost:3000](http://localhost:3000)
- Backend: [http://localhost:8080](http://localhost:8080)

## Project Structure

```
next-go-pg/
├── backend/                 # Go Backend (Clean Architecture)
│   ├── cmd/server/          # Entrypoint
│   ├── internal/
│   │   ├── domain/          # Entities (goca make entity)
│   │   ├── usecase/         # Business Logic (goca make usecase)
│   │   ├── repository/      # Data Access (goca make repository)
│   │   ├── handler/         # HTTP Handler (goca make handler)
│   │   └── middleware/      # Auth, CORS, Logging
│   ├── pkg/logger/          # zerolog Logger
│   └── docs/                # Swagger (generated)
├── frontend/                # Next.js Frontend (FSD Architecture)
│   ├── src/
│   │   ├── app/             # Next.js App Router
│   │   ├── widgets/         # Composite UI (Header)
│   │   ├── features/        # User Interactions (Auth, Stats)
│   │   ├── entities/        # Business Objects (User)
│   │   └── shared/          # Reusable (UI, API, Lib, Logger)
│   └── orval.config.ts      # API Generator Config
├── docker-compose.dev.yml   # Dev Database + Mailpit
├── docker-compose.logging.yml # Logging Stack (Grafana + Loki)
├── docker-compose.backup.yml  # Backup Stack (postgres-backup-s3 + RustFS)
├── deploy/
│   ├── loki/                # Loki & Promtail Config
│   └── grafana/             # Grafana Provisioning
├── Makefile                 # Build Commands
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

### Setup (automatic with `make install`)

```bash
# Install gitleaks
brew install gitleaks    # macOS
# or
go install github.com/gitleaks/gitleaks/v8@latest

# Setup git hooks (runs automatically with make install)
make setup-hooks
```

### What it scans for

- API keys, tokens, and passwords
- Absolute paths with usernames
- Database URLs with embedded credentials
- Private keys and certificates
- AWS/GCP/Azure credentials

### Manual scan

```bash
make security-scan    # Scan entire codebase
```

### Bypass (not recommended)

```bash
git commit --no-verify
```

## Make Commands

### Development

```bash
make dev              # Start DB + Frontend + Backend
make dev-frontend     # Frontend only (localhost:3000)
make dev-backend      # Backend only (localhost:8080)
```

### Database

```bash
make db-up            # Start PostgreSQL
make db-down          # Stop PostgreSQL
make db-reset         # Reset database
```

### API

```bash
make api              # Generate TypeScript client from OpenAPI
```

### Quality

```bash
make lint             # Biome + Steiger (FSD) Linting
make lint-fix         # Auto-fix
make typecheck        # TypeScript Check
make test             # Run tests
make security-scan    # Scan for secrets
```

### Build

```bash
make build            # Frontend + Backend
make build-frontend   # Next.js Production Build
make build-backend    # Go Binary
```

### Documentation

```bash
make search-docs q="query"          # Search docs with semantic search
make search-docs q="query" n=10     # Search with custom result count
make fetch-docs url=<url>           # Fetch LLM-friendly docs
make fetch-docs url=<url> name=<n>  # With custom filename
```

### Logging Stack

```bash
make logs-up      # Start Grafana + Loki + Promtail
make logs-down    # Stop logging stack
make logs-open    # Open Grafana (localhost:3001)
make logs-query q='{level="error"}'  # Query logs via CLI
```

### Database Backups

Fully automatic PostgreSQL backups to S3-compatible storage (RustFS).

```bash
make backup-up       # Start automatic backup system
make backup-down     # Stop backup stack
make backup-now      # Create backup immediately
make backup-list     # List all backups in S3
make backup-restore  # Restore from latest backup
```

- **Schedule**: Daily (configurable via `BACKUP_SCHEDULE`)
- **Retention**: 7 days (configurable via `BACKUP_KEEP_DAYS`)
- **Storage**: RustFS Console at http://localhost:9001

See [Disaster Recovery](.docs/disaster-recovery.md) for details.

## API Workflow

1. Add Swagger comments to Go handler
2. Generate TypeScript client: `make api`
3. Use generated hooks:

```tsx
import { useGetStats } from "@shared/api/endpoints/users/users"

function MyComponent() {
  const { data, isLoading } = useGetStats()
  // ...
}
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
make dev  # Starts Mailpit on port 8025
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
make db-up    # PostgreSQL only
```

### Production

```bash
make docker-build   # Build images
make docker-up      # Start containers
make docker-down    # Stop containers
```

## Generate Backend Features (Goca)

```bash
cd backend

# New feature with all layers
goca feature Product --fields "name:string,price:float64,stock:int"

# Entity only
goca make entity Product

# Repository only
goca make repository Product

# Swagger + Orval (one command from root!)
cd ..
make api
```

### API Workflow

`make api` automatically runs:
1. **swag init** → Generates Swagger from Go comments
2. **orval** → Generates TypeScript React Query Hooks

## Further Documentation

- [Technical Docs (.docs)](./.docs/README.md) - LLM-friendly Tech Stack Docs
- [Frontend README](./frontend/README.md)
- [Backend README](./backend/README.md)
- [Goca Documentation](https://github.com/sazardev/goca)
