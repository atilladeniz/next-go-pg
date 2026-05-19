# Production Database Migrations

This folder holds versioned SQL migrations driven by [golang-migrate](https://github.com/golang-migrate/migrate).

**The dev path does not run these.** `just dev` boots `cmd/server`, which runs GORM `AutoMigrate` (entity registry at `internal/domain/registry.go`) and River migrations on startup. The folder is currently empty — add a numbered `*.up.sql` / `*.down.sql` pair only when a production deploy needs precise schema control that AutoMigrate can't provide.

## When to use SQL migrations instead of AutoMigrate

| Scenario | Recommended path |
|----------|------------------|
| Adding a new entity in dev | Add Go struct → register in `registry.go` → AutoMigrate picks it up |
| Production with one replica (current Kamal setup) | AutoMigrate on boot is still fine |
| Production with multiple replicas | Use these SQL migrations as a deploy hook — AutoMigrate racing across replicas is unsafe |
| Schema change AutoMigrate can't do (drop column, rename, data backfill, complex index) | Add an SQL migration here |
| Reversible rollouts | SQL migrations have `.down.sql` files; AutoMigrate has no rollback path |

## Layout

```
backend/migrations/
├── 001_<name>.up.sql      # Apply
├── 001_<name>.down.sql    # Rollback
├── 002_<name>.up.sql
├── 002_<name>.down.sql
└── README.md
```

## Commands (run from repo root)

```bash
just prod-migrate-create <name>   # Scaffolds <NNN>_<name>.up.sql + .down.sql
just prod-migrate-up              # Apply all pending
just prod-migrate-down            # Rollback last
just prod-migrate-version         # Show current version
```

These wrap `backend/cmd/migrate`, which uses golang-migrate against this folder. The companion River migrations live in `cmd/river-migrate` (`just prod-river-migrate-up`).

## Best practices

1. **Write both halves.** Every `.up.sql` needs a matching `.down.sql`.
2. **Round-trip locally.** `prod-migrate-up` → `prod-migrate-down` → `prod-migrate-up` before opening a PR.
3. **Small, focused migrations.** One logical change per file.
4. **Never edit a migration that already shipped.** Add a follow-up instead.
5. **Avoid AutoMigrate drift.** If you add an SQL migration, also update the matching Go struct so AutoMigrate stays a no-op (or remove AutoMigrate from `cmd/server` entirely when you switch the project to SQL-only).
