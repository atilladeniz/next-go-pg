---
name: database-migration
description: Manage database migrations and Better Auth schema. Use when adding tables, modifying schema, running migrations, or resetting the database.
allowed-tools: Read, Bash, Grep
---

# Database Migration

Manage PostgreSQL database, Better Auth schema, and GORM AutoMigrate.

## GORM AutoMigrate (Backend)

Nach `goca feature` muss die neue Entity in `backend/internal/domain/registry.go` registriert werden:

```go
// internal/domain/registry.go
func AllEntities() []interface{} {
    return []interface{}{
        &UserStats{},
        &NewEntity{},  // Neue Entity hier hinzufuegen
    }
}
```

Das ist die **EINZIGE** Stelle - `main.go` bleibt unveraendert!
GORM erstellt die Tabelle automatisch beim Backend-Start.

## Database Commands

### Start Database

```bash
make db-up
```

### Stop Database

```bash
make db-down
```

### Reset Database (delete all data)

```bash
make db-reset
```

## Better Auth Tables

Better Auth uses these tables (auto-created via migration):

- `user` - User accounts
- `session` - Active sessions
- `account` - OAuth/credential accounts
- `verification` - Email verification tokens

### Run Better Auth Migration

```bash
cd frontend && DATABASE_URL="postgres://postgres:postgres@localhost:5432/nextgopg" bunx @better-auth/cli migrate -y
```

### Generate Migration SQL (without applying)

```bash
cd frontend && DATABASE_URL="postgres://postgres:postgres@localhost:5432/nextgopg" bunx @better-auth/cli generate
```

## Direct Database Access

### Connect to PostgreSQL

```bash
docker exec nextgopg-db-1 psql -U postgres -d nextgopg
```

### List Tables

```bash
docker exec nextgopg-db-1 psql -U postgres -d nextgopg -c "\dt"
```

### Describe Table

```bash
docker exec nextgopg-db-1 psql -U postgres -d nextgopg -c "\d user"
```

### Run SQL Query

```bash
docker exec nextgopg-db-1 psql -U postgres -d nextgopg -c "SELECT * FROM \"user\""
```

## Database Connection

### Connection String

```
postgres://postgres:postgres@localhost:5432/nextgopg
```

### Environment Variable

Set in `frontend/.env.local`:

```
DATABASE_URL=postgres://postgres:postgres@localhost:5432/nextgopg
```

## Troubleshooting

### Table doesn't exist

Run Better Auth migration:

```bash
cd frontend && DATABASE_URL="postgres://postgres:postgres@localhost:5432/nextgopg" bunx @better-auth/cli migrate -y
```

### Connection refused

Start the database:

```bash
make db-up
```

### Reset everything

```bash
make db-reset
cd frontend && DATABASE_URL="postgres://postgres:postgres@localhost:5432/nextgopg" bunx @better-auth/cli migrate -y
```
