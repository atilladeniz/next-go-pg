---
allowed-tools: Bash(make:*), Bash(cd:*)
description: Generate Swagger docs and TypeScript API client
---

# API Generation

Regenerate Swagger documentation and TypeScript API client.

## Command

```bash
make api
```

This runs:
1. `swag init` - Generate `backend/docs/swagger.json` from Go comments
2. `orval` - Generate TypeScript hooks in `frontend/src/api/`

## When to run

- After adding/modifying Handler endpoints
- After changing Swagger comments (`@Summary`, `@Router`, etc.)
- After running `goca feature` or `goca make handler`

## Generated Files

- `backend/docs/swagger.json` - OpenAPI spec
- `frontend/src/api/endpoints/` - React Query hooks
- `frontend/src/api/models/` - TypeScript types
