#!/bin/bash
# Custom hatchet-lite entrypoint. Differences vs. the stock entrypoint:
# - Lets `quickstart` generate self-signed certs (stock uses `--skip certs`,
#   which then leaves the server unable to load TLS config).
# Hatchet v0.84+ is Postgres-only — no RabbitMQ provisioning needed.

set -e

# Schema migrations against external Postgres.
/hatchet-migrate
echo "Migrations applied."

# Generate config + certs + keys + seed initial tenant.
/hatchet-admin quickstart --generated-config-dir ./config --overwrite=false

exec /hatchet-lite --config ./config
