# GocaTest - Projekt Kontext

## Projektübersicht

Full-Stack Monorepo mit Next.js 16 Frontend und Go Backend, PostgreSQL Datenbank und Better Auth für Authentifizierung.

## Tech Stack

### Frontend (`/frontend`)
- **Framework**: Next.js 16 mit App Router und Turbopack
- **Sprache**: TypeScript 5.9
- **Styling**: Tailwind CSS 4 + shadcn/ui (neutral theme)
- **State**: TanStack Query (React Query)
- **Auth Client**: Better Auth React Client
- **API Client**: Orval-generierte Hooks aus OpenAPI Spec
- **Linting**: Biome
- **Package Manager**: Bun

### Backend (`/backend`)
- **Sprache**: Go
- **Framework**: Gorilla Mux Router
- **Architektur**: Clean Architecture (Handler → Usecase → Repository → Domain)
- **Code Generator**: **Goca CLI** (Go Clean Architecture)
- **ORM**: GORM
- **Auth**: Better Auth Session Validation
- **API Docs**: Swagger/swag
- **Modul**: `github.com/atilladeniz/gocatest/backend`

### Infrastruktur
- **Datenbank**: PostgreSQL 16 (Docker)
- **Dev Environment**: Docker Compose für DB

---

## WICHTIG: Goca für Backend-Entwicklung

### Was ist Goca?
Goca ist ein CLI-Tool für Go Clean Architecture Code-Generierung. Es generiert konsistente, typsichere Code-Strukturen mit korrekten Import-Pfaden.

### Wann Goca verwenden?
**IMMER** wenn im Backend neue Strukturen erstellt werden:

| Aufgabe | Goca Befehl |
|---------|-------------|
| Neues Entity/Model | `goca make entity <Name>` |
| Neues Repository | `goca make repository <Name>` |
| Neuer UseCase | `goca make usecase <Name>` |
| Neuer Handler | `goca make handler <Name>` |
| Komplettes Feature | `goca feature <Name> --fields "..."` |

### Goca Befehle

```bash
# Im backend/ Verzeichnis ausführen!
cd backend

# Entity erstellen (Domain Layer)
goca make entity UserStats

# Repository erstellen (Data Layer)
goca make repository UserStats

# UseCase erstellen (Business Logic Layer)
goca make usecase UserStats

# Handler erstellen (HTTP Layer)
goca make handler UserStats

# Komplettes Feature mit allen Layers
goca feature Product --fields "name:string,price:float64,stock:int"

# Feature mit Validierung
goca feature Order --fields "userId:string,total:float64" --validation

# Alle Features integrieren
goca integrate --all

# Goca Version prüfen
goca version
```

### Goca Konfiguration
Die Konfiguration liegt in `backend/.goca.yaml`:
- `module`: Go Modul-Pfad (github.com/atilladeniz/gocatest/backend)
- `architecture.layers`: Aktivierte Layer (domain, usecase, repository, handler)
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

### Warum Goca statt manuell?
1. **Korrekte Imports**: Liest Modul-Pfad aus .goca.yaml
2. **Konsistenz**: Gleiche Struktur für alle Features
3. **Clean Architecture**: Erzwingt Layer-Trennung
4. **Swagger**: Generiert API-Dokumentation automatisch
5. **Tests**: Kann Test-Stubs generieren

### Beispiel: Neues Feature hinzufügen

```bash
# 1. Feature generieren
cd backend
goca feature Invoice --fields "userId:string,amount:float64,status:string"

# 2. Swagger + Orval (ein Befehl!)
cd ..
make api
```

### API Generierung Workflow

`make api` führt automatisch aus:
1. **swag init** → Generiert `backend/docs/swagger.json` aus Go-Kommentaren
2. **orval** → Generiert TypeScript Hooks in `frontend/src/api/`

```bash
# Nach jeder API-Änderung ausführen:
make api

# Oder einzeln:
make swagger     # Nur Swagger generieren
cd frontend && bunx orval  # Nur Orval ausführen
```

### Wichtig für Claude

Wenn du Backend-Endpoints änderst:
1. Swagger-Kommentare in Handler hinzufügen (`// @Summary`, `// @Router`, etc.)
2. `make api` ausführen
3. Frontend kann die neuen Hooks nutzen (`useGetX`, `usePostX`, etc.)

---

## WICHTIG: Server-Side First (KEIN FLICKER!)

### Architektur-Prinzip

```
Server Component (Session + Daten laden) → Client Component (Interaktivität)
```

### Protected Page Pattern

```tsx
// app/(protected)/dashboard/page.tsx - SERVER COMPONENT
import { auth } from "@/lib/auth"
import { headers } from "next/headers"
import { redirect } from "next/navigation"

export default async function DashboardPage() {
  // 1. Session server-side prüfen
  const session = await auth.api.getSession({ headers: await headers() })
  if (!session) redirect("/login")

  // 2. Daten server-side laden
  const data = await getData()

  // 3. An Client Components übergeben
  return (
    <Header user={session.user} />
    <Content initialData={data} />
  )
}
```

### Client Component mit initialData

```tsx
// content.tsx - CLIENT COMPONENT
"use client"

import { useGetData } from "@/api/endpoints/data/data"
import { useSSE } from "@/hooks/use-sse"

export function Content({ initialData }) {
  useSSE() // Real-time Updates

  const { data } = useGetData({
    query: {
      initialData: initialData
        ? { data: initialData, status: 200 as const, headers: new Headers() }
        : undefined,
    },
  })
  // Kein Loading State nötig!
}
```

### Was NICHT tun

❌ `useSession()` in Protected Pages → Flicker
❌ Skeleton für Session Loading
❌ Client-seitiger Redirect
❌ `isPending` Check

### Was STATTDESSEN

✅ Server Component prüft Session
✅ `redirect()` wenn keine Session
✅ User als Prop an Client Components
✅ `initialData` für React Query

---

## SSE + React Query Pattern

Real-time Updates ohne Polling:

1. **Backend** sendet SSE Events bei Änderungen
2. **Frontend** `useSSE()` Hook hört auf Events
3. **React Query** wird automatisch invalidiert

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

## Wichtige Dateien

### API Definition
- `backend/docs/swagger.json` - Generiert aus Go-Kommentaren
- `frontend/orval.config.ts` - Orval Config für API Client Generierung

### Authentifizierung
- `frontend/src/lib/auth.ts` - Better Auth Server Config
- `frontend/src/lib/auth-client.ts` - Better Auth Client (nur für Actions!)
- `frontend/src/proxy.ts` - Route Protection (Cookie-basiert)

### Real-time
- `backend/internal/sse/broker.go` - SSE Broker
- `frontend/src/hooks/use-sse.ts` - SSE Client Hook

### UI Komponenten
- `frontend/src/components/ui/` - shadcn/ui Komponenten (alle installiert)
- `frontend/src/components/mode-toggle.tsx` - Dark Mode Toggle
- `frontend/src/components/theme-provider.tsx` - Theme Provider

## Befehle

```bash
# Development
make dev              # Start DB + Frontend + Backend
make dev-frontend     # Nur Frontend (localhost:3000)
make dev-backend      # Nur Backend (localhost:8080)

# Database
make db-up            # PostgreSQL starten
make db-down          # PostgreSQL stoppen
make db-reset         # Datenbank zurücksetzen

# API Generation
make api              # TypeScript Client aus OpenAPI generieren

# Quality
make lint             # Biome Linting
make lint-fix         # Auto-fix Lint Errors
make typecheck        # TypeScript Check

# Build
make build            # Frontend + Backend bauen
make build-frontend   # Next.js Production Build
make build-backend    # Go Binary
```

## Projektstruktur

```
gocatest/
├── backend/
│   ├── api/
│   │   └── openapi.yaml      # API Specification
│   ├── cmd/server/           # Entrypoint
│   ├── internal/
│   │   ├── handler/          # HTTP Handler
│   │   └── middleware/       # Auth, CORS
│   └── pkg/config/           # Configuration
├── frontend/
│   ├── src/
│   │   ├── api/              # Generierte API Clients
│   │   ├── app/
│   │   │   ├── (auth)/       # Login, Register
│   │   │   ├── (protected)/  # Dashboard
│   │   │   └── api/auth/     # Better Auth Handler
│   │   ├── components/
│   │   │   └── ui/           # shadcn Komponenten
│   │   └── lib/              # Auth, Utils
│   └── orval.config.ts
└── docker-compose.dev.yml
```

## Konventionen

### Code Style
- Tabs für Indentation (Biome Config)
- Double Quotes für Strings
- Keine Semicolons (außer nötig)
- Deutsche UI Texte

### Git Commits
- Englische Commit Messages
- Präfix: Add, Update, Fix, Remove
- Keine "Generated by" Tags

### API
- OpenAPI 3.0 als Source of Truth
- `make api` nach Spec-Änderungen ausführen
- Generated Files nicht manuell editieren

## Biome Ignores
- `src/api/endpoints/` - Orval generiert
- `src/api/models/` - Orval generiert
- `src/components/ui/` - shadcn generiert

## Environment Variables

### Frontend (.env.local)
```
DATABASE_URL=postgres://postgres:postgres@localhost:5432/gocatest
NEXT_PUBLIC_APP_URL=http://localhost:3000
NEXT_PUBLIC_API_URL=http://localhost:8080
BETTER_AUTH_SECRET=<secret>
```

### Backend
```
DATABASE_URL=postgres://postgres:postgres@localhost:5432/gocatest
PORT=8080
```
