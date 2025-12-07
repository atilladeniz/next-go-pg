#!/bin/sh
set -e

echo "Running database migrations..."
/app/migrate -up || echo "Migration skipped or already up to date"

echo "Starting server..."
exec /app/backend "$@"
