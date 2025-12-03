# Database Migrations

This directory contains database migration files.

## Structure
- *.up.sql - Migration files (apply changes)
- *.down.sql - Rollback files (reverse changes)

## Usage

### Using golang-migrate tool
bash
# Install golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
migrate -path ./migrations -database "postgres://user:password@localhost/dbname?sslmode=disable" up

# Rollback last migration
migrate -path ./migrations -database "postgres://user:password@localhost/dbname?sslmode=disable" down 1


### Manual execution
bash
# Apply migration
psql -h localhost -U postgres -d your_db -f migrations/001_initial.up.sql

# Rollback migration  
psql -h localhost -U postgres -d your_db -f migrations/001_initial.down.sql


## Creating new migrations
1. Create new files: 002_description.up.sql and 002_description.down.sql
2. Add your changes in the .up.sql file
3. Add the reverse changes in the .down.sql file
