---
name: feature-generator
description: Generate full-stack features with Goca backend and React frontend. Use when creating new features, adding CRUD operations, or scaffolding new pages.
allowed-tools: Read, Edit, Write, Bash, Glob, Grep
---

# Feature Generator

Create complete full-stack features following Clean Architecture with Goca.

## Feature Structure

### Backend (Go mit Goca)

1. **Entity** in `backend/internal/domain/` → `goca make entity`
2. **Repository** in `backend/internal/repository/` → `goca make repository`
3. **UseCase** in `backend/internal/usecase/` → `goca make usecase`
4. **Handler** in `backend/internal/handler/` → `goca make handler`

### Frontend (React)

1. **Server Component** für initiales Laden (kein Flicker)
2. **Client Component** für Interaktivität
3. **Generated API hooks** via Orval
4. **UI components** mit shadcn/ui

## Schnellstart: Neues Feature mit Goca

```bash
cd backend

# Komplettes Feature generieren
goca feature Product --fields "name:string,price:float64,stock:int"

# Oder einzelne Layer
goca make entity Product
goca make repository Product
goca make usecase Product
goca make handler Product

# API Client generieren
cd ..
make api
```

## Entity Registry (AutoMigrate)

Nach `goca feature` muss die neue Entity in `backend/internal/domain/registry.go` registriert werden:

```go
// internal/domain/registry.go
func AllEntities() []interface{} {
    return []interface{}{
        &UserStats{},
        &Product{},  // ← Neue Entity hier hinzufuegen!
    }
}
```

Das ist die **EINZIGE** Stelle wo neue Entities registriert werden muessen!

## Kompletter Workflow

```bash
# 1. Feature generieren
cd backend
goca feature Product --fields "name:string,price:float64,stock:int"

# 2. Entity in Registry hinzufuegen
# backend/internal/domain/registry.go → &Product{} hinzufuegen

# 3. API generieren
cd ..
make api

# 4. Backend neu starten (Migration laeuft automatisch)
make dev-backend
```

## Server-Side Data Loading Pattern (HydrationBoundary)

### Schritt 1: Server Component mit prefetchQuery

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
  // 1. Session prüfen
  const session = await getSession()
  if (!session) redirect("/login")

  // 2. Cookies für Auth
  const cookieStore = await cookies()
  const cookieHeader = cookieStore.getAll().map((c) => `${c.name}=${c.value}`).join("; ")

  // 3. Prefetch mit Orval-Funktion
  const queryClient = getQueryClient()
  await queryClient.prefetchQuery({
    queryKey: getGetProductsQueryKey(),
    queryFn: () => getProducts({ headers: { Cookie: cookieHeader }, cache: "no-store" }),
  })

  // 4. HydrationBoundary wrappen
  return (
    <HydrationBoundary state={dehydrate(queryClient)}>
      <div className="container py-8">
        <h1 className="text-2xl font-bold mb-4">Produkte</h1>
        <ProductList />
      </div>
    </HydrationBoundary>
  )
}
```

### Schritt 2: Client Component (kein initialData nötig!)

```tsx
// app/(protected)/products/product-list.tsx
"use client"

import { useGetProducts } from "@/api/endpoints/products/products"
import { useSSE } from "@/hooks/use-sse"

export function ProductList() {
  useSSE() // Real-time Updates

  // Daten sind bereits hydriert - kein initialData nötig!
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
// internal/handler/product.go

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

## Nach Backend-Änderungen

```bash
# Immer ausführen nach Handler-Änderungen:
make api
```

## Konventionen

### Datei-Benennung

- Go: `snake_case.go`
- React Pages: `page.tsx` im Route-Ordner
- Client Components: `kebab-case.tsx`

### Route Protection

- Public: `frontend/src/app/`
- Protected: `frontend/src/app/(protected)/`
- Auth: `frontend/src/app/(auth)/`

### Kein Skeleton/Flicker

- Server Component lädt Daten vor dem Rendern
- Client Component erhält `initialData`
- React Query übernimmt für Updates
- SSE für Real-time Sync
