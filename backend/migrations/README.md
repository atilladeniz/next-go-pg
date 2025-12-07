# Database Migrations

SQL migrations managed with [golang-migrate](https://github.com/golang-migrate/migrate).

## Structure

```
migrations/
├── 001_initial.up.sql      # Create tables
├── 001_initial.down.sql    # Drop tables
├── 002_feature.up.sql      # Add feature
├── 002_feature.down.sql    # Remove feature
└── ...
```

## Commands

```bash
# Run all pending migrations
make migrate-up

# Rollback last migration
make migrate-down

# Show current version
make migrate-version

# Create new migration files
make migrate-create name=add_users_table
```

## Migration vs AutoMigrate

This project uses **two migration strategies**:

1. **SQL Migrations (this folder)**: For schema changes that need precise control
   - Table creation with specific constraints
   - Index optimization
   - Data migrations
   - Production deployments

2. **GORM AutoMigrate**: For development convenience
   - Runs automatically on server start
   - Only adds columns/tables, never removes
   - Good for rapid iteration

**Recommended workflow:**
- Development: AutoMigrate handles schema sync
- Production: Run `make migrate-up` before deployment

## Creating Migrations

```bash
# 1. Create migration files
make migrate-create name=add_products_table

# 2. Edit the .up.sql file
# backend/migrations/002_add_products_table.up.sql

# 3. Edit the .down.sql file (reverse the changes)
# backend/migrations/002_add_products_table.down.sql

# 4. Test the migration
make migrate-up
make migrate-down
make migrate-up
```

## Best Practices

1. **Always write down migrations** - Every up migration needs a matching down
2. **Test rollbacks** - Run `migrate-up` then `migrate-down` then `migrate-up`
3. **Keep migrations small** - One logical change per migration
4. **Never edit applied migrations** - Create a new migration instead
5. **Use transactions** - Wrap DDL in BEGIN/COMMIT when possible
