# GocaTest

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
└── kamal-deploy.md     # Kamal Deployment
```

> **Always check `.docs/` first** before searching the internet!

## Tech Stack

| Component | Technology |
|-----------|------------|
| Frontend | Next.js 16, TypeScript, Tailwind CSS, shadcn/ui |
| Backend | Go, Gorilla Mux, Clean Architecture, GORM |
| Code Generator | **Goca CLI** (Go Clean Architecture) |
| Database | PostgreSQL 16 |
| Auth | Better Auth |
| API | Swagger/swag → Orval |

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
cd gocatest

# Install dependencies
make install

# Start database
make db-up

# Create Better Auth tables
cd frontend && DATABASE_URL="postgres://postgres:postgres@localhost:5432/gocatest" bunx @better-auth/cli migrate -y

# Start development
make dev
```

Open:
- Frontend: [http://localhost:3000](http://localhost:3000)
- Backend: [http://localhost:8080](http://localhost:8080)

## Project Structure

```
gocatest/
├── backend/                 # Go Backend
│   ├── api/
│   │   └── openapi.yaml     # API Specification (Source of Truth)
│   ├── cmd/server/          # Entrypoint
│   ├── internal/
│   │   ├── handler/         # HTTP Handler
│   │   └── middleware/      # Auth, CORS Middleware
│   └── pkg/config/          # Configuration
├── frontend/                # Next.js Frontend
│   ├── src/
│   │   ├── api/             # Generated API Clients
│   │   ├── app/             # Next.js App Router
│   │   ├── components/      # React Components
│   │   └── lib/             # Utilities, Auth
│   └── orval.config.ts      # API Generator Config
├── docker-compose.dev.yml   # Dev Database
├── Makefile                 # Build Commands
└── README.md
```

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
make lint             # Linting
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
make fetch-docs url=<url>           # Fetch LLM-friendly docs
make fetch-docs url=<url> name=<n>  # With custom filename
```

## API Workflow

1. Edit OpenAPI Spec: `backend/api/openapi.yaml`
2. Generate TypeScript client: `make api`
3. Use generated hooks:

```tsx
import { useGetHello } from "@/api/endpoints/public/public"

function MyComponent() {
  const { data, isLoading } = useGetHello()
  // ...
}
```

## Authentication

Better Auth with Email/Password Login.

### Pages
- `/login` - Sign in
- `/register` - Sign up
- `/dashboard` - Protected area

### Client Usage

```tsx
import { signIn, signUp, signOut, useSession } from "@/lib/auth-client"

// Get session
const { data: session } = useSession()

// Sign in
await signIn.email({ email, password })

// Sign up
await signUp.email({ name, email, password })

// Sign out
await signOut()
```

## Environment Variables

### Frontend (`frontend/.env.local`)

```
DATABASE_URL=postgres://postgres:postgres@localhost:5432/gocatest
NEXT_PUBLIC_APP_URL=http://localhost:3000
NEXT_PUBLIC_API_URL=http://localhost:8080
BETTER_AUTH_SECRET=<at-least-32-characters>
```

### Backend

```
DATABASE_URL=postgres://postgres:postgres@localhost:5432/gocatest
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
