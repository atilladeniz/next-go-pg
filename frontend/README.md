# Frontend

Next.js 16 Frontend mit TypeScript, Tailwind CSS und shadcn/ui.

## Tech Stack

- **Next.js 16** - App Router mit Turbopack
- **TypeScript 5.9** - Typsicherheit
- **Tailwind CSS 4** - Styling
- **shadcn/ui** - UI Komponenten (neutral theme)
- **TanStack Query** - Server State Management
- **Better Auth** - Authentifizierung
- **Orval** - API Client Generierung
- **Biome** - Linting & Formatting

## Erste Schritte

```bash
# Dependencies installieren
bun install

# Development Server starten
bun run dev
```

Öffne [http://localhost:3000](http://localhost:3000) im Browser.

## Verfügbare Scripts

```bash
bun run dev          # Development Server mit Turbopack
bun run build        # Production Build
bun run start        # Production Server
bun run lint         # Biome Linting
bun run lint:fix     # Auto-fix Lint Errors
bun run typecheck    # TypeScript Check
bun run api:generate # API Client aus OpenAPI generieren
```

## Projektstruktur

```
src/
├── api/                    # Generierte API Clients (Orval)
│   ├── endpoints/          # React Query Hooks
│   ├── models/             # TypeScript Types
│   └── custom-fetch.ts     # Fetch Wrapper
├── app/
│   ├── (auth)/             # Auth Pages (Login, Register)
│   ├── (protected)/        # Protected Pages (Dashboard)
│   ├── api/auth/           # Better Auth API Route
│   ├── layout.tsx          # Root Layout
│   └── page.tsx            # Home Page
├── components/
│   ├── ui/                 # shadcn/ui Komponenten
│   ├── mode-toggle.tsx     # Dark Mode Toggle
│   ├── providers.tsx       # App Providers
│   └── theme-provider.tsx  # Theme Provider
├── hooks/
│   └── use-mobile.ts       # Responsive Hook
├── lib/
│   ├── auth.ts             # Better Auth Server Config
│   ├── auth-client.ts      # Better Auth Client
│   └── utils.ts            # Utility Functions
└── proxy.ts                # Route Protection
```

## API Client Generierung

Der TypeScript API Client wird aus der OpenAPI Spec generiert:

```bash
bun run api:generate
```

Die Spec liegt unter `../backend/api/openapi.yaml`.

## Authentifizierung

Better Auth mit Email/Passwort:

```tsx
import { signIn, signUp, signOut, useSession } from "@/lib/auth-client"

// Session Hook
const { data: session, isPending } = useSession()

// Sign In
await signIn.email({ email, password })

// Sign Up
await signUp.email({ name, email, password })

// Sign Out
await signOut()
```

## UI Komponenten

Alle shadcn/ui Komponenten sind installiert:

```tsx
import { Button } from "@/components/ui/button"
import { Card } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
// ... etc
```

## Dark Mode

Theme Toggle im Dashboard Header. Unterstützt:
- Hell
- Dunkel
- System

```tsx
import { ModeToggle } from "@/components/mode-toggle"
```

## Environment Variables

Erstelle `.env.local`:

```
DATABASE_URL=postgres://postgres:postgres@localhost:5432/gocatest
NEXT_PUBLIC_APP_URL=http://localhost:3000
NEXT_PUBLIC_API_URL=http://localhost:8080
BETTER_AUTH_SECRET=<mindestens-32-zeichen>
```
