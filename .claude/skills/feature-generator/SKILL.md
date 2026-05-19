---
name: feature-generator
description: Generate full-stack features with Goca backend and React frontend. Use when creating new features, adding CRUD operations, or scaffolding new pages.
allowed-tools: Read, Edit, Write, Bash, Glob, Grep
---

# Feature Generator

Create complete full-stack features. Backend follows **bounded contexts** (DDD-strategic) with Clean-Architecture layers (DDD-tactical) inside each context. Frontend follows Feature-Sliced Design.

## Feature Structure

### Backend (Go) — bounded-context layout

Each aggregate lives inside one bounded context (e.g. `stats/`, `auth/`, `notifications/`, `exports/`, or a brand-new one). Goca scaffolds the four layers into `internal/_goca_inbox/` (ignored by `go build`); you then **relocate** the output into the target context.

1. **Domain** → `goca make entity` → move into `backend/internal/<ctx>/domain/` (strip GORM tags; embed `shared.AggregateBase` if it raises events; add value-object constructors)
2. **Application port + use case** → `goca make usecase` → split into `backend/internal/<ctx>/application/ports.go` and `<ctx>/application/<aggregate>_usecases.go`
3. **Infrastructure (persistence)** → `goca make repository` → split into `backend/internal/<ctx>/infrastructure/persistence/{gorm_models.go, <aggregate>_mapper.go, <aggregate>_repo.go, registry.go}`. Add the compile-time assertion `var _ <ctx>app.<Aggregate>Repository = (*Repository)(nil)`.
4. **HTTP adapter** → `goca make handler` → move into `backend/internal/<ctx>/interfaces/http/handler.go`. Imports ONLY this context's `application/`.
5. **Wire in `internal/composition/composition.go`** — build the repository, use cases, handler; append `Entities()` to `runAutoMigrations`. If the context needs data from another context, write an Anti-Corruption Layer in `composition.go` (mirror `statsToExportsReader` / `authToNotificationsDirectory`).

### Frontend (React)

1. **Server Component** for initial loading (no flicker)
2. **Client Component** for interactivity
3. **Generated API hooks** via Orval
4. **UI components** with shadcn/ui

## Quick Start: New Feature with Goca

```bash
cd backend

# Komplettes Feature generieren
goca feature Product --fields "name:string,price:float64,stock:int"

# Or individual layers
goca make entity Product
goca make repository Product
goca make usecase Product
goca make handler Product

# Generate API Client
cd ..
just api
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

## Complete Workflow

```bash
# 1. Scaffold into _goca_inbox/
cd backend
goca feature Product --fields "name:string,price:float64,stock:int"

# 2. Relocate the four generated files into a bounded context:
#    - domain  → internal/<ctx>/domain/
#    - usecase → internal/<ctx>/application/
#    - repo    → internal/<ctx>/infrastructure/persistence/
#    - handler → internal/<ctx>/interfaces/http/
#    Strip GORM tags from the domain type, write a GORM twin + mapper,
#    add Entities() to the context's persistence/registry.go.

# 3. Wire in internal/composition/composition.go
#    (repo, usecase, handler; append Entities() to AutoMigrate)

# 4. Generate Swagger + Orval API client
cd ..
just api

# 5. Restart backend (AutoMigrate runs on startup)
just dev-backend
```

## Server-Side Data Loading Pattern (HydrationBoundary)

### Step 1: Server Component with prefetchQuery

```tsx
// app/(protected)/products/page.tsx (Server Component)
import { dehydrate, HydrationBoundary } from "@tanstack/react-query"
import { cookies } from "next/headers"
import { redirect } from "next/navigation"
import { getGetProductsQueryKey, getProducts } from "@/api/endpoints/products/products"
import { getQueryClient } from "@/lib/get-query-client"
import { getSession } from "@/lib/auth-server"
import { ProductList } from "./product-list"

export default async function ProductsPage() {
  // 1. Check session
  const session = await getSession()
  if (!session) redirect("/login")

  // 2. Cookies for auth
  const cookieStore = await cookies()
  const cookieHeader = cookieStore.getAll().map((c) => `${c.name}=${c.value}`).join("; ")

  // 3. Prefetch with Orval function
  const queryClient = getQueryClient()
  await queryClient.prefetchQuery({
    queryKey: getGetProductsQueryKey(),
    queryFn: () => getProducts({ headers: { Cookie: cookieHeader }, cache: "no-store" }),
  })

  // 4. Wrap with HydrationBoundary
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

### Step 2: Client Component (no initialData needed!)

```tsx
// app/(protected)/products/product-list.tsx
"use client"

import { useGetProducts } from "@/api/endpoints/products/products"
import { useSSE } from "@/hooks/use-sse"

export function ProductList() {
  useSSE() // Real-time Updates

  // Data is already hydrated - no initialData needed!
  const { data: productsResponse } = useGetProducts()

  const products = productsResponse?.status === 200 ? productsResponse.data : null

  return (
    <div className="grid gap-4">
      {products?.map(product => (
        <div key={product.id}>{product.name} - {product.price}€</div>
      ))}
    </div>
  )
}
```

## Backend Handler mit Swagger

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
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
    // Implementation
}
```

## After Backend Changes

```bash
# Always run after handler changes:
just api
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
- Client Component receives `initialData`
- React Query takes over for updates
- SSE for real-time sync
