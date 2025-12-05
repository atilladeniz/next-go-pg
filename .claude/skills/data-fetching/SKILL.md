---
name: data-fetching
description: Server-Side + Client-Side Data Fetching with Orval + TanStack Query HydrationBoundary Pattern. ALWAYS use Orval - NEVER manual fetch()!
allowed-tools: Read, Edit, Write, Glob, Grep
---

# Data Fetching Strategy (FSD)

**Core Rule**: ALWAYS use Orval-generated functions - NEVER manual `fetch()` calls!

## FSD Paths

```
src/shared/
├── api/                        # Orval-generated
│   ├── endpoints/              # React Query Hooks
│   ├── models/                 # TypeScript Types
│   └── custom-fetch.ts         # Fetch Wrapper
└── lib/
    ├── query-client.ts         # getQueryClient()
    ├── auth-server/            # Server-only: getSession()
    └── auth-client/            # Client-safe: signIn, signOut
```

## The HydrationBoundary Pattern (TanStack Recommended)

```
Server Component (prefetchQuery) → HydrationBoundary → Client Component (useQuery)
```

### Advantages over initialData

- Cleaner: No manual response mapping
- Streaming Support: Supports React 18 Streaming
- Correct Cache: Query cache is properly hydrated
- Type-Safe: Better TypeScript integration

## Setup

### 1. Query Client Helper (`shared/lib/query-client.ts`)

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

### 2. Providers (`shared/config/providers.tsx`)

```tsx
"use client"

import { QueryClientProvider } from "@tanstack/react-query"
import { getQueryClient } from "@shared/lib/query-client"

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
import { getStats, getGetStatsQueryKey } from "@shared/api/endpoints/users/users"
import { getQueryClient } from "@shared/lib/query-client"
import { getSession } from "@shared/lib/auth-server"
import { StatsGrid } from "@features/stats"

export default async function DashboardPage() {
  // 1. Check session
  const session = await getSession()
  if (!session) redirect("/login")

  // 2. Get cookies for server fetch
  const cookieStore = await cookies()
  const cookieHeader = cookieStore
    .getAll()
    .map((c) => `${c.name}=${c.value}`)
    .join("; ")

  // 3. Prefetch with Orval function
  const queryClient = getQueryClient()
  await queryClient.prefetchQuery({
    queryKey: getGetStatsQueryKey(),
    queryFn: () =>
      getStats({
        headers: { Cookie: cookieHeader },
        cache: "no-store",
      }),
  })

  // 4. Wrap with HydrationBoundary
  return (
    <HydrationBoundary state={dehydrate(queryClient)}>
      <StatsGrid />
    </HydrationBoundary>
  )
}
```

## Client Component Pattern

```tsx
// features/stats/ui/stats-grid.tsx - CLIENT COMPONENT
"use client"

import { useGetStats, usePostStats } from "@shared/api/endpoints/users/users"
import { useSSE } from "../model/use-sse"

export function StatsGrid() {
  // SSE for real-time updates
  useSSE()

  // Data is already hydrated - no initialData needed!
  const { data: response } = useGetStats()

  // Mutation Hook
  const { mutate: updateStats } = usePostStats()

  const stats = response?.status === 200 ? response.data : null

  return (
    <div>
      <p>Projects: {stats?.projectCount}</p>
      <button onClick={() => updateStats({ data: { field: "projects", delta: 1 } })}>
        +1
      </button>
    </div>
  )
}
```

## Multiple Data Sources

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

  // Parallel prefetch
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

## FORBIDDEN: Manual Fetch Calls

```tsx
// ❌ NEVER DO THIS:
async function getStats() {
  const res = await fetch("http://localhost:8080/api/v1/stats")
  return res.json()
}

// ✅ ALWAYS DO THIS (Orval function):
import { getStats, getGetStatsQueryKey } from "@shared/api/endpoints/users/users"

await queryClient.prefetchQuery({
  queryKey: getGetStatsQueryKey(),
  queryFn: () => getStats({ headers: { Cookie: cookieHeader } }),
})
```

## When to Server-Side Prefetch?

✅ **Server-Side** (prefetchQuery in Server Component):

- Initial Page Load (SEO, no flicker)
- Protected Pages (check session before render)
- Critical "above-the-fold" content
- Data that must be immediately visible

## When Client-Side Only?

✅ **Client-Side Only** (useQuery without prefetch):

- After user interaction (click, form submit)
- Lazy-loaded content (below the fold)
- Pagination, Infinite Scroll
- Data that doesn't need to be immediately visible

## SSE + React Query Integration

```tsx
// features/stats/model/use-sse.ts
"use client"

import { useQueryClient } from "@tanstack/react-query"
import { useEffect } from "react"
import { getGetStatsQueryKey } from "@shared/api/endpoints/users/users"

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

For streaming without await:

```tsx
export default function PostsPage() {
  const queryClient = getQueryClient()

  // No await - starts fetch, doesn't block
  queryClient.prefetchQuery({
    queryKey: getGetPostsQueryKey(),
    queryFn: () => getPosts(),
  })

  return (
    <HydrationBoundary state={dehydrate(queryClient)}>
      <Posts /> {/* useSuspenseQuery here for streaming */}
    </HydrationBoundary>
  )
}
```

## Summary

```
┌─────────────────────────────────────────────────────────────┐
│                    SERVER COMPONENT                          │
│  1. Check session (getSession)                              │
│  2. Get cookies for auth                                    │
│  3. prefetchQuery with Orval function                       │
│  4. Wrap with HydrationBoundary                             │
└─────────────────────────────────────────────────────────────┘
                              ↓
                     dehydrate(queryClient)
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                    CLIENT COMPONENT                          │
│  1. useQuery() - Data is already there!                     │
│  2. useSSE() for real-time updates                          │
│  3. useMutation() for changes                               │
└─────────────────────────────────────────────────────────────┘
```
