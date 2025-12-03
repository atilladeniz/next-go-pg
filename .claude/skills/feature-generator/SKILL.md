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

## Server-Side Data Loading Pattern (KEIN FLICKER!)

### Schritt 1: Server Component für Daten

```tsx
// app/(protected)/products/page.tsx (Server Component)
import { auth } from "@/lib/auth"
import { headers } from "next/headers"
import { redirect } from "next/navigation"
import { ProductList } from "./product-list"

async function getProducts() {
  const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/v1/products`, {
    headers: { Cookie: (await headers()).get("cookie") ?? "" },
    cache: "no-store",
  })
  if (!res.ok) return null
  return res.json()
}

export default async function ProductsPage() {
  const session = await auth.api.getSession({ headers: await headers() })
  if (!session) redirect("/login")

  const products = await getProducts()

  return (
    <div className="container py-8">
      <h1 className="text-2xl font-bold mb-4">Produkte</h1>
      <ProductList initialProducts={products} />
    </div>
  )
}
```

### Schritt 2: Client Component mit initialData

```tsx
// app/(protected)/products/product-list.tsx
"use client"

import { useGetProducts } from "@/api/endpoints/products/products"
import { useSSE } from "@/hooks/use-sse"

type Product = { id: string; name: string; price: number }

export function ProductList({ initialProducts }: { initialProducts: Product[] | null }) {
  // SSE für Real-time Updates
  useSSE()

  // React Query mit Server-Daten als Initial
  const { data: productsResponse } = useGetProducts({
    query: {
      initialData: initialProducts
        ? { data: initialProducts, status: 200 as const, headers: new Headers() }
        : undefined,
    },
  })

  const products = productsResponse?.status === 200 ? productsResponse.data : initialProducts

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
