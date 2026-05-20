-- Bootstrap a dedicated database for Hatchet (workflow engine).
-- Runs only on FIRST Postgres init (empty data volume). For an already
-- populated volume, create it manually:
--   docker compose -f infra/compose/docker-compose.dev.yml \
--     exec db psql -U postgres -c 'CREATE DATABASE hatchet'
SELECT 'CREATE DATABASE hatchet'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'hatchet')\gexec
