# Next-Go-PG - Project Context

## IMPORTANT: Technical Docs (.docs)

**ALWAYS check `.docs/` first** before searching the internet!

```
.docs/
├── nextjs.md           # Next.js 16 App Router
├── tanstack-query.md   # TanStack Query / React Query
├── better-auth.md      # Better Auth
├── gorm.md             # GORM ORM
├── goca.md             # Goca CLI
├── orval.md            # Orval API Client Generator
├── kamal-deploy.md     # Kamal Deployment (Docker)
└── ...                 # More Tech Stack Docs
```

LLM-friendly documentation for the entire tech stack can be found there.

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
- **Package Manager**: Bun

### Backend (`/backend`)

- **Language**: Go
- **Framework**: Gorilla Mux Router
- **Architecture**: Clean Architecture (Handler → Usecase → Repository → Domain)
- **Code Generator**: **Goca CLI** (Go Clean Architecture)
- **ORM**: GORM
- **Auth**: Better Auth Session Validation
- **API Docs**: Swagger/swag
- **Module**: `github.com/atilladeniz/next-go-pg/backend`

### Infrastructure

- **Database**: PostgreSQL 16 (Docker)
- **Dev Environment**: Docker Compose for DB

---

## IMPORTANT: Goca for Backend Development

### What is Goca?

Goca is a CLI tool for Go Clean Architecture code generation. It generates consistent, type-safe code structures with correct import paths.

### When to use Goca?

**ALWAYS** when creating new structures in the backend:

| Task | Goca Command |
|------|--------------|
| New Entity/Model | `goca make entity <Name>` |
| New Repository | `goca make repository <Name>` |
| New UseCase | `goca make usecase <Name>` |
| New Handler | `goca make handler <Name>` |
| Complete Feature | `goca feature <Name> --fields "..."` |

### Goca Commands

```bash
# Run in backend/ directory!
cd backend

# Create Entity (Domain Layer)
goca make entity UserStats

# Create Repository (Data Layer)
goca make repository UserStats

# Create UseCase (Business Logic Layer)
goca make usecase UserStats

# Create Handler (HTTP Layer)
goca make handler UserStats

# Complete Feature with all Layers
goca feature Product --fields "name:string,price:float64,stock:int"

# Feature with Validation
goca feature Order --fields "userId:string,total:float64" --validation

# Integrate all Features
goca integrate --all

# Check Goca Version
goca version
```

### Goca Configuration

Configuration is in `backend/.goca.yaml`:

- `module`: Go module path (github.com/atilladeniz/next-go-pg/backend)
- `architecture.layers`: Enabled layers (domain, usecase, repository, handler)
- `database.type`: postgres
- `generation.swagger.enabled`: true

### Clean Architecture Layers

```
backend/internal/
├── domain/           # Entities, Business Rules (goca make entity)
├── usecase/          # Application Logic (goca make usecase)
├── repository/       # Data Access (goca make repository)
├── handler/          # HTTP/API (goca make handler)
└── middleware/       # Cross-cutting concerns
```

### Why Goca instead of manual?

1. **Correct Imports**: Reads module path from .goca.yaml
2. **Consistency**: Same structure for all features
3. **Clean Architecture**: Enforces layer separation
4. **Swagger**: Generates API documentation automatically
5. **Tests**: Can generate test stubs

### Example: Adding a New Feature

```bash
# 1. Generate feature
cd backend
goca feature Invoice --fields "userId:string,amount:float64,status:string"

# 2. Add entity to registry (internal/domain/registry.go)
# Add &Invoice{} to AllEntities() function

# 3. Swagger + Orval (one command!)
cd ..
make api

# 4. Restart backend (migration runs automatically)
make dev-backend
```

### Entity Registry

New entities must be registered in `backend/internal/domain/registry.go`:

```go
func AllEntities() []interface{} {
    return []interface{}{
        &UserStats{},
        &Invoice{},  // ← New entity here
    }
}
```

This is the **ONLY** place for AutoMigrate!

### API Generation Workflow

`make api` automatically runs:

1. **swag init** → Generates `backend/docs/swagger.json` from Go comments
2. **orval** → Generates TypeScript Hooks in `frontend/src/shared/api/`

```bash
# Run after every API change:
make api

# Or separately:
make swagger     # Only generate Swagger
cd frontend && bunx orval  # Only run Orval
```

### Important for Claude

When you modify backend endpoints:

1. Add Swagger comments to Handler (`// @Summary`, `// @Router`, etc.)
2. Run `make api`
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
├── widgets/
│   └── header/                 # Header + ModeToggle
│       ├── ui/
│       └── index.ts
├── features/
│   ├── auth/                   # Login, Register, AuthSync
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
    │   ├── auth-client/        # Client-safe Auth
    │   ├── auth-server/        # Server-only Auth
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
export { RegisterForm } from "./ui/register-form"
export { useAuthSync, broadcastSignOut } from "./model/use-auth-sync"
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
- `frontend/src/shared/lib/auth-client/` - Better Auth Client (only for actions!)
- `frontend/src/features/auth/` - Login/Register Forms + Cross-Tab Sync

### Real-time

- `backend/internal/sse/broker.go` - SSE Broker
- `frontend/src/features/stats/model/use-sse.ts` - SSE Client Hook

### UI Components (FSD Paths)

- `frontend/src/shared/ui/` - shadcn/ui components
- `frontend/src/widgets/header/` - App Header with Mode Toggle
- `frontend/src/shared/config/` - Providers, Theme Config

## Commands

```bash
# Development
make dev              # Start DB + Frontend + Backend
make dev-frontend     # Frontend only (localhost:3000)
make dev-backend      # Backend only (localhost:8080)

# Database
make db-up            # Start PostgreSQL
make db-down          # Stop PostgreSQL
make db-reset         # Reset database

# API Generation
make api              # Generate TypeScript client from OpenAPI

# Quality
make lint             # Biome + Steiger Linting
make lint-fix         # Auto-fix Lint Errors
make typecheck        # TypeScript Check

# Security
make security-scan    # Scan for secrets/sensitive data
make setup-hooks      # Setup Git Hooks

# Build
make build            # Build Frontend + Backend
make build-frontend   # Next.js Production Build
make build-backend    # Go Binary
```

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
make security-scan

# Emergency bypass (NOT RECOMMENDED)
git commit --no-verify
```

## Project Structure

```
next-go-pg/
├── backend/
│   ├── cmd/server/           # Entrypoint
│   ├── internal/
│   │   ├── domain/           # Entities (goca make entity)
│   │   ├── usecase/          # Business Logic (goca make usecase)
│   │   ├── repository/       # Data Access (goca make repository)
│   │   ├── handler/          # HTTP Handler (goca make handler)
│   │   └── middleware/       # Auth, CORS
│   └── docs/                 # Swagger JSON (generated)
├── frontend/
│   ├── src/
│   │   ├── app/              # Next.js App Router
│   │   ├── widgets/          # Composite UI (Header)
│   │   ├── features/         # User Interactions (Auth, Stats)
│   │   ├── entities/         # Business Objects (User)
│   │   └── shared/           # Reusable (UI, API, Lib)
│   └── orval.config.ts
└── docker-compose.dev.yml
```

## Conventions

### Code Style

- Tabs for indentation (Biome Config)
- Double quotes for strings
- No semicolons (unless required)
- German UI texts

### Git Commits

- English commit messages
- Prefix: Add, Update, Fix, Remove
- No "Generated by" tags

### API

- OpenAPI 3.0 as source of truth
- Run `make api` after spec changes
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
```

### Backend

```
DATABASE_URL=postgres://postgres:postgres@localhost:5432/nextgopg
PORT=8080
```
