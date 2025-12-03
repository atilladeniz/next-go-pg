---
name: api-generation
description: Generate TypeScript API client from OpenAPI spec. Use when updating API endpoints, adding new routes, or regenerating the frontend API client after backend changes.
allowed-tools: Read, Edit, Write, Bash, Glob, Grep
---

# API Generation

Generate typed TypeScript API client from OpenAPI specification using Orval.

## Workflow

1. Edit the OpenAPI spec at `backend/api/openapi.yaml`
2. Run `make api` to generate TypeScript client
3. Use generated hooks in frontend components

## OpenAPI Spec Location

```
backend/api/openapi.yaml
```

## Adding a New Endpoint

### Step 1: Add to OpenAPI spec

```yaml
paths:
  /your-endpoint:
    get:
      operationId: getYourEndpoint
      tags: [your-tag]
      summary: Description of endpoint
      responses:
        "200":
          description: Success
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/YourResponse"
```

### Step 2: Add schema if needed

```yaml
components:
  schemas:
    YourResponse:
      type: object
      required:
        - id
        - name
      properties:
        id:
          type: string
        name:
          type: string
```

### Step 3: Generate client

```bash
make api
```

### Step 4: Use in frontend

```tsx
import { useGetYourEndpoint } from "@/api/endpoints/your-tag/your-tag"

function MyComponent() {
  const { data, isLoading, error } = useGetYourEndpoint()
  // ...
}
```

## Generated Files

- `frontend/src/api/endpoints/` - React Query hooks by tag
- `frontend/src/api/models/` - TypeScript types

## Important Notes

- Never edit generated files manually
- Always run `make api` after changing `openapi.yaml`
- Use tags to organize endpoints into separate files
- operationId becomes the hook name (camelCase with use prefix)
