# Frontend

Next.js 16 Frontend with TypeScript, Tailwind CSS and shadcn/ui.

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
bun run lint         # Biome linting
bun run lint:fix     # Auto-fix lint errors
bun run typecheck    # TypeScript check
bun run api:generate # Generate API client from OpenAPI
```

## Project Structure

```
src/
├── api/                    # Generated API Clients (Orval)
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
│   ├── ui/                 # shadcn/ui Components
│   ├── mode-toggle.tsx     # Dark Mode Toggle
│   ├── providers.tsx       # App Providers
│   └── theme-provider.tsx  # Theme Provider
├── hooks/
│   ├── use-auth-sync.ts    # Cross-Tab Auth Sync
│   ├── use-mobile.ts       # Responsive Hook
│   └── use-sse.ts          # SSE Real-Time Updates
├── lib/
│   ├── auth.ts             # Better Auth Server Config
│   ├── auth-client.ts      # Better Auth Client
│   └── utils.ts            # Utility Functions
└── proxy.ts                # Route Protection
```

## API Client Generation

The TypeScript API client is generated from the OpenAPI spec:

```bash
bun run api:generate
```

The spec is located at `../backend/api/openapi.yaml`.

## Authentication

Better Auth with Email/Password:

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

### Cross-Tab Synchronization

Better Auth has no built-in cross-tab synchronization. Use `broadcast-channel`:

```bash
bun add broadcast-channel
```

```tsx
// In Client Components (e.g. Header)
import { broadcastSignOut, useAuthSync } from "@/hooks/use-auth-sync"
import { signOut } from "@/lib/auth-client"

export function Header({ user }: { user: User }) {
  useAuthSync() // Cross-Tab Listener

  const handleSignOut = async () => {
    await signOut()
    await broadcastSignOut() // Notify all other tabs
    router.push("/")
  }
}
```

Benefits:
- Safari Private Mode Support
- localStorage Fallback
- TypeScript Support

## UI Components

All shadcn/ui components are installed:

```tsx
import { Button } from "@/components/ui/button"
import { Card } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
// ... etc
```

## Dark Mode

Theme Toggle in Dashboard Header. Supports:
- Light
- Dark
- System

```tsx
import { ModeToggle } from "@/components/mode-toggle"
```

## Environment Variables

Create `.env.local`:

```
DATABASE_URL=postgres://postgres:postgres@localhost:5432/nextgopg
NEXT_PUBLIC_APP_URL=http://localhost:3000
NEXT_PUBLIC_API_URL=http://localhost:8080
BETTER_AUTH_SECRET=<at-least-32-characters>
```
