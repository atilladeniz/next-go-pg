---
name: auth-pages
description: Create and manage authentication pages with server-side session handling. Use when adding login, register, or protected pages WITHOUT flicker/skeleton.
allowed-tools: Read, Edit, Write, Glob
---

# Authentication Pages

Better Auth integration mit Next.js 16 - **Server-Side First** für flicker-freie UX.

## Architektur-Prinzip

```
Server Component (Session prüfen) → Client Component (Interaktivität)
```

**KEIN useSession() in Protected Pages!** → Verursacht Flicker.

## Bestehende Auth Pages

- `/login` - Login mit Email/Password
- `/register` - Registrierung
- `/dashboard` - Protected Dashboard

## Server-Side Session (KEIN FLICKER!)

### Protected Page Pattern

```tsx
// app/(protected)/dashboard/page.tsx - SERVER COMPONENT
import { auth } from "@/lib/auth"
import { headers } from "next/headers"
import { redirect } from "next/navigation"
import { Header } from "./header"
import { DashboardContent } from "./dashboard-content"

async function getData() {
  const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/v1/data`, {
    headers: { Cookie: (await headers()).get("cookie") ?? "" },
    cache: "no-store",
  })
  if (!res.ok) return null
  return res.json()
}

export default async function DashboardPage() {
  // Server-side Session Check - KEIN Flicker!
  const session = await auth.api.getSession({ headers: await headers() })
  if (!session) redirect("/login")

  // Daten server-side laden
  const data = await getData()

  return (
    <div className="min-h-screen bg-background">
      {/* User als Prop übergeben, nicht useSession! */}
      <Header user={session.user} />
      {/* Initial Data für React Query */}
      <DashboardContent initialData={data} />
    </div>
  )
}
```

### Header mit User Prop

```tsx
// app/(protected)/dashboard/header.tsx
"use client"

import { signOut } from "@/lib/auth-client"
import { Button } from "@/components/ui/button"

type User = { name: string; email: string }

export function Header({ user }: { user: User }) {
  return (
    <header className="border-b">
      <div className="container flex h-16 items-center justify-between">
        <span>Willkommen, {user.name}</span>
        <Button variant="outline" onClick={() => signOut()}>
          Abmelden
        </Button>
      </div>
    </header>
  )
}
```

### Client Component mit initialData

```tsx
// app/(protected)/dashboard/dashboard-content.tsx
"use client"

import { useGetData } from "@/api/endpoints/data/data"
import { useSSE } from "@/hooks/use-sse"

export function DashboardContent({ initialData }: { initialData: Data | null }) {
  // Real-time Updates via SSE
  useSSE()

  // React Query mit Server-Daten als Initial
  const { data: response } = useGetData({
    query: {
      initialData: initialData
        ? { data: initialData, status: 200 as const, headers: new Headers() }
        : undefined,
    },
  })

  const data = response?.status === 200 ? response.data : initialData

  return <div>{/* Render data */}</div>
}
```

## Auth Client (nur für Actions!)

```tsx
import { signIn, signUp, signOut } from "@/lib/auth-client"

// Sign In (in Login Page)
await signIn.email({ email, password })

// Sign Up (in Register Page)
await signUp.email({ name, email, password })

// Sign Out (in Header Button)
await signOut()
```

## Route Protection

### Proxy (middleware)

`frontend/src/proxy.ts` prüft Cookie-basiert:

```typescript
export function proxy(request: NextRequest) {
  const sessionCookie = request.cookies.get("better-auth.session_token")

  // Protected Route ohne Session → Login
  if (isProtectedRoute && !sessionCookie) {
    return NextResponse.redirect(new URL("/login", request.url))
  }

  // Auth Route mit Session → Dashboard
  if (isAuthRoute && sessionCookie) {
    return NextResponse.redirect(new URL("/dashboard", request.url))
  }
}
```

### Ordner-Struktur

```
src/app/
├── (auth)/           # Login, Register (redirect wenn eingeloggt)
│   ├── login/
│   └── register/
├── (protected)/      # Erfordert Auth
│   └── dashboard/
│       ├── page.tsx           # Server Component
│       ├── header.tsx         # Client Component
│       └── dashboard-content.tsx
└── page.tsx          # Public Home
```

## WICHTIG: Was NICHT tun

❌ **NICHT** `useSession()` in Protected Pages
❌ **NICHT** Skeleton für Session Loading
❌ **NICHT** `isPending` Check für Session
❌ **NICHT** Client-seitiger Redirect

✅ **STATTDESSEN**:
- Server Component prüft Session
- `redirect()` wenn keine Session
- User als Prop an Client Components
- `initialData` für React Query

## Real-Time Updates

SSE Hook für automatische Daten-Aktualisierung:

```tsx
"use client"

import { useSSE } from "@/hooks/use-sse"

export function MyComponent() {
  useSSE() // Verbindet zu /api/v1/events
  // React Query wird automatisch invalidiert bei SSE Events
}
```
