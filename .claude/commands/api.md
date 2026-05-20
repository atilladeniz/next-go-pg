---
allowed-tools: Bash(make:*), Bash(cd:*)
description: Generate Swagger docs and TypeScript API client
---

# API Generation

Regenerate Swagger documentation and TypeScript API client.

## Command

```bash
just api
```

This runs:
1. `swag init` - Generate `backend/docs/swagger.json` from Go comments
2. `orval` - Generate TypeScript hooks in `frontend/src/api/`

## When to run

- After adding/modifying handler endpoints (e.g. in `backend/internal/<ctx>/interfaces/http/`)
- After changing Swagger comments (`@Summary`, `@Router`, etc.)

## Generated Files

- `backend/docs/swagger.json` - OpenAPI spec
- `frontend/src/api/endpoints/` - React Query hooks
- `frontend/src/api/models/` - TypeScript types
