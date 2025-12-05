---
allowed-tools: Bash(docker:*), Bash(make:*), Bash(cd:*)
description: Run database migrations (Better Auth + GORM AutoMigrate)
---

# Database Migration

## Better Auth Migration

Run Better Auth schema migration:

```bash
cd frontend && DATABASE_URL="postgres://postgres:postgres@localhost:5432/nextgopg" bunx @better-auth/cli migrate -y
```

## GORM AutoMigrate

GORM migrations run automatically on backend startup. To add a new entity:

1. Add to `backend/internal/domain/registry.go`:
   ```go
   func AllEntities() []interface{} {
       return []interface{}{
           &UserStats{},
           &NewEntity{},  // Add here
       }
   }
   ```

2. Restart backend:
   ```bash
   make dev-backend
   ```

## Database Commands

- `make db-up` - Start PostgreSQL
- `make db-down` - Stop PostgreSQL
- `make db-reset` - Reset database (delete all data)

## Check Tables

```bash
docker exec nextgopg-db-1 psql -U postgres -d nextgopg -c "\dt"
```
