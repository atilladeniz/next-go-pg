# Deployment with Kamal

Zero-downtime Docker deployments to any server.

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

## Logging (Self-Hosted)

Loki + Grafana run as Kamal accessories on the same server:

- **Loki** (:3100) - Log aggregation, stores logs locally
- **Grafana** (:3001) - Dashboard for viewing logs

Both Backend and Frontend send logs directly to Loki via HTTP.

### Access Grafana

After deployment:
```
https://your-server:3001
User: admin
Password: (from GF_SECURITY_ADMIN_PASSWORD secret)
```

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
vim deploy/config/deploy.staging.yml
```

### 3. Configure Secrets

```bash
# Create secrets file
cp deploy/.kamal/secrets.example deploy/.kamal/secrets

# Fill in secrets
vim deploy/.kamal/secrets
```

### 4. Create GitHub Container Registry Token

1. https://github.com/settings/tokens
2. "Generate new token (classic)"
3. Scope: `write:packages`, `read:packages`
4. Enter token in `deploy/.kamal/secrets`

### 5. First Deployment

```bash
# Bootstrap server + deploy
make deploy-setup
# → Choose "staging" or "production"
```

## Directory Structure

```
deploy/
├── config/
│   ├── deploy.yml              # Base configuration
│   ├── deploy.staging.yml      # Staging-specific
│   └── deploy.production.yml   # Production-specific
├── .kamal/
│   ├── secrets                 # Secrets (DO NOT commit!)
│   ├── secrets.example         # Secrets template
│   └── hooks/
│       ├── pre-build           # Before Docker build
│       ├── pre-deploy          # Before deployment
│       └── post-deploy         # After deployment
├── Dockerfile                  # Multi-stage build
├── supervisord.conf            # Process manager
└── README.md                   # This file
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
# deploy/config/deploy.production.yml
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
# deploy/config/deploy.production.yml
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
kamal app logs -c deploy/config/deploy.yml -d staging

# Container status
kamal details -c deploy/config/deploy.yml -d staging
```

### Rollback needed

```bash
make deploy-rollback
# Choose version from list
```

### Manually restart container

```bash
kamal app boot -c deploy/config/deploy.yml -d staging
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
