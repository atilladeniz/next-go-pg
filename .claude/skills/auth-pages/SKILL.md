---
name: auth-pages
description: Create and manage authentication pages and flows. Use when adding login, register, password reset, or other auth-related pages.
allowed-tools: Read, Edit, Write, Glob
---

# Authentication Pages

Better Auth integration with Next.js 16 for authentication flows.

## Existing Auth Pages

- `/login` - Login with email/password
- `/register` - User registration
- `/dashboard` - Protected dashboard (requires auth)

## Auth Client

```tsx
import { signIn, signUp, signOut, useSession } from "@/lib/auth-client"
```

### Check Session

```tsx
const { data: session, isPending } = useSession()

if (isPending) return <div>Laden...</div>
if (!session) return <div>Nicht angemeldet</div>

// Access user data
session.user.name
session.user.email
```

### Sign In

```tsx
const result = await signIn.email({
  email: "user@example.com",
  password: "password123",
})

if (result.error) {
  // Handle error
  console.error(result.error.message)
} else {
  // Redirect to dashboard
  router.push("/dashboard")
}
```

### Sign Up

```tsx
const result = await signUp.email({
  name: "John Doe",
  email: "user@example.com",
  password: "password123",
})

if (result.error) {
  console.error(result.error.message)
} else {
  router.push("/dashboard")
}
```

### Sign Out

```tsx
await signOut()
router.push("/")
```

## Route Protection

### Proxy (middleware)

Located at `frontend/src/proxy.ts`:

```typescript
export function proxy(request: NextRequest) {
  const sessionCookie = request.cookies.get("better-auth.session_token")

  // Redirect to login if accessing protected route without session
  if (isProtectedRoute && !sessionCookie) {
    return NextResponse.redirect(new URL("/login", request.url))
  }

  // Redirect to dashboard if accessing auth routes with active session
  if (isAuthRoute && sessionCookie) {
    return NextResponse.redirect(new URL("/dashboard", request.url))
  }
}
```

### Protected Routes

Place protected pages in `frontend/src/app/(protected)/`:

```
src/app/
├── (auth)/           # Auth pages (login, register)
│   ├── login/
│   └── register/
├── (protected)/      # Requires authentication
│   └── dashboard/
└── page.tsx          # Public home page
```

## Creating New Auth Page

### Example: Password Reset Page

Create `frontend/src/app/(auth)/reset-password/page.tsx`:

```tsx
"use client"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"

export default function ResetPasswordPage() {
  const [email, setEmail] = useState("")
  const [loading, setLoading] = useState(false)
  const [sent, setSent] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)

    // TODO: Implement password reset
    // await authClient.forgetPassword({ email })

    setSent(true)
    setLoading(false)
  }

  if (sent) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <Card className="w-full max-w-md">
          <CardContent className="pt-6">
            <p>Prüfe deine E-Mails für den Reset-Link.</p>
          </CardContent>
        </Card>
      </div>
    )
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-background">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle>Passwort zurücksetzen</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="email">E-Mail</Label>
              <Input
                id="email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
              />
            </div>
            <Button type="submit" className="w-full" disabled={loading}>
              {loading ? "Wird gesendet..." : "Reset-Link senden"}
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
```

## UI Conventions

- Use german labels (Anmelden, Registrieren, Abmelden)
- Use shadcn/ui Card for auth forms
- Show loading state on buttons
- Display errors in red alert box
- Redirect after successful auth
