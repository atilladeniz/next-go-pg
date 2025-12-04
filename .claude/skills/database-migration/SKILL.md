---
name: database-migration
description: Manage database migrations and Better Auth schema. Use when adding tables, modifying schema, running migrations, or resetting the database.
allowed-tools: Read, Bash, Grep
---

# Database Migration

Manage PostgreSQL database, Better Auth schema, and GORM AutoMigrate.

## GORM AutoMigrate (Backend)

Nach `goca feature` muss die neue Entity in `backend/cmd/server/main.go` registriert werden:

```go
// Nach DB-Verbindung
db.AutoMigrate(
    &domain.User{},
    &domain.UserStats{},
    &domain.NewEntity{},  // Neue Entity hinzufuegen
)
```

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
cd frontend && DATABASE_URL="postgres://postgres:postgres@localhost:5432/gocatest" bunx @better-auth/cli migrate -y
```

### Generate Migration SQL (without applying)

```bash
cd frontend && DATABASE_URL="postgres://postgres:postgres@localhost:5432/gocatest" bunx @better-auth/cli generate
```

## Direct Database Access

### Connect to PostgreSQL

```bash
docker exec gocatest-db-1 psql -U postgres -d gocatest
```

### List Tables

```bash
docker exec gocatest-db-1 psql -U postgres -d gocatest -c "\dt"
```

### Describe Table

```bash
docker exec gocatest-db-1 psql -U postgres -d gocatest -c "\d user"
```

### Run SQL Query

```bash
docker exec gocatest-db-1 psql -U postgres -d gocatest -c "SELECT * FROM \"user\""
```

## Database Connection

### Connection String

```
postgres://postgres:postgres@localhost:5432/gocatest
```

### Environment Variable

Set in `frontend/.env.local`:

```
DATABASE_URL=postgres://postgres:postgres@localhost:5432/gocatest
```

## Troubleshooting

### Table doesn't exist

Run Better Auth migration:

```bash
cd frontend && DATABASE_URL="postgres://postgres:postgres@localhost:5432/gocatest" bunx @better-auth/cli migrate -y
```

### Connection refused

Start the database:

```bash
make db-up
```

### Reset everything

```bash
make db-reset
cd frontend && DATABASE_URL="postgres://postgres:postgres@localhost:5432/gocatest" bunx @better-auth/cli migrate -y
```
