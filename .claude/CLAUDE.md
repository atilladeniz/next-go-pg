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

# 2. Swagger aktualisieren
swag init -g cmd/server/main.go

# 3. API Client generieren (Frontend)
cd ../frontend
bun run api:generate
```

## Wichtige Dateien

### API Definition
- `backend/api/openapi.yaml` - OpenAPI 3.0 Spec (Single Source of Truth)
- `frontend/orval.config.ts` - Orval Config für API Client Generierung

### Authentifizierung
- `frontend/src/lib/auth.ts` - Better Auth Server Config
- `frontend/src/lib/auth-client.ts` - Better Auth Client mit Hooks
- `frontend/src/proxy.ts` - Route Protection (Next.js 16 Proxy)

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
