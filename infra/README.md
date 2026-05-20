# Infrastructure

Everything that lives outside the application code: Docker build, Kamal
deployment, observability (Loki + Grafana), and the docker-compose
files used for dev and CI.

```
infra/
├── docker/             # Dockerfile + supervisord.conf (production image)
├── compose/            # docker-compose.{yml,dev,logging,backup}.yml
│   └── postgres-init/  # SQL run once on first Postgres init
├── kamal/              # deploy.yml + deploy.{staging,production}.yml
│                       # + hooks/ + secrets.example
├── loki/               # Loki + Promtail configs (dev + prod)
└── grafana/            # Grafana provisioning (datasources + dashboards)
```

This document covers Kamal deployment and the optional AI dev stack
(Hatchet workflow engine + Ollama LLM runtime).

## AI dev stack (optional)

The AI dev stack adds a self-hosted `hatchet-lite` workflow engine and a
local `ollama` LLM runtime to the dev Docker Compose file. Both run
behind the `ai` Compose profile, so they are NOT started by `just dev` —
opt in explicitly with `just ai-up`.

### Quick start

```bash
# 1. Bring up core dev stack (db, mailpit) and AI services
just ai-up

# 2. First-run only: pull the default model (~5 GB for gemma4:e4b)
just ai-pull-model              # uses default gemma4:e4b
# or pick a different one:
just ai-pull-model gemma4:26b-a4b

# 3. Hatchet dashboard at http://localhost:8888  (anonymous in dev mode)
```

### Endpoints

| Service        | URL                                  | Notes                                  |
|----------------|--------------------------------------|----------------------------------------|
| Hatchet gRPC   | `hatchet-lite:7077` (docker net)     | Backend worker dials this              |
| Hatchet UI     | `http://localhost:8888`              | Browser dashboard                      |
| Ollama HTTP    | `http://ollama:11434` (docker net)   | Backend HTTP client target             |

### What gets created

- A dedicated `hatchet` Postgres database (script:
  `infra/compose/postgres-init/01-create-hatchet-db.sql`). Created
  automatically on a **fresh** Postgres data volume. If your `postgres_data`
  volume already exists, create it manually:
  ```bash
  docker compose -f infra/compose/docker-compose.dev.yml \
    exec db psql -U postgres -c 'CREATE DATABASE hatchet'
  ```
- An `ollama_data` named volume that persists pulled models between
  restarts.

### Expected resource usage

- `hatchet-lite`: ~150 MB RAM idle.
- `ollama` + `gemma4:e4b`: ~5 GB disk for the model, 6–8 GB RAM during
  inference on a dev MacBook (M1/M2). Switch to `gemma4:26b-a4b` only on
  a workstation with ≥32 GB RAM.

### Known gotchas

- **Postgres init script only fires on a fresh volume.** Existing dev
  volumes need the manual `CREATE DATABASE hatchet` step above.
- **Hatchet broadcast address is fragile.** `SERVER_GRPC_BROADCAST_ADDRESS`
  must be reachable from any worker. Inside the docker network we use the
  service name `hatchet-lite:7077`. If you point a worker from outside the
  docker network at it, the worker will fail to resume long-running tasks.
- **First model pull is slow** (~5 GB over the network). The Ollama
  healthcheck does not pre-pull — that's an explicit `just ai-pull-model`
  step so contributors choose when to spend the disk and bandwidth.
- **Hatchet Go SDK has no OpenTelemetry support as of May 2026.** Engine
  and worker logs flow via stdout → Promtail → Loki, queryable in Grafana
  with `{service="hatchet"}`.

### Stopping

```bash
just ai-down            # stops hatchet-lite + ollama, keeps db running
```

### Production deploy

The AI dev stack is **not** wired into `infra/kamal/deploy.yml`. Promoting
`hatchet-lite` (or a managed Hatchet Cloud workspace) to production is
Phase 2 of the orchestrator plan — see
[`.docs/orchestrator-decision.md`](../.docs/orchestrator-decision.md).

## Quick Reference

```bash
# Deploy to staging
make deploy-staging

# Deploy to production (with confirmation)
make deploy-production

# Rollback
make deploy-rollback

# Show logs
make deploy-logs

# Console on server
make deploy-console
```

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        Your Server                          │
│  ┌───────────────────────────────────────────────────────┐  │
│  │                    kamal-proxy                        │  │
│  │              (SSL termination, routing)               │  │
│  │                    :80 / :443                         │  │
│  └───────────────────┬───────────────────────────────────┘  │
│                      │                                      │
│  ┌───────────────────▼───────────────────────────────────┐  │
│  │              Docker Container (next-go-pg)            │  │
│  │  ┌─────────────────┐    ┌─────────────────────────┐   │  │
│  │  │   Backend :8080 │────▶   Loki :3100            │   │  │
│  │  │   (Go + Mux)    │    │   (Log Aggregation)     │   │  │
│  │  └─────────────────┘    └─────────────────────────┘   │  │
│  │  ┌─────────────────┐    ┌─────────────────────────┐   │  │
│  │  │   Frontend :3000│────▶   Grafana :3001         │   │  │
│  │  │   (Next.js)     │    │   (Log Visualization)   │   │  │
│  │  └─────────────────┘    └─────────────────────────┘   │  │
│  │              Managed by supervisord                   │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
│  ┌───────────────────────────────────────────────────────┐  │
│  │            PostgreSQL (Managed DB recommended)        │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

## Logging (Self-Hosted + Tailscale)

Loki + Grafana run as Kamal accessories, **accessible only via Tailscale** (not public internet).

- **Loki** (:3100) - Log aggregation, stores logs locally
- **Grafana** (:3001) - Dashboard for viewing logs

Both Backend and Frontend send logs directly to Loki via HTTP.

### Tailscale Setup (Required)

1. **Install Tailscale on VPS**:
   ```bash
   ssh root@your-server-ip
   curl -fsSL https://tailscale.com/install.sh | sh
   tailscale up
   ```

2. **Get Tailscale credentials**:
   ```bash
   # Get Tailscale IP
   tailscale ip -4
   # Output: 100.x.x.x

   # Get MagicDNS hostname
   tailscale status --self --json | jq -r '.Self.DNSName' | sed 's/\.$//'
   # Output: your-vps.your-tailnet.ts.net
   ```

3. **Add to secrets**:
   ```bash
   # infra/kamal/secrets
   TAILSCALE_IP=100.x.x.x
   TAILSCALE_HOSTNAME=your-vps.your-tailnet.ts.net
   ```

### Access Grafana

After deployment (from any device in your Tailnet):
```
https://your-vps.your-tailnet.ts.net:3001
User: admin
Password: (from GF_SECURITY_ADMIN_PASSWORD secret)
```

**Note**: Grafana is NOT accessible from the public internet - only from devices connected to your Tailscale network.

## Setup (Initial Deployment)

### 1. Prerequisites

```bash
# Install Kamal
gem install kamal

# Or via Docker
docker run -it ghcr.io/basecamp/kamal:latest
```

### 2. Prepare Server

- Ubuntu 22.04+ (or Debian)
- SSH key access as root
- Ports 80, 443 open

```bash
# Enter server IP in deploy.staging.yml or deploy.production.yml
vim infra/kamal/deploy.staging.yml
# (or infra/kamal/deploy.production.yml)
```

### 3. Configure Secrets

```bash
# Create secrets file
cp infra/kamal/secrets.example infra/kamal/secrets

# Fill in secrets
vim infra/kamal/secrets
```

### 4. Create GitHub Container Registry Token

1. https://github.com/settings/tokens
2. "Generate new token (classic)"
3. Scope: `write:packages`, `read:packages`
4. Enter token in `infra/kamal/secrets`

### 5. First Deployment

```bash
# Bootstrap server + deploy
make deploy-setup
# → Choose "staging" or "production"
```

## Kamal Directory Structure

```
infra/kamal/
├── deploy.yml              # Base configuration
├── deploy.staging.yml      # Staging-specific overrides
├── deploy.production.yml   # Production-specific overrides
├── secrets                 # Secrets (gitignored — copy from secrets.example)
├── secrets.example         # Secrets template
└── hooks/
    ├── pre-build           # Before Docker build
    ├── pre-deploy          # Before deployment
    └── post-deploy         # After deployment

infra/docker/
├── Dockerfile              # Multi-stage build
└── supervisord.conf        # Process manager
```

## Environments

### Staging
- Server: `staging.example.com`
- Purpose: Testing before production
- Automatic deploys: On merge to `develop`

### Production
- Server: `app.example.com`
- Purpose: Live system
- Deploys: Manual with confirmation

## Scaling

### Horizontal (more servers)

```yaml
# infra/kamal/deploy.production.yml
servers:
  web:
    hosts:
      - web1.example.com
      - web2.example.com    # New server
      - web3.example.com    # Another server
```

Then add load balancer in front (Hetzner LB, Cloudflare, etc.)

### Vertical (larger server)

```yaml
# infra/kamal/deploy.production.yml
servers:
  web:
    hosts:
      - web1.example.com
    options:
      memory: 2g    # More RAM
      cpus: 4       # More CPUs
```

## Troubleshooting

### Deployment fails

```bash
# Show logs
kamal app logs -c infra/kamal/deploy.yml -d staging

# Container status
kamal details -c infra/kamal/deploy.yml -d staging
```

### Rollback needed

```bash
make deploy-rollback
# Choose version from list
```

### Manually restart container

```bash
kamal app boot -c infra/kamal/deploy.yml -d staging
```

### SSH to server

```bash
make deploy-console
```

## Best Practices

1. **Staging first**: Always test on staging first
2. **Database backups**: Before every production deploy
3. **Monitoring**: Set up uptime check (e.g., Uptime Kuma)
4. **Secrets**: Never commit to Git
5. **Rolling deploys**: With multiple servers (is default)

## Cost Overview

| Setup | Server | DB | Cost/Month |
|-------|--------|------|------------|
| Minimal | Hetzner CX22 (2 vCPU, 4GB) | SQLite/local | ~€5 |
| Standard | Hetzner CX32 (4 vCPU, 8GB) | Managed PostgreSQL | ~€30 |
| HA | 2x CX22 + LB | Managed PostgreSQL | ~€50 |

## Links

- [Kamal Docs](https://kamal-deploy.org/)
- [GitHub Container Registry](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry)
- [Hetzner Cloud](https://www.hetzner.com/cloud)
