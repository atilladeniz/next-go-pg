---
name: api-generation
description: Generate TypeScript API client from Swagger/Go comments. Use when updating API endpoints, adding new routes, or regenerating the frontend API client after backend changes.
allowed-tools: Read, Edit, Write, Bash, Glob, Grep
---

# API Generation

Generate typed TypeScript API client from Go Swagger comments using swag + Orval.

## Workflow

```
Go Handler → swag init → swagger.json → Orval → TypeScript Hooks
```

## Ein Befehl für alles

```bash
make api
```

Dies führt aus:
1. `swag init` → Generiert `backend/docs/swagger.json` aus Go-Kommentaren
2. `orval` → Generiert TypeScript Hooks in `frontend/src/api/`

## Neuen Endpoint hinzufügen

### Schritt 1: Go Handler mit Swagger-Kommentaren

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

// CreateProduct godoc
// @Summary Create a product
// @Description Create a new product
// @Tags products
// @Accept json
// @Produce json
// @Param request body CreateProductRequest true "Product data"
// @Success 201 {object} ProductResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Security BearerAuth
// @Router /products [post]
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
    // Implementation
}
```

### Schritt 2: Response/Request Types definieren

```go
// In handler oder separater types.go

type ProductResponse struct {
    ID    string  `json:"id" example:"prod_123"`
    Name  string  `json:"name" example:"Widget"`
    Price float64 `json:"price" example:"29.99"`
}

type CreateProductRequest struct {
    Name  string  `json:"name" example:"Widget"`
    Price float64 `json:"price" example:"29.99"`
}
```

### Schritt 3: API Client generieren

```bash
make api
```

### Schritt 4: Im Frontend nutzen

```tsx
import { useGetProducts, usePostProducts } from "@/api/endpoints/products/products"

function ProductsPage() {
  const { data, isLoading } = useGetProducts()
  const createProduct = usePostProducts()

  const handleCreate = () => {
    createProduct.mutate({ data: { name: "New", price: 10 } })
  }

  return (...)
}
```

## Swagger Annotation Reference

| Annotation | Beschreibung |
|------------|--------------|
| `@Summary` | Kurze Beschreibung |
| `@Description` | Ausführliche Beschreibung |
| `@Tags` | Gruppierung (wird zu Ordner in endpoints/) |
| `@Accept` | Request Content-Type |
| `@Produce` | Response Content-Type |
| `@Param` | Parameter (body, query, path) |
| `@Success` | Erfolgs-Response |
| `@Failure` | Fehler-Response |
| `@Security` | Auth-Anforderung |
| `@Router` | HTTP Path und Method |

## Generierte Dateien

```
frontend/src/api/
├── endpoints/
│   ├── products/        # Nach @Tags gruppiert
│   │   └── products.ts  # useGetProducts, usePostProducts, etc.
│   └── users/
│       └── users.ts
├── models/              # TypeScript Types
└── custom-fetch.ts      # Fetch Wrapper mit Auth
```

## Wichtige Regeln

- **Nie** generierte Dateien manuell editieren
- **Immer** `make api` nach Handler-Änderungen
- Tags werden zu Ordnernamen → `@Tags products` → `endpoints/products/`
- operationId wird automatisch aus Router-Path generiert
