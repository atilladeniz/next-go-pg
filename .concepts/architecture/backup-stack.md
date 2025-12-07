# Backup Stack Architecture

Fully automatic PostgreSQL backup system with RustFS (S3-compatible storage).

## Stack Overview

```mermaid
flowchart TB
    subgraph BackupStack["Backup Stack (make backup-up)"]
        subgraph AutoBackup["Automatic Backup"]
            PGB[postgres-backup-s3]
            INIT[rustfs-init]
        end

        subgraph Storage
            RFS[RustFS :9000/:9001]
            RFSV[(rustfs_data)]
        end
    end

    subgraph DevStack["Dev Stack (make dev)"]
        DB[(db :5432)]
    end

    INIT -->|"create bucket"| RFS
    PGB -->|"pg_dump"| DB
    PGB -->|"S3 upload"| RFS
    RFS --> RFSV

    style BackupStack fill:#e1f5fe
    style DevStack fill:#f3e5f5
```

## Component Details

### postgres-backup-s3

Automatic PostgreSQL backup container with S3 upload.

```mermaid
flowchart LR
    subgraph Container["postgres-backup-s3"]
        Cron[Cron Scheduler]
        Dump[pg_dump]
        Upload[aws s3 cp]
        Cleanup[Retention Cleanup]
    end

    subgraph Target
        DB[(App Database)]
    end

    subgraph Destination
        S3[RustFS S3]
    end

    Cron -->|"@daily"| Dump
    Dump --> DB
    Dump --> Upload
    Upload --> S3
    Upload --> Cleanup
```

**Features:**
- Fully automatic (cron-based)
- Zero configuration via environment variables
- S3-compatible storage support
- Configurable retention policy
- Compression (custom format)

### RustFS

S3-compatible object storage (MinIO alternative, written in Rust).

```mermaid
flowchart TB
    subgraph RustFS["RustFS Container"]
        API[S3 API :9000]
        Console[Web Console :9001]
        Engine[Storage Engine]
    end

    subgraph Volume
        Data[(rustfs_data)]
    end

    subgraph Clients
        Backup[postgres-backup-s3]
        CLI[mc/aws-cli]
    end

    Backup -->|"S3 Protocol"| API
    CLI -->|"S3 Protocol"| API
    API --> Engine
    Console --> Engine
    Engine --> Data
```

**Features:**
- S3-compatible API
- Web console for management
- Bucket management
- Lightweight (~20MB)
- Active development (Rust-based)

## Data Flow

### Backup Flow

```mermaid
sequenceDiagram
    participant C as Cron
    participant P as postgres-backup-s3
    participant D as App Database
    participant R as RustFS

    C->>P: Trigger backup (@daily)
    P->>D: pg_dump --format=custom
    D-->>P: Backup data (.dump)
    P->>R: aws s3 cp to backups/postgres/
    R-->>P: Upload complete
    P->>P: Delete backups older than 7 days
```

### Restore Flow

```mermaid
sequenceDiagram
    participant A as Admin
    participant P as postgres-backup-s3
    participant R as RustFS
    participant D as App Database

    A->>P: make backup-restore
    P->>R: Download latest backup
    R-->>P: Backup data
    P->>D: pg_restore
    D-->>P: Restore complete
    P-->>A: Success
```

## Network Architecture

```mermaid
flowchart TB
    subgraph DockerNetwork["Docker Network: next-go-pg_default"]
        DB[(db)]
        PGB[postgres-backup-s3]
        RFS[rustfs]
        INIT[rustfs-init]
    end

    subgraph Host["Host Machine"]
        Browser[Browser]
    end

    Browser -->|":9001"| RFS
    INIT -->|"create bucket"| RFS
    PGB -->|"internal"| DB
    PGB -->|"internal"| RFS
```

## Docker Compose Services

| Service | Image | Ports | Purpose |
|---------|-------|-------|---------|
| `rustfs` | rustfs/rustfs | 9000, 9001 | S3-compatible storage |
| `rustfs-init` | minio/mc | - | Auto-create bucket on startup |
| `pg-backup` | eeshugerman/postgres-backup-s3:16 | - | Automatic backups |

## Volume Structure

```mermaid
flowchart TB
    subgraph Volumes
        V1[rustfs_data]
    end

    subgraph Content
        C1[S3 bucket: backups/postgres/]
    end

    V1 --- C1
```

## Setup Workflow

```mermaid
flowchart TB
    A[make backup-up] --> B[rustfs-init creates bucket]
    B --> C[pg-backup waits for DB healthy]
    C --> D[Automatic daily backups]
    D --> E[Retention: 7 days]

    style A fill:#4caf50
    style D fill:#2196f3
```

**No manual configuration required!**

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `S3_ACCESS_KEY` | rustfsadmin | RustFS access key |
| `S3_SECRET_KEY` | rustfsadmin | RustFS secret key |
| `BACKUP_SCHEDULE` | @daily | Cron schedule for backups |
| `BACKUP_KEEP_DAYS` | 7 | Days to keep backups |
| `POSTGRES_DB` | nextgopg | Database name |
| `POSTGRES_USER` | postgres | Database user |
| `POSTGRES_PASSWORD` | postgres | Database password |

## Makefile Commands

```bash
make backup-up       # Start automatic backup system
make backup-down     # Stop backup stack
make backup-now      # Create backup immediately
make backup-list     # List all backups in S3
make backup-restore  # Restore from latest backup
```

## Console Access

RustFS Web Console: http://localhost:9001/rustfs/console/

Default credentials:
- Access Key: `rustfsadmin`
- Secret Key: `rustfsadmin`
