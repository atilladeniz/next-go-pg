# GocaTest

Full-Stack Monorepo mit Next.js Frontend und Go Backend.

## Tech Stack

| Komponente | Technologie |
|------------|-------------|
| Frontend | Next.js 16, TypeScript, Tailwind CSS, shadcn/ui |
| Backend | Go, Chi Router, Clean Architecture |
| Datenbank | PostgreSQL 16 |
| Auth | Better Auth |
| API | OpenAPI 3.0, Orval |

## Voraussetzungen

- [Bun](https://bun.sh/) (Frontend)
- [Go 1.21+](https://go.dev/) (Backend)
- [Docker](https://www.docker.com/) (Datenbank)

## Schnellstart

```bash
# Repository klonen
git clone <repo-url>
cd gocatest

# Dependencies installieren
make install

# Datenbank starten
make db-up

# Better Auth Tabellen erstellen
cd frontend && DATABASE_URL="postgres://postgres:postgres@localhost:5432/gocatest" bunx @better-auth/cli migrate -y

# Development starten
make dev
```

Öffne:
- Frontend: [http://localhost:3000](http://localhost:3000)
- Backend: [http://localhost:8080](http://localhost:8080)

## Projektstruktur

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
│   │   ├── api/             # Generierte API Clients
│   │   ├── app/             # Next.js App Router
│   │   ├── components/      # React Components
│   │   └── lib/             # Utilities, Auth
│   └── orval.config.ts      # API Generator Config
├── docker-compose.dev.yml   # Dev Database
├── Makefile                 # Build Commands
└── README.md
```

## Make Commands

### Development

```bash
make dev              # DB + Frontend + Backend starten
make dev-frontend     # Nur Frontend (localhost:3000)
make dev-backend      # Nur Backend (localhost:8080)
```

### Datenbank

```bash
make db-up            # PostgreSQL starten
make db-down          # PostgreSQL stoppen
make db-reset         # Datenbank zurücksetzen
```

### API

```bash
make api              # TypeScript Client aus OpenAPI generieren
```

### Quality

```bash
make lint             # Linting
make lint-fix         # Auto-fix
make typecheck        # TypeScript Check
make test             # Tests ausführen
```

### Build

```bash
make build            # Frontend + Backend
make build-frontend   # Next.js Production Build
make build-backend    # Go Binary
```

## API Workflow

1. OpenAPI Spec bearbeiten: `backend/api/openapi.yaml`
2. TypeScript Client generieren: `make api`
3. Generierte Hooks nutzen:

```tsx
import { useGetHello } from "@/api/endpoints/public/public"

function MyComponent() {
  const { data, isLoading } = useGetHello()
  // ...
}
```

## Authentifizierung

Better Auth mit Email/Passwort Login.

### Seiten
- `/login` - Anmeldung
- `/register` - Registrierung
- `/dashboard` - Protected Bereich

### Client Usage

```tsx
import { signIn, signUp, signOut, useSession } from "@/lib/auth-client"

// Session abrufen
const { data: session } = useSession()

// Anmelden
await signIn.email({ email, password })

// Registrieren
await signUp.email({ name, email, password })

// Abmelden
await signOut()
```

## Environment Variables

### Frontend (`frontend/.env.local`)

```
DATABASE_URL=postgres://postgres:postgres@localhost:5432/gocatest
NEXT_PUBLIC_APP_URL=http://localhost:3000
NEXT_PUBLIC_API_URL=http://localhost:8080
BETTER_AUTH_SECRET=<mindestens-32-zeichen>
```

### Backend

```
DATABASE_URL=postgres://postgres:postgres@localhost:5432/gocatest
PORT=8080
```

## Docker

### Development

```bash
make db-up    # Nur PostgreSQL
```

### Production

```bash
make docker-build   # Images bauen
make docker-up      # Container starten
make docker-down    # Container stoppen
```

## Weitere Dokumentation

- [Frontend README](./frontend/README.md)
- [API Specification](./backend/api/openapi.yaml)
