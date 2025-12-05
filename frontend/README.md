# Frontend

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
│   ├── (auth)/                 # Auth Pages (Login, Register)
│   ├── (protected)/            # Protected Pages (Dashboard)
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
│   ├── auth/
│   │   ├── ui/
│   │   │   ├── login-form.tsx
│   │   │   └── register-form.tsx
│   │   ├── model/
│   │   │   └── use-auth-sync.ts
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
    │   ├── auth-client/        # Better Auth Client
    │   ├── auth-server/        # Better Auth Server
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
export { RegisterForm } from "./ui/register-form"
export { useAuthSync, broadcastSignOut } from "./model/use-auth-sync"
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

Better Auth with Email/Password:

```tsx
import { signIn, signUp, signOut, useSession } from "@shared/lib/auth-client"

// Session Hook
const { data: session, isPending } = useSession()

// Sign In
await signIn.email({ email, password })

// Sign Up
await signUp.email({ name, email, password })

// Sign Out
await signOut()
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

## UI Components

All shadcn/ui components are in `@shared/ui/`:

```tsx
import { Button } from "@shared/ui/button"
import { Card } from "@shared/ui/card"
import { Input } from "@shared/ui/input"
```

## Dark Mode

Theme Toggle in Header. Supports Light, Dark, System.

## Environment Variables

Create `.env.local`:

```
DATABASE_URL=postgres://postgres:postgres@localhost:5432/nextgopg
NEXT_PUBLIC_APP_URL=http://localhost:3000
NEXT_PUBLIC_API_URL=http://localhost:8080
BETTER_AUTH_SECRET=<at-least-32-characters>
```
