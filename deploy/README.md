# Deployment mit Kamal

Zero-downtime Docker deployments auf beliebige Server.

## Quick Reference

```bash
# Staging deployen
make deploy-staging

# Production deployen (mit Bestätigung)
make deploy-production

# Rollback
make deploy-rollback

# Logs anzeigen
make deploy-logs

# Console auf Server
make deploy-console
```

## Architektur

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
│  │  │   Backend :8080 │    │   Frontend :3000        │   │  │
│  │  │   (Go + Mux)    │    │   (Next.js standalone)  │   │  │
│  │  └─────────────────┘    └─────────────────────────┘   │  │
│  │              Managed by supervisord                   │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
│  ┌───────────────────────────────────────────────────────┐  │
│  │            PostgreSQL (Managed DB empfohlen)          │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

## Setup (Erstmaliges Deployment)

### 1. Voraussetzungen

```bash
# Kamal installieren
gem install kamal

# Oder via Docker
docker run -it ghcr.io/basecamp/kamal:latest
```

### 2. Server vorbereiten

- Ubuntu 22.04+ (oder Debian)
- SSH-Key Zugang als root
- Ports 80, 443 offen

```bash
# Server IP in deploy.staging.yml oder deploy.production.yml eintragen
vim deploy/config/deploy.staging.yml
```

### 3. Secrets konfigurieren

```bash
# Secrets-Datei erstellen
cp deploy/.kamal/secrets.example deploy/.kamal/secrets

# Secrets ausfüllen
vim deploy/.kamal/secrets
```

### 4. GitHub Container Registry Token erstellen

1. https://github.com/settings/tokens
2. "Generate new token (classic)"
3. Scope: `write:packages`, `read:packages`
4. Token in `deploy/.kamal/secrets` eintragen

### 5. Erstes Deployment

```bash
# Server bootstrappen + deployen
make deploy-setup
# → Wähle "staging" oder "production"
```

## Verzeichnisstruktur

```
deploy/
├── config/
│   ├── deploy.yml              # Basis-Konfiguration
│   ├── deploy.staging.yml      # Staging-spezifisch
│   └── deploy.production.yml   # Production-spezifisch
├── .kamal/
│   ├── secrets                 # Secrets (NICHT committen!)
│   ├── secrets.example         # Secrets-Vorlage
│   └── hooks/
│       ├── pre-build           # Vor Docker Build
│       ├── pre-deploy          # Vor Deployment
│       └── post-deploy         # Nach Deployment
├── Dockerfile                  # Multi-stage Build
├── supervisord.conf            # Process Manager
└── README.md                   # Diese Datei
```

## Environments

### Staging
- Server: `staging.example.com`
- Zweck: Testen vor Production
- Automatische Deploys: Bei Merge in `develop`

### Production
- Server: `app.example.com`
- Zweck: Live-System
- Deploys: Manuell mit Bestätigung

## Skalierung

### Horizontal (mehr Server)

```yaml
# deploy/config/deploy.production.yml
servers:
  web:
    hosts:
      - web1.example.com
      - web2.example.com    # Neuer Server
      - web3.example.com    # Noch ein Server
```

Dann Load Balancer davor (Hetzner LB, Cloudflare, etc.)

### Vertikal (größerer Server)

```yaml
# deploy/config/deploy.production.yml
servers:
  web:
    hosts:
      - web1.example.com
    options:
      memory: 2g    # Mehr RAM
      cpus: 4       # Mehr CPUs
```

## Troubleshooting

### Deployment schlägt fehl

```bash
# Logs anzeigen
kamal app logs -c deploy/config/deploy.yml -d staging

# Container Status
kamal details -c deploy/config/deploy.yml -d staging
```

### Rollback nötig

```bash
make deploy-rollback
# Version auswählen aus Liste
```

### Container manuell neustarten

```bash
kamal app boot -c deploy/config/deploy.yml -d staging
```

### SSH auf Server

```bash
make deploy-console
```

## Best Practices

1. **Staging first**: Immer erst auf Staging testen
2. **Database Backups**: Vor jedem Production Deploy
3. **Monitoring**: Uptime-Check einrichten (z.B. Uptime Kuma)
4. **Secrets**: Nie in Git committen
5. **Rolling Deploys**: Bei mehreren Servern (ist default)

## Kosten-Übersicht

| Setup | Server | DB | Kosten/Monat |
|-------|--------|------|--------------|
| Minimal | Hetzner CX22 (2 vCPU, 4GB) | SQLite/local | ~€5 |
| Standard | Hetzner CX32 (4 vCPU, 8GB) | Managed PostgreSQL | ~€30 |
| HA | 2x CX22 + LB | Managed PostgreSQL | ~€50 |

## Links

- [Kamal Docs](https://kamal-deploy.org/)
- [GitHub Container Registry](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry)
- [Hetzner Cloud](https://www.hetzner.com/cloud)
