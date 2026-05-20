---
name: feature-generator
description: Generate full-stack features. Backend = hand-written bounded-context aggregates (DDD); frontend = FSD slices with HydrationBoundary. Use when creating new features, adding CRUD operations, or scaffolding new pages.
allowed-tools: Read, Edit, Write, Bash, Glob, Grep
---

# Feature Generator

Create complete full-stack features. Backend follows **bounded contexts** (DDD-strategic) with Clean-Architecture layers (DDD-tactical) inside each context. Frontend follows Feature-Sliced Design.

**No code generator is used.** The backend used to be scaffolded with Goca; it was removed in `refactor/backend-clean-architecture` because its flat-layer output is structurally incompatible with this repo's bounded-context layout. The five backend files per aggregate are short — copy the existing `internal/stats/` context as the canonical template and modify.

## Backend (Go) — bounded-context layout

Each aggregate lives inside one bounded context (e.g. `stats/`, `auth/`, `notifications/`, `exports/`, or a brand-new one). Five files plus a wiring step:

1. **Domain** — `backend/internal/<ctx>/domain/<aggregate>.go`
   - Pure types only. No `gorm.io/gorm` import. No I/O.
   - Embed `shared.AggregateBase` if it raises events.
   - Define value objects with constructor invariants (`func NewMoney(...) (Money, error)`).
   - Define domain events implementing `EventName() string`.
2. **Application port + use case** — `backend/internal/<ctx>/application/`
   - `ports.go` declares interfaces (`Repository`, `JobEnqueuer`, ...).
   - `<aggregate>_usecases.go` holds use-case structs with `Execute(ctx, ...)`.
   - Pull events with `agg.PullEvents()` **before** `repo.Save(...)`.
3. **Persistence** — `backend/internal/<ctx>/infrastructure/persistence/`
   - `gorm_models.go` (unexported GORM-tagged twin) + `<aggregate>_mapper.go` + `<aggregate>_repo.go` + `registry.go` exposing `Entities() []any`.
   - Assert the port: `var _ <ctx>app.Repository = (*Repository)(nil)`.
   - `Save` must mutate only DB-owned fields back into `*agg`, never replace it whole (would wipe pending events).
4. **HTTP adapter** — `backend/internal/<ctx>/interfaces/http/handler.go`
   - Depends only on this context's `application/` package. Never imports gorm or another bounded context.
   - Swagger annotations on every endpoint.
5. **Wire in composition root** — `backend/internal/composition/composition.go`
   - Build repo → use cases → handler. Register routes. Append `<ctx>persist.Entities()` to `runAutoMigrations`.
   - If cross-context data is needed, add an Anti-Corruption Layer adapter right here (mirror `statsToExportsReader` / `authToNotificationsDirectory`).

## Frontend (React)

1. **Server Component** for initial loading (no flicker)
2. **Client Component** for interactivity
3. **Generated API hooks** via Orval
4. **UI components** with shadcn/ui

## Quick Start: Copy the Stats Context

```bash
# Look at the canonical example
ls backend/internal/stats/
# domain/  application/  infrastructure/  interfaces/

# Create your context by mirroring its layout
mkdir -p backend/internal/products/{domain,application,infrastructure/persistence,interfaces/http}
# Write your five files, then:
cd .. && just api   # regenerate Swagger + Orval TS client
just dev-backend    # AutoMigrate runs on startup
```

## Entity Registry (AutoMigrate)

There is **no central registry**. Each bounded context owns its own `Entities()` function under `internal/<ctx>/infrastructure/persistence/registry.go`. The composition root aggregates them:

```go
// internal/<ctx>/infrastructure/persistence/registry.go
func Entities() []any {
    return []any{&gormProduct{}}
}

// internal/composition/composition.go (runAutoMigrations)
entities = append(entities, productspersist.Entities()...)
```

## Server-Side Data Loading Pattern (HydrationBoundary)

### Step 1: Server Component with prefetchQuery

```tsx
// app/(protected)/products/page.tsx (Server Component)
import { dehydrate, HydrationBoundary } from "@tanstack/react-query"
import { cookies } from "next/headers"
import { redirect } from "next/navigation"
import { getGetProductsQueryKey, getProducts } from "@shared/api/endpoints/products/products"
import { getQueryClient } from "@shared/lib/query-client"
import { getSession } from "@shared/lib/auth-server"
import { ProductList } from "./product-list"

export default async function ProductsPage() {
  const session = await getSession()
  if (!session) redirect("/login")

  const cookieStore = await cookies()
  const cookieHeader = cookieStore.getAll().map((c) => `${c.name}=${c.value}`).join("; ")

  const queryClient = getQueryClient()
  await queryClient.prefetchQuery({
    queryKey: getGetProductsQueryKey(),
    queryFn: () => getProducts({ headers: { Cookie: cookieHeader }, cache: "no-store" }),
  })

  return (
    <HydrationBoundary state={dehydrate(queryClient)}>
      <div className="container py-8">
        <h1 className="text-2xl font-bold mb-4">Products</h1>
        <ProductList />
      </div>
    </HydrationBoundary>
  )
}
```

### Step 2: Client Component (no initialData needed)

```tsx
// app/(protected)/products/product-list.tsx
"use client"

import { useGetProducts } from "@shared/api/endpoints/products/products"
import { useSSE } from "@features/stats"

export function ProductList() {
  useSSE()
  const { data: productsResponse } = useGetProducts()
  const products = productsResponse?.status === 200 ? productsResponse.data : null

  return (
    <div className="grid gap-4">
      {products?.map((product) => (
        <div key={product.id}>{product.name} - {product.price}€</div>
      ))}
    </div>
  )
}
```

## Backend Handler with Swagger

```go
// internal/<ctx>/interfaces/http/handler.go

// GetProducts godoc
// @Summary List all products
// @Description Get all products for the authenticated user
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {array} ProductResponse
// @Failure 401 {object} ErrorResponse
// @Security BearerAuth
// @Router /products [get]
func (h *Handler) GetProducts(w http.ResponseWriter, r *http.Request) {
    // Implementation depends only on this context's application use cases.
}
```

## After Backend Changes

```bash
just api    # always run after handler changes
```

## Conventions

### File Naming

- Go: `snake_case.go`
- React Pages: `page.tsx` in route folder
- Client Components: `kebab-case.tsx`

### Route Protection

- Public: `frontend/src/app/`
- Protected: `frontend/src/app/(protected)/`
- Auth: `frontend/src/app/(auth)/`

### No Skeleton/Flicker

- Server Component loads data before rendering
- Client Component reads from the hydrated cache (no `initialData` prop needed)
- React Query takes over for updates
- SSE for real-time sync
