# Deployment Architecture

## Deployment-Strategien

### Single Server (Phase 1)

Einfachste Konfiguration für Start und kleine Last.

```mermaid
flowchart TB
    subgraph Internet
        U[User]
    end

    subgraph Server["Hetzner CX22 (~€5/mo)"]
        subgraph Docker
            P[kamal-proxy :80/:443]
            C[gocatest Container]
            subgraph Container
                F[Frontend :3000]
                B[Backend :8080]
            end
        end
        D[(PostgreSQL :5432)]
    end

    U -->|HTTPS| P
    P --> C
    F <--> B
    B --> D
```

**Vorteile:**
- Günstig (~€5-15/Monat)
- Einfach zu verwalten
- Zero-Downtime Deploys

**Nachteile:**
- Single Point of Failure
- Keine horizontale Skalierung

### Multi-Server mit Load Balancer (Phase 2)

Für höhere Verfügbarkeit und Last.

```mermaid
flowchart TB
    subgraph Internet
        U[User]
        CF[Cloudflare CDN]
    end

    subgraph Infrastructure
        LB[Load Balancer :443]

        subgraph Web1["Web Server 1"]
            P1[kamal-proxy]
            C1[gocatest]
        end

        subgraph Web2["Web Server 2"]
            P2[kamal-proxy]
            C2[gocatest]
        end

        subgraph Database["Database Cluster"]
            DB1[(Primary)]
            DB2[(Replica)]
        end
    end

    U --> CF
    CF --> LB
    LB --> P1
    LB --> P2
    P1 --> C1
    P2 --> C2
    C1 --> DB1
    C2 --> DB1
    DB1 -.->|Replication| DB2
```

**Vorteile:**
- High Availability
- Horizontale Skalierung
- Failover möglich

**Kosten:**
- 2x Web Server: ~€20
- Load Balancer: ~€5
- Managed DB: ~€15
- **Total: ~€40-50/Monat**

## Kamal Deployment Flow

```mermaid
sequenceDiagram
    participant D as Developer
    participant K as Kamal CLI
    participant R as Registry
    participant S as Server
    participant P as kamal-proxy

    D->>K: make deploy-staging
    K->>K: Run pre-build hook
    K->>K: Build Docker Image
    K->>R: Push Image
    K->>S: Pull Image
    K->>S: Start new Container
    K->>P: Health Check
    P-->>K: 200 OK
    K->>P: Route Traffic to New
    K->>S: Stop Old Container
    K-->>D: Deploy Complete
```

## Zero-Downtime Deploy

```mermaid
stateDiagram-v2
    [*] --> Running: Container v1 läuft

    Running --> Deploying: kamal deploy
    Deploying --> HealthCheck: Neuer Container startet

    HealthCheck --> Switching: Health OK
    HealthCheck --> Rollback: Health Failed

    Switching --> Running: Traffic umgeleitet
    Rollback --> Running: Alter Container weiter

    note right of Switching
        Kein Downtime!
        Proxy wechselt
        nahtlos
    end note
```

## Server Setup

```mermaid
flowchart LR
    subgraph Preparation
        A[Ubuntu 22.04] --> B[SSH Key]
        B --> C[Firewall: 80, 443, 22]
    end

    subgraph Kamal
        D[kamal setup] --> E[Docker Install]
        E --> F[kamal-proxy Start]
        F --> G[App Deploy]
    end

    Preparation --> Kamal
```

## Monitoring & Logging

```mermaid
flowchart TB
    subgraph Server
        App[gocatest]
        Logs[/var/log/]
    end

    subgraph Monitoring
        UK[Uptime Kuma]
        BT[Better Stack]
    end

    subgraph Alerts
        E[Email]
        S[Slack]
    end

    App -->|stdout| Logs
    App -->|/health| UK
    App -->|/health| BT
    UK -->|Alert| E
    UK -->|Alert| S
    BT -->|Alert| E
```

## Backup Strategy

```mermaid
flowchart LR
    subgraph Daily
        D1[DB Dump]
        D2[Upload to S3]
    end

    subgraph Weekly
        W1[Full Server Snapshot]
    end

    subgraph Recovery
        R1[Restore from S3]
        R2[Restore Snapshot]
    end

    D1 --> D2
    D2 -.->|Recovery| R1
    W1 -.->|Recovery| R2
```
