---
name: database-migration
description: Manage database migrations and Better Auth schema. Use when adding tables, modifying schema, running migrations, or resetting the database.
allowed-tools: Read, Bash, Grep
---

# Database Migration

Manage PostgreSQL database, Better Auth schema, and GORM AutoMigrate.

## GORM AutoMigrate (Backend)

Jeder Bounded Context unter `backend/internal/<ctx>/` hat seine eigene `infrastructure/persistence/registry.go` mit einer `Entities() []any` Funktion. Der Composition Root sammelt die Listen ein:

```go
// backend/internal/<ctx>/infrastructure/persistence/registry.go
func Entities() []any {
    return []any{&gormNewEntity{}}  // GORM-tagged twin der pure-domain Entity
}

// backend/internal/composition/composition.go (runAutoMigrations)
entities := []any{}
entities = append(entities, statspersist.Entities()...)
entities = append(entities, <newctx>persist.Entities()...)  // ← neu
```

Es gibt **keine zentrale Registry** mehr. `cmd/server/main.go` ruft nur `composition.Build` auf — die Migrationen laufen automatisch beim Backend-Start.

## Database Commands

### Start Database

```bash
just db-up
```

### Stop Database

```bash
just db-down
```

### Reset Database (delete all data)

```bash
just db-reset
```

## Better Auth Tables

Better Auth uses these tables (auto-created via migration):

- `user` - User accounts
- `session` - Active sessions
- `account` - OAuth/credential accounts
- `verification` - Email verification tokens

### Run Better Auth Migration

```bash
just db-migrate
```

Or manually:

```bash
cd frontend && bunx dotenv-cli -e .env.local -- bunx @better-auth/cli@latest migrate --config src/shared/lib/auth-server/auth.ts --yes
```

### Generate Migration SQL (without applying)

```bash
cd frontend && bunx dotenv-cli -e .env.local -- bunx @better-auth/cli@latest generate --config src/shared/lib/auth-server/auth.ts
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
just db-up
```

### Reset everything

```bash
just db-reset
just db-migrate
```
