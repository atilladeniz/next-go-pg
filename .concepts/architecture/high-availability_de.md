# High Availability Architecture

Enterprise-Grade Ausfallsicherheit nach dem Vorbild von Notion, Spotify, Shopify.

## Executive Summary

Dieses Dokument beschreibt die Evolution von Single-Server zu einer hochverfügbaren, global verteilten Architektur.

---

## Phase 1: Current State (Single Server)

Aktueller Stand der Infrastruktur.

```mermaid
flowchart TB
    subgraph Internet
        U[Users]
    end

    subgraph SingleServer["Single Server (SPOF)"]
        P[kamal-proxy]
        F[Frontend]
        B[Backend]
        D[(PostgreSQL)]
    end

    U -->|HTTPS| P
    P --> F & B
    B --> D

    style SingleServer fill:#ffcccc,stroke:#cc0000
```

**Probleme:**

- Single Point of Failure (SPOF)
- Keine horizontale Skalierung
- Downtime bei Deployments möglich
- Keine geografische Redundanz

---

## Phase 2: Regional High Availability

Erste Stufe der Ausfallsicherheit - alles in einer Region, aber redundant.

```mermaid
flowchart TB
    subgraph Internet
        U[Users]
        CDN[Cloudflare CDN]
    end

    subgraph Region["Region: EU-Central (Frankfurt)"]
        subgraph LoadBalancer["Layer 7 Load Balancer"]
            LB[HAProxy / Traefik]
            LB2[HAProxy / Traefik]
        end

        subgraph AppTier["Application Tier (Stateless)"]
            subgraph Node1["Node 1"]
                P1[kamal-proxy]
                F1[Frontend]
                B1[Backend]
            end
            subgraph Node2["Node 2"]
                P2[kamal-proxy]
                F2[Frontend]
                B2[Backend]
            end
            subgraph Node3["Node 3"]
                P3[kamal-proxy]
                F3[Frontend]
                B3[Backend]
            end
        end

        subgraph DataTier["Data Tier"]
            subgraph DBCluster["PostgreSQL Cluster"]
                DB1[(Primary)]
                DB2[(Sync Replica)]
                DB3[(Async Replica)]
            end

            subgraph Cache["Redis Cluster"]
                R1[Redis Primary]
                R2[Redis Replica]
            end
        end

        subgraph JobTier["Background Jobs"]
            W1[River Worker 1]
            W2[River Worker 2]
        end
    end

    U --> CDN
    CDN --> LB & LB2
    LB --> P1 & P2 & P3
    LB2 --> P1 & P2 & P3

    B1 & B2 & B3 --> DB1
    B1 & B2 & B3 --> R1

    DB1 -.->|Sync| DB2
    DB1 -.->|Async| DB3
    R1 -.-> R2

    W1 & W2 --> DB1

    style LoadBalancer fill:#90EE90
    style AppTier fill:#87CEEB
    style DataTier fill:#DDA0DD
```

### Komponenten-Details

| Komponente | Replicas | Strategie | Failover Zeit |
|------------|----------|-----------|---------------|
| Load Balancer | 2 | Active-Passive | < 10s |
| App Nodes | 3+ | Active-Active | 0s (instant) |
| PostgreSQL | 3 | Primary + Replicas | < 30s |
| Redis | 2 | Primary + Replica | < 5s |
| River Workers | 2+ | Competing Consumers | 0s |

---

## Phase 3: Multi-Region Active-Active

Globale Verfügbarkeit wie Notion/Shopify.

```mermaid
flowchart TB
    subgraph Global["Global Layer"]
        DNS[Route53 / Cloudflare<br/>GeoDNS + Health Checks]
        GSLB[Global Load Balancer]
    end

    subgraph EU["Region: EU (Frankfurt)"]
        subgraph EU_Edge["Edge"]
            EU_CDN[CDN PoP]
            EU_LB[Load Balancer]
        end

        subgraph EU_App["App Cluster"]
            EU_N1[Node 1-3]
        end

        subgraph EU_Data["Data"]
            EU_DB[(PostgreSQL<br/>Primary)]
            EU_R[Redis]
        end
    end

    subgraph US["Region: US-East (Virginia)"]
        subgraph US_Edge["Edge"]
            US_CDN[CDN PoP]
            US_LB[Load Balancer]
        end

        subgraph US_App["App Cluster"]
            US_N1[Node 1-3]
        end

        subgraph US_Data["Data"]
            US_DB[(PostgreSQL<br/>Read Replica)]
            US_R[Redis]
        end
    end

    subgraph APAC["Region: APAC (Singapore)"]
        subgraph APAC_Edge["Edge"]
            APAC_CDN[CDN PoP]
            APAC_LB[Load Balancer]
        end

        subgraph APAC_App["App Cluster"]
            APAC_N1[Node 1-3]
        end

        subgraph APAC_Data["Data"]
            APAC_DB[(PostgreSQL<br/>Read Replica)]
            APAC_R[Redis]
        end
    end

    DNS --> EU_CDN & US_CDN & APAC_CDN
    GSLB --> EU_LB & US_LB & APAC_LB

    EU_LB --> EU_N1
    US_LB --> US_N1
    APAC_LB --> APAC_N1

    EU_N1 --> EU_DB & EU_R
    US_N1 --> US_DB & US_R
    APAC_N1 --> APAC_DB & APAC_R

    EU_DB -.->|"Cross-Region<br/>Replication"| US_DB
    EU_DB -.->|"Cross-Region<br/>Replication"| APAC_DB

    style Global fill:#FFD700
    style EU fill:#90EE90
    style US fill:#87CEEB
    style APAC fill:#DDA0DD
```

---

## Database Tier: PostgreSQL High Availability (Deep Dive)

Dies ist die kritischste Komponente - hier liegt die Single Source of Truth.

### Warum ist die Datenbank so kritisch?

```mermaid
flowchart LR
    subgraph Risiken["Ausfallrisiken"]
        HW[Hardware-Ausfall]
        SW[Software-Bug]
        NET[Netzwerk-Partition]
        DC[Datacenter-Ausfall]
        HUM[Menschlicher Fehler]
        COR[Daten-Korruption]
    end

    subgraph Auswirkungen["Auswirkungen"]
        DOWN[Kompletter Ausfall]
        LOSS[Datenverlust]
        INC[Inkonsistenz]
    end

    HW & SW & NET --> DOWN
    DC & HUM --> DOWN & LOSS
    COR --> INC & LOSS

    style Risiken fill:#ffcccc
    style Auswirkungen fill:#ff6666
```

### PostgreSQL HA Architektur mit Patroni

```mermaid
flowchart TB
    subgraph PGCluster["PostgreSQL HA Cluster"]
        subgraph Primary["Primary Node (Leader)"]
            PG1[(PostgreSQL<br/>Read + Write)]
            Patroni1[Patroni Agent]
            WAL1[WAL Sender]
        end

        subgraph SyncReplica["Synchronous Replica"]
            PG2[(PostgreSQL<br/>Read Only)]
            Patroni2[Patroni Agent]
            WAL2[WAL Receiver]
        end

        subgraph AsyncReplica["Async Replicas (2x)"]
            PG3[(PostgreSQL<br/>Read Only)]
            PG4[(PostgreSQL<br/>Read Only)]
        end

        subgraph Consensus["Distributed Consensus (etcd)"]
            ETCD1[etcd Node 1]
            ETCD2[etcd Node 2]
            ETCD3[etcd Node 3]
        end

        subgraph ConnectionPool["Connection Pooling"]
            PGB1[PgBouncer 1]
            PGB2[PgBouncer 2]
        end
    end

    subgraph App["Application Layer"]
        B1[Backend 1]
        B2[Backend 2]
        B3[Backend 3]
    end

    B1 & B2 & B3 -->|Connection| PGB1 & PGB2

    PGB1 & PGB2 -->|"Write (Primary)"| PG1
    PGB1 & PGB2 -->|"Read (Replicas)"| PG2 & PG3 & PG4

    Patroni1 <-->|Leader Election| ETCD1 & ETCD2 & ETCD3
    Patroni2 <-->|Health Status| ETCD1 & ETCD2 & ETCD3

    WAL1 -->|"Sync Stream<br/>(commit wait)"| WAL2
    PG1 -.->|"Async Stream<br/>(no wait)"| PG3 & PG4

    style Primary fill:#90EE90,stroke:#006400
    style SyncReplica fill:#87CEEB,stroke:#00008B
    style AsyncReplica fill:#DDA0DD,stroke:#800080
    style Consensus fill:#FFD700,stroke:#B8860B
```

### Replikations-Modi erklärt

```mermaid
flowchart LR
    subgraph SyncReplication["Synchrone Replikation"]
        direction TB
        S_Write[1. Write Request]
        S_Primary[2. Primary schreibt]
        S_WAL[3. WAL an Replica]
        S_Confirm[4. Replica bestätigt]
        S_Commit[5. Commit an Client]

        S_Write --> S_Primary --> S_WAL --> S_Confirm --> S_Commit
    end

    subgraph AsyncReplication["Asynchrone Replikation"]
        direction TB
        A_Write[1. Write Request]
        A_Primary[2. Primary schreibt]
        A_Commit[3. Commit an Client]
        A_WAL[4. WAL an Replica<br/>verzögert]

        A_Write --> A_Primary --> A_Commit
        A_Primary -.-> A_WAL
    end

    style SyncReplication fill:#90EE90
    style AsyncReplication fill:#FFB6C1
```

| Modus | Datenverlust | Latenz | Use Case |
|-------|--------------|--------|----------|
| **Synchron** | 0 (garantiert) | +2-5ms | Finanzdaten, kritische Writes |
| **Asynchron** | Möglich (< 1s) | Minimal | Read Replicas, Analytics |

### Failover-Szenario: Primary fällt aus

```mermaid
sequenceDiagram
    participant App as Application
    participant PGB as PgBouncer
    participant P1 as Primary
    participant P2 as Sync Replica
    participant ETCD as etcd Cluster
    participant Alert as Alerting

    Note over P1: ❌ Primary fällt aus!

    rect rgb(255, 200, 200)
        P1-xP2: Heartbeat fehlt
        P1-xETCD: Health Check failed
    end

    ETCD->>ETCD: Leader Election startet
    ETCD->>P2: Du bist neuer Leader!

    rect rgb(200, 255, 200)
        P2->>P2: Promote to Primary
        P2->>ETCD: Registriere als Primary
        ETCD->>PGB: Update Routing Config
    end

    PGB->>PGB: Reconnect zu P2

    App->>PGB: Write Request
    PGB->>P2: Forward to new Primary
    P2-->>App: Success

    ETCD->>Alert: Primary Failover Event
    Alert->>Alert: PagerDuty / Slack

    Note over App,P2: Failover komplett: < 30 Sekunden
```

### Failover-Szenarien im Detail

```mermaid
flowchart TB
    subgraph Scenarios["Ausfallszenarien"]
        S1[Primary Node Crash]
        S2[Netzwerk-Partition]
        S3[Datacenter-Ausfall]
        S4[Disk Corruption]
        S5[Overload / OOM]
    end

    subgraph Responses["Automatische Reaktion"]
        R1[Patroni Failover]
        R2[Split-Brain Prevention]
        R3[Cross-Region Failover]
        R4[Point-in-Time Recovery]
        R5[Connection Draining]
    end

    subgraph Recovery["Recovery-Zeit"]
        T1["< 30s"]
        T2["< 30s"]
        T3["< 5min"]
        T4["< 1h"]
        T5["< 1min"]
    end

    S1 --> R1 --> T1
    S2 --> R2 --> T2
    S3 --> R3 --> T3
    S4 --> R4 --> T4
    S5 --> R5 --> T5

    style Scenarios fill:#ffcccc
    style Responses fill:#90EE90
    style Recovery fill:#87CEEB
```

### Split-Brain Prevention

Das gefährlichste Szenario: Zwei Nodes glauben, sie sind Primary.

```mermaid
flowchart TB
    subgraph Problem["Split-Brain Problem"]
        P1_old[(Primary 1<br/>glaubt Leader)]
        P2_old[(Primary 2<br/>glaubt Leader)]
        Client1[Client A]
        Client2[Client B]

        Client1 -->|Write X=1| P1_old
        Client2 -->|Write X=2| P2_old
    end

    subgraph Solution["Lösung: etcd Quorum"]
        ETCD1[etcd 1]
        ETCD2[etcd 2]
        ETCD3[etcd 3]
        Leader[(Einziger Leader)]

        ETCD1 & ETCD2 & ETCD3 -->|"Mehrheit (2/3)<br/>entscheidet"| Leader
    end

    Problem -->|"❌ Ohne Consensus"| Conflict[Datenkonflikt!]
    Solution -->|"✅ Mit Consensus"| Safe[Konsistente Daten]

    style Problem fill:#ffcccc
    style Solution fill:#90EE90
    style Conflict fill:#ff0000,color:#fff
    style Safe fill:#00ff00
```

### Connection Pooling mit PgBouncer

```mermaid
flowchart TB
    subgraph Apps["100+ App Instances"]
        A1[Backend 1]
        A2[Backend 2]
        A3[Backend ...]
        AN[Backend N]
    end

    subgraph PgBouncer["PgBouncer Pool"]
        Pool["Connection Pool<br/>max 20 DB Connections"]

        subgraph Modes["Pool Modes"]
            Session[Session Pooling]
            Transaction[Transaction Pooling]
            Statement[Statement Pooling]
        end
    end

    subgraph Database["PostgreSQL"]
        DB[(Primary<br/>max_connections: 100)]
    end

    A1 & A2 & A3 & AN -->|"1000+ App Connections"| Pool
    Pool -->|"20 DB Connections"| DB

    style Apps fill:#87CEEB
    style PgBouncer fill:#FFD700
    style Database fill:#90EE90
```

| Pool Mode | Beschreibung | Use Case |
|-----------|--------------|----------|
| **Session** | Conn bleibt bei Client | Prepared Statements, Temp Tables |
| **Transaction** | Conn pro Transaction | Standard OLTP (empfohlen) |
| **Statement** | Conn pro Query | Simple Queries only |

### Backup & Point-in-Time Recovery

```mermaid
flowchart TB
    subgraph Continuous["Continuous Backup"]
        PG[(PostgreSQL)]
        WAL[WAL Archiving<br/>alle 5 Minuten]
        BaseBackup[Base Backup<br/>täglich]
    end

    subgraph Storage["Backup Storage"]
        S3_Local[S3 Primary Region]
        S3_Remote[S3 DR Region]
        Glacier[S3 Glacier<br/>90 Tage Retention]
    end

    subgraph Recovery["Recovery Optionen"]
        PITR[Point-in-Time Recovery<br/>bis zu 5 Min vor Ausfall]
        Clone[DB Clone<br/>für Testing]
        DR[Disaster Recovery<br/>neue Region]
    end

    PG --> WAL --> S3_Local
    PG --> BaseBackup --> S3_Local
    S3_Local --> S3_Remote
    S3_Local --> Glacier

    S3_Local --> PITR & Clone
    S3_Remote --> DR

    style Continuous fill:#90EE90
    style Storage fill:#FF9900
    style Recovery fill:#87CEEB
```

### RPO & RTO Targets

| Tier | RPO (max. Datenverlust) | RTO (max. Downtime) | Strategie |
|------|-------------------------|---------------------|-----------|
| **Enterprise** | 0 (zero loss) | < 1 min | Multi-Region Active-Active |
| **Business Critical** | < 1 min | < 15 min | Hot Standby + Auto-Failover |
| **Standard** | < 1 hour | < 4 hours | Warm Standby |
| **Archive** | < 24 hours | < 24 hours | Cold Backup |

### Monitoring der Datenbank

```mermaid
flowchart TB
    subgraph Metrics["Kritische Metriken"]
        Rep[Replication Lag]
        Conn[Connection Count]
        QPS[Queries per Second]
        Lock[Lock Waits]
        Disk[Disk Usage]
        CPU[CPU / Memory]
    end

    subgraph Alerts["Alert Schwellen"]
        A1["Lag > 10s → Warning"]
        A2["Lag > 60s → Critical"]
        A3["Connections > 80% → Warning"]
        A4["Disk > 85% → Critical"]
    end

    subgraph Actions["Automatische Aktionen"]
        Scale[Connection Pool skalieren]
        Failover[Failover triggern]
        Cleanup[Log Cleanup]
        Page[On-Call anrufen]
    end

    Rep --> A1 & A2
    Conn --> A3
    Disk --> A4

    A2 --> Failover
    A3 --> Scale
    A4 --> Cleanup & Page

    style Metrics fill:#87CEEB
    style Alerts fill:#FFD700
    style Actions fill:#90EE90
```

---

## Caching Layer (Redis)

```mermaid
flowchart TB
    subgraph RedisCluster["Redis Sentinel Cluster"]
        subgraph Master["Master"]
            RM[Redis Master]
        end

        subgraph Replicas["Replicas"]
            RR1[Redis Replica 1]
            RR2[Redis Replica 2]
        end

        subgraph Sentinels["Sentinel Nodes (Monitoring)"]
            S1[Sentinel 1]
            S2[Sentinel 2]
            S3[Sentinel 3]
        end
    end

    subgraph UseCases["Cache Use Cases"]
        Sessions[Session Store<br/>TTL: 24h]
        RateLimit[Rate Limiting<br/>TTL: 1min]
        QueryCache[Query Cache<br/>TTL: 5min]
        SSE[SSE Pub/Sub<br/>Real-time]
    end

    RM --> RR1 & RR2
    S1 & S2 & S3 -.->|Monitor| RM & RR1 & RR2

    Sessions & RateLimit & QueryCache & SSE --> RM

    style RedisCluster fill:#DC382D,color:#fff
```

---

## Disaster Recovery

### Backup-Strategie

```mermaid
flowchart TB
    subgraph Production["Production Region"]
        DB[(PostgreSQL)]
        Files[Object Storage]
    end

    subgraph Backups["Backup Strategy"]
        subgraph Continuous["Continuous"]
            WAL[WAL Archiving<br/>alle 5 min]
            Replication[Streaming Replication]
        end

        subgraph Periodic["Periodic"]
            Hourly[Hourly Snapshots]
            Daily[Daily Full Backup]
            Weekly[Weekly Archive]
        end
    end

    subgraph Storage["Backup Storage"]
        S3_Primary[S3 Primary Region]
        S3_DR[S3 DR Region]
        Glacier[S3 Glacier<br/>Long-term]
    end

    DB --> WAL --> S3_Primary
    DB --> Replication
    DB --> Hourly & Daily --> S3_Primary
    Weekly --> Glacier

    S3_Primary -.->|Cross-Region| S3_DR

    style Production fill:#90EE90
    style Storage fill:#FF9900
```

---

## Monitoring & Observability

### Observability Stack

```mermaid
flowchart TB
    subgraph Apps["Applications"]
        FE[Frontend]
        BE[Backend]
        WK[Workers]
    end

    subgraph Collection["Data Collection"]
        subgraph Metrics
            Prom[Prometheus]
            Node[Node Exporter]
            PGExp[PG Exporter]
        end

        subgraph Logs
            Loki[Grafana Loki]
            Promtail[Promtail]
        end

        subgraph Traces
            Tempo[Grafana Tempo]
            OTEL[OpenTelemetry]
        end
    end

    subgraph Visualization["Visualization & Alerting"]
        Graf[Grafana]
        Alert[Alertmanager]
        PD[PagerDuty]
        Slack[Slack]
    end

    FE & BE & WK -->|metrics| Prom
    FE & BE & WK -->|logs| Promtail --> Loki
    FE & BE & WK -->|traces| OTEL --> Tempo

    Node & PGExp --> Prom

    Prom & Loki & Tempo --> Graf
    Prom --> Alert --> PD & Slack

    style Collection fill:#FF6B6B
    style Visualization fill:#4ECDC4
```

---

## Security in HA Environment

### Zero Trust Architecture

```mermaid
flowchart TB
    subgraph Internet
        U[User]
        Attacker[Attacker]
    end

    subgraph Edge["Edge Security"]
        WAF[WAF]
        DDoS[DDoS Shield]
        Bot[Bot Detection]
    end

    subgraph Network["Network Security"]
        VPC[Private VPC]
        SG[Security Groups]
        NL[Network Policies]
    end

    subgraph Identity["Identity & Access"]
        IAM[IAM / RBAC]
        SSO[SSO]
        MFA[MFA]
    end

    subgraph Data["Data Security"]
        Encrypt[Encryption at Rest]
        TLS[TLS 1.3]
        Secrets[Vault / Secrets Manager]
    end

    U --> WAF --> VPC
    Attacker --x WAF

    VPC --> SG --> NL
    NL --> Encrypt --> TLS

    style Edge fill:#FF6B6B
    style Network fill:#4ECDC4
    style Identity fill:#45B7D1
    style Data fill:#96CEB4
```

---

## Migration Path

### Von Current State zu HA

```mermaid
gantt
    title HA Migration Roadmap
    dateFormat YYYY-MM

    section Phase 1 - Foundation
    Redis Cache einführen          :p1a, 2024-01, 1M
    Load Balancer hinzufügen       :p1b, after p1a, 1M
    PostgreSQL Replica aufsetzen   :p1c, after p1b, 2M

    section Phase 2 - Scaling
    Multi-Node Deployment          :p2a, after p1c, 2M
    Kubernetes Migration           :p2b, after p2a, 3M
    Auto-Scaling einrichten        :p2c, after p2b, 1M

    section Phase 3 - Multi-Region
    Zweite Region aufsetzen        :p3a, after p2c, 2M
    Cross-Region Replication       :p3b, after p3a, 2M
    Global Load Balancer           :p3c, after p3b, 1M
```

### Kosten-Übersicht

| Phase | Beschreibung | Kosten/Monat |
|-------|--------------|--------------|
| Current | Single Server | ~15 EUR |
| Phase 1 | Regional HA (Basic) | ~150 EUR |
| Phase 2 | Regional HA (Full) | ~800 EUR |
| Phase 3 | Multi-Region (2) | ~2500 EUR |
| Phase 4 | Multi-Region (3) | ~4500 EUR |

---

## Entscheidungsmatrix: Build vs Buy

| Komponente | Empfehlung | Grund |
|------------|------------|-------|
| **Load Balancing** | Buy (Cloudflare/AWS) | Commodity, günstig |
| **Database HA** | Managed (RDS/Cloud SQL) | Kritisch, komplex |
| **Monitoring** | SaaS (Grafana Cloud) | Focus on Core Business |
| **CDN** | Buy (Cloudflare) | Global PoPs nötig |
| **Auth** | Build (Better Auth) | Schon implementiert |
| **Kubernetes** | Managed (EKS/GKE) | Operational overhead |

---

## Zusammenfassung

### Key Takeaways

1. **Stateless Application Layer** - Horizontale Skalierung durch stateless Design
2. **Database als Single Source of Truth** - PostgreSQL mit HA-Cluster ist kritisch
3. **Caching für Performance** - Redis für Sessions und häufige Queries
4. **Event-Driven Architecture** - River für Jobs, optional Kafka für Events
5. **Defense in Depth** - Mehrere Security-Layer
6. **Observability First** - Metrics, Logs, Traces von Anfang an

### Nächste Schritte

1. [ ] Redis Cache einführen (Sessions, Rate Limiting)
2. [ ] PostgreSQL Replica aufsetzen
3. [ ] Load Balancer vor App-Server
4. [ ] Kubernetes Migration evaluieren
5. [ ] Disaster Recovery Plan dokumentieren
