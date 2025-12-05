---
name: auth-pages
description: Create and manage authentication pages with server-side session handling. Use when adding login, register, or protected pages WITHOUT flicker/skeleton.
allowed-tools: Read, Edit, Write, Glob
---

# Authentication Pages (FSD)

Better Auth integration with Next.js 16 - **Server-Side First** for flicker-free UX.

## FSD Paths

```
src/
├── features/auth/              # Auth Feature
│   ├── ui/
│   │   ├── login-form.tsx
│   │   └── register-form.tsx
│   ├── model/
│   │   └── use-auth-sync.ts
│   └── index.ts
├── shared/lib/
│   ├── auth-client/            # Client-safe: signIn, signOut
│   └── auth-server/            # Server-only: getSession, auth
└── widgets/header/             # Header with Auth Sync
```

## Architecture Principle

```
Server Component (check session) → Client Component (interactivity)
```

**NO useSession() in Protected Pages!** → Causes flicker.

## Existing Auth Pages

- `/login` - Login with Email/Password
- `/register` - Registration
- `/dashboard` - Protected Dashboard

## Server-Side Session (NO FLICKER!)

### Protected Page Pattern

```tsx
// app/(protected)/dashboard/page.tsx - SERVER COMPONENT
import { getSession } from "@shared/lib/auth-server"
import { redirect } from "next/navigation"
import { Header } from "@widgets/header"
import { StatsGrid } from "@features/stats"

export default async function DashboardPage() {
  // Server-side Session Check - NO Flicker!
  const session = await getSession()
  if (!session) redirect("/login")

  return (
    <div className="min-h-screen bg-background">
      {/* Pass user as prop, not useSession! */}
      <Header user={session.user} />
      {/* HydrationBoundary for React Query - see data-fetching skill */}
      <StatsGrid />
    </div>
  )
}
```

### Header Widget

```tsx
// widgets/header/ui/header.tsx
"use client"

import type { SessionUser } from "@entities/user"
import { broadcastSignOut, useAuthSync } from "@features/auth"
import { signOut } from "@shared/lib/auth-client"
import { Button } from "@shared/ui/button"

export function Header({ user }: { user: SessionUser }) {
  useAuthSync() // Cross-Tab Listener

  const handleSignOut = async () => {
    await signOut()
    await broadcastSignOut()
    router.push("/")
  }

  return (
    <header className="border-b">
      <div className="container flex h-16 items-center justify-between">
        <span>Welcome, {user.name}</span>
        <Button variant="outline" onClick={handleSignOut}>
          Sign Out
        </Button>
      </div>
    </header>
  )
}
```

## Auth Client (only for actions!)

```tsx
import { signIn, signUp, signOut } from "@shared/lib/auth-client"

// Sign In (in Login Page)
await signIn.email({ email, password })

// Sign Up (in Register Page)
await signUp.email({ name, email, password })

// Sign Out (in Header Button)
await signOut()
```

## Cross-Tab Synchronization

Better Auth has **no built-in cross-tab synchronization**. Use `broadcast-channel` library:

```bash
bun add broadcast-channel
```

### useAuthSync Hook (FSD)

```tsx
// features/auth/model/use-auth-sync.ts
"use client"

import { BroadcastChannel } from "broadcast-channel"
import { useRouter } from "next/navigation"
import { useEffect } from "react"

type AuthMessage = { type: "SIGN_OUT" | "SIGN_IN"; timestamp: number }

export function useAuthSync() {
  const router = useRouter()

  useEffect(() => {
    const channel = new BroadcastChannel<AuthMessage>("auth-sync", {
      type: "localstorage", // Fallback for Safari Private Mode
    })

    channel.onmessage = (msg) => {
      if (msg.type === "SIGN_OUT") {
        router.push("/login")
        router.refresh()
      }
    }

    return () => channel.close()
  }, [router])
}

export async function broadcastSignOut() {
  const channel = new BroadcastChannel<AuthMessage>("auth-sync", { type: "localstorage" })
  await channel.postMessage({ type: "SIGN_OUT", timestamp: Date.now() })
}
```

### Public API

```tsx
// features/auth/index.ts
export { LoginForm } from "./ui/login-form"
export { RegisterForm } from "./ui/register-form"
export { useAuthSync, broadcastSignOut, broadcastSignIn } from "./model/use-auth-sync"
```

### Integration in Header Widget

```tsx
// widgets/header/ui/header.tsx
import { broadcastSignOut, useAuthSync } from "@features/auth"
import { signOut } from "@shared/lib/auth-client"

export function Header({ user }: { user: SessionUser }) {
  useAuthSync() // Cross-Tab Listener

  const handleSignOut = async () => {
    await signOut()
    await broadcastSignOut()
    router.push("/")
  }
  // ...
}
```

### Benefits of broadcast-channel

- Safari Private Mode Support
- localStorage Fallback
- TypeScript Support
- WebWorker Support

## Route Protection

### Proxy (middleware)

`frontend/src/proxy.ts` checks cookie-based:

```typescript
export function proxy(request: NextRequest) {
  const sessionCookie = request.cookies.get("better-auth.session_token")

  // Protected Route without Session → Login
  if (isProtectedRoute && !sessionCookie) {
    return NextResponse.redirect(new URL("/login", request.url))
  }

  // Auth Route with Session → Dashboard
  if (isAuthRoute && sessionCookie) {
    return NextResponse.redirect(new URL("/dashboard", request.url))
  }
}
```

### Folder Structure

```
src/app/
├── (auth)/           # Login, Register (redirect when logged in)
│   ├── login/
│   └── register/
├── (protected)/      # Requires Auth
│   └── dashboard/
│       ├── page.tsx           # Server Component
│       ├── header.tsx         # Client Component
│       └── dashboard-content.tsx
└── page.tsx          # Public Home
```

## IMPORTANT: What NOT to do

❌ **DO NOT** use `useSession()` in Protected Pages
❌ **DO NOT** use Skeleton for Session Loading
❌ **DO NOT** use `isPending` check for Session
❌ **DO NOT** use client-side redirect

✅ **INSTEAD**:

- Server Component checks session
- `redirect()` if no session
- Pass user as prop to Client Components
- Use `initialData` for React Query

## Real-Time Updates

SSE Hook for automatic data updates:

```tsx
"use client"

import { useSSE } from "@features/stats"

export function MyComponent() {
  useSSE() // Connects to /api/v1/events
  // React Query is automatically invalidated on SSE events
}
```
