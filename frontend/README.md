# Frontend

[![Frontend CI](https://github.com/atilladeniz/next-go-pg/actions/workflows/ci.yml/badge.svg)](https://github.com/atilladeniz/next-go-pg/actions/workflows/ci.yml)

Next.js 16 Frontend with TypeScript, Tailwind CSS, shadcn/ui and **Feature-Sliced Design (FSD)** architecture.

> **Tip:** LLM-friendly documentation for Next.js, TanStack Query, Better Auth, etc. can be found in `../.docs/`

## Tech Stack

- **Next.js 16** - App Router with Turbopack
- **TypeScript 5.9** - Type safety
- **Tailwind CSS 4** - Styling
- **shadcn/ui** - UI Components (neutral theme)
- **TanStack Query** - Server State Management
- **Better Auth** - Authentication
- **Orval** - API Client Generation
- **Pino** - Structured Logging
- **Biome** - Linting & Formatting
- **Steiger** - FSD Architecture Linting

## Getting Started

```bash
# Install dependencies
bun install

# Start development server
bun run dev
```

Open [http://localhost:3000](http://localhost:3000) in your browser.

## Available Scripts

```bash
bun run dev          # Development server with Turbopack
bun run build        # Production build
bun run start        # Production server
bun run lint         # Biome + Steiger linting
bun run lint:fix     # Auto-fix lint errors
bun run typecheck    # TypeScript check
bun run api:generate # Generate API client from OpenAPI
```

## Feature-Sliced Design (FSD)

This project follows the [Feature-Sliced Design](https://feature-sliced.design/) methodology.

### Layer Hierarchy

```
app/        → widgets, features, entities, shared
widgets/    → features, entities, shared
features/   → entities, shared
entities/   → shared
shared/     → (only external libs)
```

### Project Structure

```
src/
├── app/                        # Next.js App Router
│   ├── (auth)/                 # Auth Pages
│   │   ├── login/              # Magic Link Login
│   │   ├── magic-link/verify/  # Magic Link Verification UI
│   │   └── verify-email/       # Email Verification UI
│   ├── (protected)/            # Protected Pages
│   │   ├── dashboard/
│   │   └── settings/           # Session Management
│   └── api/auth/               # Better Auth API Route
│
├── widgets/                    # Composite UI Blocks
│   └── header/
│       ├── ui/
│       │   ├── header.tsx
│       │   └── mode-toggle.tsx
│       └── index.ts
│
├── features/                   # User Interactions
│   ├── auth/                   # Magic Link Login
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
│   │   │   ├── sessions-list.tsx
│   │   │   └── session-card.tsx
│   │   ├── model/
│   │   │   └── use-sessions.ts
│   │   └── index.ts
│   └── stats/
│       ├── ui/
│       │   └── stats-grid.tsx
│       ├── model/
│       │   └── use-sse.ts
│       └── index.ts
│
├── entities/                   # Business Objects
│   └── user/
│       ├── ui/
│       │   └── user-info.tsx
│       ├── model/
│       │   └── types.ts
│       └── index.ts
│
└── shared/                     # Reusable Code
    ├── ui/                     # shadcn/ui Components
    ├── api/                    # Orval-generated (endpoints, models)
    ├── lib/
    │   ├── auth-client/        # Better Auth Client (Magic Link)
    │   ├── auth-server/        # Better Auth Server Config
    │   ├── geo/                # User Agent Parsing
    │   ├── logger/             # Pino Logger
    │   ├── query-client.ts     # TanStack Query Client
    │   └── utils.ts            # cn() helper
    └── config/
        ├── providers.tsx       # App Providers
        └── theme-provider.tsx  # Theme Provider
```

### Import Rules

```tsx
// ✅ ALLOWED: Import from lower layer
import { Button } from "@shared/ui/button"
import { SessionUser } from "@entities/user"
import { LoginForm } from "@features/auth"

// ❌ FORBIDDEN: Import from higher layer
// In shared/ NEVER import features/!
// In entities/ NEVER import features/!
// In features/ NEVER import other features/!
```

### TypeScript Path Aliases

```tsx
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
export { useAuthSync, broadcastSignOut } from "./model/use-auth-sync"

// features/user-settings/index.ts
export { SessionsList } from "./ui/sessions-list"
export { useSessions } from "./model/use-sessions"
```

### Adding New Features

1. Create feature folder: `src/features/<name>/`
2. Add segments: `ui/`, `model/`, `api/` (as needed)
3. Create public API: `index.ts` with all exports
4. Import only from `shared/` or `entities/`

## API Client Generation

The TypeScript API client is generated from Swagger:

```bash
bun run api:generate
```

Uses Orval to generate React Query hooks from `../backend/docs/swagger.json`.

## Authentication

**Magic Link** (passwordless) authentication with Better Auth:

```tsx
import { signIn, signOut } from "@shared/lib/auth-client"

// Request Magic Link
await signIn.magicLink({
  email,
  callbackURL: "/dashboard",
  newUserCallbackURL: "/dashboard",
  errorCallbackURL: "/login?error=verification_failed",
})

// Sign Out
await signOut()
```

### Flow

1. User enters email on `/login`
2. Backend sends Magic Link email
3. User clicks link → Token verified → Logged in

### Session Management

View and revoke sessions at `/settings`:

```tsx
import { useSessions } from "@features/user-settings"

const { sessions, revokeSession, revokeOtherSessions } = useSessions()

// Revoke single session
await revokeSession(token)

// Revoke all other sessions
await revokeOtherSessions()
```

### Cross-Tab Synchronization

```tsx
import { broadcastSignOut, useAuthSync } from "@features/auth"
import { signOut } from "@shared/lib/auth-client"

export function Header({ user }: { user: User }) {
  useAuthSync() // Cross-Tab Listener

  const handleSignOut = async () => {
    await signOut()
    await broadcastSignOut() // Notify all other tabs
    router.push("/")
  }
}
```

### Features

- **Passwordless**: No passwords, just Magic Links
- **Rate Limiting**: 3 requests per minute
- **Session Management**: View/revoke sessions
- **Login Notifications**: Email on new device login
- **Cross-Tab Sync**: Logout syncs across tabs

## UI Components

All shadcn/ui components are in `@shared/ui/`:

```tsx
import { Button } from "@shared/ui/button"
import { Card } from "@shared/ui/card"
import { Input } from "@shared/ui/input"
```

## Dark Mode

Theme Toggle in Header. Supports Light, Dark, System.

## Logging

Structured logging with Pino:

```typescript
import { log, createLogger } from "@shared/lib/logger"

// Simple logging
log.info("Page loaded")
log.error("API failed", { endpoint: "/api/users" })

// Component-specific logger
const authLogger = createLogger("auth")
authLogger.info("Login submitted")
```

Development: Pretty colored output
Production: JSON for log aggregation (ELK, Datadog, etc.)

## Environment Variables

Create `.env.local`:

```
DATABASE_URL=postgres://postgres:postgres@localhost:5432/nextgopg
NEXT_PUBLIC_APP_URL=http://localhost:3000
NEXT_PUBLIC_API_URL=http://localhost:8080
BETTER_AUTH_SECRET=<at-least-32-characters>
LOG_LEVEL=info
```
