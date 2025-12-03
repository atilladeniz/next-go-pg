---
name: data-fetching
description: Server-Side + Client-Side Data Fetching mit Orval + TanStack Query HydrationBoundary Pattern. IMMER Orval nutzen - NIEMALS manuelles fetch()!
allowed-tools: Read, Edit, Write, Glob, Grep
---

# Data Fetching Strategy

**Grundregel**: IMMER Orval-generierte Funktionen nutzen - NIEMALS manuelle `fetch()` Aufrufe!

## Das HydrationBoundary Pattern (TanStack Recommended)

```
Server Component (prefetchQuery) → HydrationBoundary → Client Component (useQuery)
```

### Vorteile gegenüber initialData

- Sauberer: Kein manuelles Response-Mapping
- Streaming-Support: Unterstützt React 18 Streaming
- Korrekter Cache: Query Cache wird korrekt hydriert
- Type-Safe: Bessere TypeScript-Integration

## Setup

### 1. Query Client Helper (`lib/get-query-client.ts`)

```tsx
import {
  isServer,
  QueryClient,
  defaultShouldDehydrateQuery,
} from "@tanstack/react-query"

function makeQueryClient() {
  return new QueryClient({
    defaultOptions: {
      queries: {
        staleTime: 60 * 1000,
      },
      dehydrate: {
        // Include pending queries for streaming
        shouldDehydrateQuery: (query) =>
          defaultShouldDehydrateQuery(query) ||
          query.state.status === "pending",
      },
    },
  })
}

let browserQueryClient: QueryClient | undefined = undefined

export function getQueryClient() {
  if (isServer) {
    return makeQueryClient()
  }
  if (!browserQueryClient) browserQueryClient = makeQueryClient()
  return browserQueryClient
}
```

### 2. Providers (`components/providers.tsx`)

```tsx
"use client"

import { QueryClientProvider } from "@tanstack/react-query"
import { getQueryClient } from "@/lib/get-query-client"

export function Providers({ children }: { children: ReactNode }) {
  const queryClient = getQueryClient()

  return (
    <QueryClientProvider client={queryClient}>
      {children}
    </QueryClientProvider>
  )
}
```

## Server Component Pattern

```tsx
// app/(protected)/dashboard/page.tsx - SERVER COMPONENT
import { dehydrate, HydrationBoundary } from "@tanstack/react-query"
import { cookies } from "next/headers"
import { redirect } from "next/navigation"
import { getStats, getGetStatsQueryKey } from "@/api/endpoints/users/users"
import { getQueryClient } from "@/lib/get-query-client"
import { getSession } from "@/lib/auth-server"
import { StatsGrid } from "./stats-grid"

export default async function DashboardPage() {
  // 1. Session prüfen
  const session = await getSession()
  if (!session) redirect("/login")

  // 2. Cookies für Server-Fetch
  const cookieStore = await cookies()
  const cookieHeader = cookieStore
    .getAll()
    .map((c) => `${c.name}=${c.value}`)
    .join("; ")

  // 3. Prefetch mit Orval-Funktion
  const queryClient = getQueryClient()
  await queryClient.prefetchQuery({
    queryKey: getGetStatsQueryKey(),
    queryFn: () =>
      getStats({
        headers: { Cookie: cookieHeader },
        cache: "no-store",
      }),
  })

  // 4. HydrationBoundary wrappen
  return (
    <HydrationBoundary state={dehydrate(queryClient)}>
      <StatsGrid />
    </HydrationBoundary>
  )
}
```

## Client Component Pattern

```tsx
// app/(protected)/dashboard/stats-grid.tsx - CLIENT COMPONENT
"use client"

import { useGetStats, usePostStats } from "@/api/endpoints/users/users"
import { useSSE } from "@/hooks/use-sse"

export function StatsGrid() {
  // SSE für Real-time Updates
  useSSE()

  // Daten sind bereits hydriert - kein initialData nötig!
  const { data: response } = useGetStats()

  // Mutation Hook
  const { mutate: updateStats } = usePostStats()

  const stats = response?.status === 200 ? response.data : null

  return (
    <div>
      <p>Projekte: {stats?.projectCount}</p>
      <button onClick={() => updateStats({ data: { field: "projects", delta: 1 } })}>
        +1
      </button>
    </div>
  )
}
```

## Multiple Datenquellen

```tsx
// Server Component
export default async function DashboardPage() {
  const session = await getSession()
  if (!session) redirect("/login")

  const cookieStore = await cookies()
  const cookieHeader = cookieStore
    .getAll()
    .map((c) => `${c.name}=${c.value}`)
    .join("; ")

  const queryClient = getQueryClient()

  // Parallel prefetchen
  await Promise.all([
    queryClient.prefetchQuery({
      queryKey: getGetStatsQueryKey(),
      queryFn: () => getStats({ headers: { Cookie: cookieHeader }, cache: "no-store" }),
    }),
    queryClient.prefetchQuery({
      queryKey: getGetNotificationsQueryKey(),
      queryFn: () => getNotifications({ headers: { Cookie: cookieHeader }, cache: "no-store" }),
    }),
  ])

  return (
    <HydrationBoundary state={dehydrate(queryClient)}>
      <StatsGrid />
      <NotificationList />
    </HydrationBoundary>
  )
}
```

## VERBOTEN: Manuelle Fetch-Aufrufe

```tsx
// ❌ NIEMALS SO:
async function getStats() {
  const res = await fetch("http://localhost:8080/api/v1/stats")
  return res.json()
}

// ✅ IMMER SO (Orval-Funktion):
import { getStats, getGetStatsQueryKey } from "@/api/endpoints/users/users"

await queryClient.prefetchQuery({
  queryKey: getGetStatsQueryKey(),
  queryFn: () => getStats({ headers: { Cookie: cookieHeader } }),
})
```

## Wann Server-Side prefetchen?

✅ **Server-Side** (prefetchQuery in Server Component):
- Initial Page Load (SEO, kein Flicker)
- Protected Pages (Session prüfen vor Render)
- Kritische "above-the-fold" Inhalte
- Daten die sofort sichtbar sein müssen

## Wann nur Client-Side?

✅ **Nur Client-Side** (useQuery ohne prefetch):
- Nach User-Interaktion (Klick, Form Submit)
- Lazy-loaded Content (unterhalb des Folds)
- Pagination, Infinite Scroll
- Daten die nicht sofort sichtbar sein müssen

## SSE + React Query Integration

```tsx
// hooks/use-sse.ts
"use client"

import { useQueryClient } from "@tanstack/react-query"
import { useEffect } from "react"
import { getGetStatsQueryKey } from "@/api/endpoints/users/users"

export function useSSE() {
  const queryClient = useQueryClient()

  useEffect(() => {
    const eventSource = new EventSource(
      `${process.env.NEXT_PUBLIC_API_URL}/api/v1/events`
    )

    eventSource.addEventListener("stats-updated", () => {
      queryClient.invalidateQueries({ queryKey: getGetStatsQueryKey() })
    })

    return () => eventSource.close()
  }, [queryClient])
}
```

## Streaming (Optional)

Für Streaming ohne await:

```tsx
export default function PostsPage() {
  const queryClient = getQueryClient()

  // Kein await - startet Fetch, blockiert nicht
  queryClient.prefetchQuery({
    queryKey: getGetPostsQueryKey(),
    queryFn: () => getPosts(),
  })

  return (
    <HydrationBoundary state={dehydrate(queryClient)}>
      <Posts /> {/* useSuspenseQuery hier für Streaming */}
    </HydrationBoundary>
  )
}
```

## Zusammenfassung

```
┌─────────────────────────────────────────────────────────────┐
│                    SERVER COMPONENT                          │
│  1. Session prüfen (getSession)                             │
│  2. Cookies holen für Auth                                  │
│  3. prefetchQuery mit Orval-Funktion                        │
│  4. HydrationBoundary wrappen                               │
└─────────────────────────────────────────────────────────────┘
                              ↓
                     dehydrate(queryClient)
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                    CLIENT COMPONENT                          │
│  1. useQuery() - Daten sind bereits da!                     │
│  2. useSSE() für Real-time Updates                          │
│  3. useMutation() für Änderungen                            │
└─────────────────────────────────────────────────────────────┘
```
