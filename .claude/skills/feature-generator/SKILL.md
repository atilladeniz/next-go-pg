---
name: feature-generator
description: Generate full-stack features with Go backend handler and React frontend page. Use when creating new features, adding CRUD operations, or scaffolding new pages.
allowed-tools: Read, Edit, Write, Bash, Glob, Grep
---

# Feature Generator

Create complete full-stack features following project conventions.

## Feature Structure

A complete feature includes:

### Backend (Go)

1. **Handler** in `backend/internal/handler/`
2. **OpenAPI endpoint** in `backend/api/openapi.yaml`
3. **Middleware** if authentication required

### Frontend (React)

1. **Page** in `frontend/src/app/`
2. **Generated API hooks** via Orval
3. **UI components** using shadcn/ui

## Creating a New Feature

### Step 1: Add OpenAPI endpoint

Edit `backend/api/openapi.yaml`:

```yaml
paths:
  /features:
    get:
      operationId: getFeatures
      tags: [features]
      summary: List all features
      responses:
        "200":
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: "#/components/schemas/Feature"
    post:
      operationId: createFeature
      tags: [features]
      summary: Create a feature
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: "#/components/schemas/CreateFeatureRequest"
      responses:
        "201":
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/Feature"
```

### Step 2: Add Go handler

Create `backend/internal/handler/feature.go`:

```go
package handler

import (
    "encoding/json"
    "net/http"
)

type FeatureHandler struct {
    // dependencies
}

func NewFeatureHandler() *FeatureHandler {
    return &FeatureHandler{}
}

func (h *FeatureHandler) List(w http.ResponseWriter, r *http.Request) {
    // implementation
}

func (h *FeatureHandler) Create(w http.ResponseWriter, r *http.Request) {
    // implementation
}
```

### Step 3: Register routes

In `backend/internal/handler/api.go`:

```go
featureHandler := NewFeatureHandler()
r.Route("/features", func(r chi.Router) {
    r.Get("/", featureHandler.List)
    r.Post("/", featureHandler.Create)
})
```

### Step 4: Generate API client

```bash
make api
```

### Step 5: Create frontend page

Create `frontend/src/app/(protected)/features/page.tsx`:

```tsx
"use client"

import { useGetFeatures, useCreateFeature } from "@/api/endpoints/features/features"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"

export default function FeaturesPage() {
    const { data: features, isLoading } = useGetFeatures()
    const createFeature = useCreateFeature()

    if (isLoading) return <div>Laden...</div>

    return (
        <div className="container py-8">
            <h1 className="text-2xl font-bold mb-4">Features</h1>
            <div className="grid gap-4">
                {features?.map(feature => (
                    <Card key={feature.id}>
                        <CardHeader>
                            <CardTitle>{feature.name}</CardTitle>
                        </CardHeader>
                    </Card>
                ))}
            </div>
        </div>
    )
}
```

## Conventions

### File Naming

- Go handlers: `snake_case.go`
- React pages: `page.tsx` in route folder
- Components: `kebab-case.tsx`

### Route Protection

- Public routes: `frontend/src/app/`
- Protected routes: `frontend/src/app/(protected)/`
- Auth routes: `frontend/src/app/(auth)/`

### UI Components

Always use shadcn/ui components:

```tsx
import { Button } from "@/components/ui/button"
import { Card } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
```
