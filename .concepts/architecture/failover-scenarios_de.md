# Failover-Szenarien: Konkrete Fallbeispiele

Dieses Dokument zeigt anhand realer Szenarien, wie die HA-Architektur Ausfälle abfängt.

---

## Szenario 1: Datenbank-Server stürzt ab

**Situation:** Der Primary PostgreSQL-Server hat um 03:14 Uhr einen Kernel Panic und fällt komplett aus.

### Timeline des Ausfalls

```mermaid
sequenceDiagram
    participant U as User
    participant App as Backend
    participant PGB as PgBouncer
    participant P1 as Primary DB
    participant P2 as Sync Replica
    participant ETCD as etcd Cluster
    participant Alert as Alerting

    Note over P1: 03:14:00 - Kernel Panic!

    rect rgb(255, 200, 200)
        Note over P1: Server nicht erreichbar
        App->>PGB: INSERT INTO orders...
        PGB->>P1: Forward Query
        P1--xPGB: Connection refused
        PGB-->>App: Error: Connection lost
    end

    Note over ETCD: 03:14:03 - Health Check failed (3x)

    rect rgb(255, 255, 200)
        ETCD->>ETCD: Leader Election startet
        ETCD->>P2: Du bist neuer Leader!
        P2->>P2: Promote to Primary (5s)
        P2->>ETCD: Registriert als Primary
        ETCD->>PGB: Routing Update
        PGB->>PGB: Reconnect zu P2
    end

    Note over P2: 03:14:22 - Neuer Primary aktiv

    rect rgb(200, 255, 200)
        App->>PGB: INSERT INTO orders... (retry)
        PGB->>P2: Forward to new Primary
        P2-->>App: Success!
        App-->>U: Bestellung erfolgreich
    end

    ETCD->>Alert: CRITICAL: Primary Failover
    Alert->>Alert: PagerDuty + Slack
```

### Was passiert auf jeder Ebene?

```mermaid
flowchart TB
    subgraph Timeline["Zeitlicher Ablauf"]
        T0["03:14:00<br/>Ausfall"]
        T1["03:14:03<br/>Erkannt"]
        T2["03:14:08<br/>Election"]
        T3["03:14:15<br/>Promotion"]
        T4["03:14:22<br/>Wiederhergestellt"]
    end

    subgraph Actions["Automatische Aktionen"]
        A1[Patroni erkennt<br/>fehlenden Heartbeat]
        A2[etcd startet<br/>Leader Election]
        A3[Sync Replica wird<br/>zu Primary promoted]
        A4[PgBouncer routet<br/>zu neuem Primary]
    end

    T0 --> T1 --> T2 --> T3 --> T4
    T1 --- A1
    T2 --- A2
    T3 --- A3
    T4 --- A4

    style T0 fill:#ff6666
    style T4 fill:#66ff66
```

### User Experience

| Ohne HA | Mit HA |
|---------|--------|
| Kompletter Ausfall | 22 Sekunden Unterbrechung |
| Manuelles Eingreifen nötig | Automatisches Failover |
| Potentieller Datenverlust | 0 Datenverlust (sync replication) |
| Downtime: Stunden | Downtime: < 30 Sekunden |

---

## Szenario 2: Netzwerk-Partition (Split-Brain Gefahr)

**Situation:** Ein Netzwerkfehler trennt den Primary vom Rest des Clusters.

### Das Problem ohne Schutz

```mermaid
flowchart TB
    subgraph Danger["GEFAHR: Split-Brain"]
        subgraph Network1["Netzwerk A"]
            Client1[Client A]
            P1_bad[(Primary<br/>isoliert)]
        end

        subgraph Network2["Netzwerk B"]
            Client2[Client B]
            P2_bad[(Replica<br/>promoted)]
        end

        Client1 -->|"Write: balance=100"| P1_bad
        Client2 -->|"Write: balance=50"| P2_bad

        P1_bad -.->|"❌ Keine Verbindung"| P2_bad
    end

    subgraph Result["Resultat"]
        Conflict[Datenkonflikt!<br/>Welcher Wert ist korrekt?]
    end

    Danger --> Conflict

    style Danger fill:#ffcccc
    style Conflict fill:#ff0000,color:#fff
```

### Die Lösung: etcd Quorum + Fencing

```mermaid
sequenceDiagram
    participant P1 as Primary (isoliert)
    participant P2 as Sync Replica
    participant E1 as etcd Node 1
    participant E2 as etcd Node 2
    participant E3 as etcd Node 3
    participant App as Application

    Note over P1,E3: Netzwerk-Partition tritt auf

    rect rgb(255, 200, 200)
        P1-xE1: Heartbeat
        P1-xE2: Heartbeat
        P1-xE3: Heartbeat
        Note over P1: Kann etcd nicht erreichen
    end

    rect rgb(200, 255, 200)
        P2->>E1: Heartbeat OK
        P2->>E2: Heartbeat OK
        P2->>E3: Heartbeat OK
        Note over E1,E3: Quorum: 3/3 sehen P2
    end

    E1->>E2: P1 nicht erreichbar
    E2->>E3: Bestätigt
    Note over E1,E3: Mehrheit entscheidet: P1 ist tot

    E1->>P2: Du bist neuer Leader!
    P2->>P2: Promote to Primary

    rect rgb(255, 255, 200)
        Note over P1: FENCING: P1 erkennt<br/>keine etcd Verbindung
        P1->>P1: SELBST-DEGRADIERUNG!
        P1->>P1: Akzeptiere keine Writes mehr
    end

    App->>P2: Write Request
    P2-->>App: Success (einziger Primary)

    Note over P1,App: Kein Split-Brain!<br/>Nur EIN Primary aktiv
```

### Fencing-Mechanismus im Detail

```mermaid
flowchart TB
    subgraph Primary["Isolierter Primary"]
        Check[Patroni prüft etcd]
        Decision{Kann etcd<br/>erreichen?}
        Continue[Weiter als Primary]
        Demote[SELBST-DEGRADIERUNG]
        ReadOnly[Nur noch Read-Only]
    end

    subgraph Protection["Schutzmaßnahmen"]
        NoWrite[Writes werden abgelehnt]
        NoLease[Leader-Lease läuft ab]
        Alert[Alert: Primary degradiert]
    end

    Check --> Decision
    Decision -->|Ja| Continue
    Decision -->|Nein| Demote
    Demote --> ReadOnly
    ReadOnly --> NoWrite & NoLease & Alert

    style Demote fill:#FFD700
    style ReadOnly fill:#87CEEB
```

---

## Szenario 3: Deployment schlägt fehl

**Situation:** Ein neues Backend-Deployment hat einen Bug und die App crasht beim Start.

### Rolling Deployment mit Health Checks

```mermaid
sequenceDiagram
    participant K as Kamal/K8s
    participant P as Proxy/LB
    participant Old as Backend v1.0
    participant New as Backend v1.1
    participant HC as Health Check

    Note over K: Deployment v1.1 startet

    K->>New: Container starten
    New->>New: App initialisieren

    loop Health Check (alle 5s)
        HC->>New: GET /health
        New--xHC: 500 Error (Bug!)
    end

    Note over New: 3x Health Check failed

    rect rgb(255, 200, 200)
        K->>K: Deployment FAILED
        K->>New: Container stoppen
    end

    rect rgb(200, 255, 200)
        Note over Old: v1.0 läuft weiter!
        P->>Old: Traffic weiterhin zu v1.0
    end

    K->>K: Rollback zu v1.0
    Note over K,Old: Kein Ausfall für User!
```

### Zero-Downtime Deployment (Happy Path)

```mermaid
flowchart TB
    subgraph Phase1["Phase 1: Vorbereitung"]
        Pull[Image pullen]
        Start[Neuen Container starten]
    end

    subgraph Phase2["Phase 2: Validierung"]
        HC1[Health Check 1]
        HC2[Health Check 2]
        HC3[Health Check 3]
        Ready{Alle Checks OK?}
    end

    subgraph Phase3["Phase 3: Umschaltung"]
        Route[Traffic umleiten]
        Drain[Alte Connections draining]
        Stop[Alten Container stoppen]
    end

    subgraph Rollback["Rollback Path"]
        Fail[Deployment abbrechen]
        Keep[Alten Container behalten]
    end

    Pull --> Start --> HC1 --> HC2 --> HC3 --> Ready
    Ready -->|Ja| Route --> Drain --> Stop
    Ready -->|Nein| Fail --> Keep

    style Phase3 fill:#90EE90
    style Rollback fill:#FFB6C1
```

---

## Szenario 4: DDoS-Angriff

**Situation:** Die Anwendung wird mit 10 Millionen Requests/Sekunde angegriffen.

### Schutzschichten

```mermaid
flowchart TB
    subgraph Attack["DDoS Angriff"]
        Bot1[Botnet 1<br/>1M req/s]
        Bot2[Botnet 2<br/>3M req/s]
        Bot3[Botnet 3<br/>6M req/s]
    end

    subgraph Layer1["Layer 1: Cloudflare Edge"]
        CF[Cloudflare DDoS Shield]
        Block1[99% geblockt<br/>bekannte Botnets]
    end

    subgraph Layer2["Layer 2: WAF"]
        WAF[Web Application Firewall]
        Block2[Rate Limiting<br/>100 req/s pro IP]
    end

    subgraph Layer3["Layer 3: Application"]
        LB[Load Balancer]
        Redis[Redis Rate Limit]
        Block3[User Rate Limit<br/>10 req/s pro User]
    end

    subgraph Protected["Geschützte App"]
        App[Backend]
        DB[(Database)]
    end

    Bot1 & Bot2 & Bot3 -->|"10M req/s"| CF
    CF -->|"100K req/s"| WAF
    WAF -->|"10K req/s"| LB
    LB -->|"1K req/s"| App
    App --> DB

    Block1 -.- CF
    Block2 -.- WAF
    Block3 -.- Redis

    style Attack fill:#ff6666
    style Layer1 fill:#FF9900
    style Layer2 fill:#FFD700
    style Layer3 fill:#87CEEB
    style Protected fill:#90EE90
```

### Rate Limiting Kaskade

```mermaid
flowchart LR
    subgraph Levels["Rate Limit Ebenen"]
        L1["Edge: 1000 req/s/IP"]
        L2["WAF: 100 req/s/IP"]
        L3["App: 10 req/s/User"]
        L4["API: 1 req/s/Endpoint"]
    end

    subgraph Response["Antwort bei Überschreitung"]
        R1["403 Forbidden<br/>IP geblockt"]
        R2["429 Too Many<br/>Retry-After: 60s"]
        R3["429 Too Many<br/>Retry-After: 10s"]
        R4["429 Too Many<br/>Retry-After: 1s"]
    end

    L1 --> R1
    L2 --> R2
    L3 --> R3
    L4 --> R4
```

---

## Szenario 5: Datacenter-Ausfall

**Situation:** Das gesamte Datacenter in Frankfurt fällt aus (Stromausfall, Naturkatastrophe).

### Multi-Region Failover

```mermaid
sequenceDiagram
    participant User as User (EU)
    participant DNS as Cloudflare DNS
    participant EU as EU Region (Frankfurt)
    participant US as US Region (Virginia)
    participant Monitor as Health Monitor

    Note over EU: ❌ Datacenter Ausfall!

    loop Health Check (alle 30s)
        Monitor->>EU: GET /health
        EU--xMonitor: Timeout
    end

    Monitor->>Monitor: 3x Timeout = Region Down

    rect rgb(255, 200, 200)
        Monitor->>DNS: EU Region unhealthy
        DNS->>DNS: Entferne EU aus Rotation
    end

    rect rgb(200, 255, 200)
        User->>DNS: app.example.com
        DNS-->>User: US Region IP
        User->>US: Request
        US-->>User: Response (höhere Latenz)
    end

    Note over User,US: Service verfügbar!<br/>Latenz: +80ms
```

### GeoDNS Failover

```mermaid
flowchart TB
    subgraph Normal["Normalbetrieb"]
        EU_User[EU User] -->|"app.example.com"| DNS1[DNS]
        DNS1 -->|"Frankfurt IP"| EU_DC[EU Datacenter]

        US_User[US User] -->|"app.example.com"| DNS2[DNS]
        DNS2 -->|"Virginia IP"| US_DC[US Datacenter]
    end

    subgraph Failover["Nach EU-Ausfall"]
        EU_User2[EU User] -->|"app.example.com"| DNS3[DNS]
        DNS3 -->|"Virginia IP<br/>(nächste Region)"| US_DC2[US Datacenter]

        US_User2[US User] -->|"app.example.com"| DNS4[DNS]
        DNS4 -->|"Virginia IP"| US_DC3[US Datacenter]
    end

    style EU_DC fill:#90EE90
    style US_DC fill:#90EE90
    style US_DC2 fill:#FFD700
    style US_DC3 fill:#90EE90
```

### Daten-Synchronisation zwischen Regionen

```mermaid
flowchart TB
    subgraph EU["EU Region (Primary)"]
        EU_App[App Cluster]
        EU_DB[(PostgreSQL<br/>Primary)]
        EU_Redis[Redis]
    end

    subgraph US["US Region (Secondary)"]
        US_App[App Cluster]
        US_DB[(PostgreSQL<br/>Read Replica)]
        US_Redis[Redis]
    end

    subgraph Sync["Synchronisation"]
        DBSync[PostgreSQL<br/>Async Replication<br/>Lag: 100-500ms]
        CacheSync[Redis<br/>Pub/Sub Sync]
    end

    EU_DB -->|WAL Stream| DBSync -->|Apply| US_DB
    EU_Redis <-->|Sync| CacheSync <-->|Sync| US_Redis

    EU_App --> EU_DB
    EU_App --> EU_Redis

    US_App -->|Read| US_DB
    US_App -->|Write via EU| EU_DB
    US_App --> US_Redis

    style EU fill:#90EE90
    style US fill:#87CEEB
```

---

## Szenario 6: Datenkorruption durch Bug

**Situation:** Ein Bug im Code löscht versehentlich User-Daten.

### Point-in-Time Recovery

```mermaid
sequenceDiagram
    participant Dev as Developer
    participant App as Application
    participant DB as Database
    participant WAL as WAL Archive
    participant S3 as S3 Backup

    Note over App: 14:30 - Bug deployed
    App->>DB: DELETE FROM users WHERE active=true
    Note over DB: 😱 50.000 User gelöscht!

    DB->>WAL: WAL Segment geschrieben
    WAL->>S3: Archiviert

    Note over Dev: 15:45 - Bug entdeckt!

    rect rgb(255, 255, 200)
        Dev->>Dev: Entscheidung: PITR zu 14:29
    end

    rect rgb(200, 255, 200)
        Dev->>S3: Hole Base Backup
        S3-->>Dev: Backup von 00:00
        Dev->>S3: Hole WAL Segmente
        S3-->>Dev: WAL 00:00 - 14:29

        Dev->>Dev: Restore + Replay WAL
        Dev->>Dev: Stoppe bei 14:29:59
    end

    Note over Dev,DB: Daten wiederhergestellt!<br/>Verlust: nur Änderungen 14:29-15:45
```

### Recovery-Optionen

```mermaid
flowchart TB
    subgraph Problem["Problem erkannt"]
        Bug[Datenkorruption<br/>durch Bug]
    end

    subgraph Options["Recovery-Optionen"]
        PITR[Point-in-Time Recovery<br/>Restore bis Zeitpunkt X]
        Clone[Clone + Repair<br/>Daten extrahieren]
        Replica[Replica nutzen<br/>wenn noch nicht repliziert]
    end

    subgraph Outcome["Ergebnis"]
        Full[Vollständige<br/>Wiederherstellung]
        Partial[Partielle<br/>Wiederherstellung]
        Manual[Manuelle<br/>Korrektur]
    end

    Bug --> PITR & Clone & Replica
    PITR --> Full
    Clone --> Partial
    Replica --> Full

    style Problem fill:#ff6666
    style Full fill:#90EE90
    style Partial fill:#FFD700
```

---

## Szenario 7: Memory Leak / OOM

**Situation:** Ein Memory Leak führt dazu, dass der Backend-Pod OOM (Out of Memory) killed wird.

### Kubernetes Self-Healing

```mermaid
sequenceDiagram
    participant K8s as Kubernetes
    participant Pod1 as Backend Pod 1
    participant Pod2 as Backend Pod 2
    participant Pod3 as Backend Pod 3
    participant LB as Load Balancer

    Note over Pod1: Memory: 80%... 90%... 95%...

    rect rgb(255, 200, 200)
        K8s->>Pod1: OOM Kill!
        Pod1--xLB: Unhealthy
        LB->>LB: Entferne Pod1 aus Rotation
    end

    Note over LB: Traffic nur zu Pod2 & Pod3

    rect rgb(200, 255, 200)
        K8s->>K8s: Restart Policy: Always
        K8s->>Pod1: Container neu starten
        Pod1->>Pod1: Initialisierung...
        Pod1->>K8s: Health Check OK
        K8s->>LB: Pod1 wieder healthy
        LB->>LB: Füge Pod1 zur Rotation
    end

    Note over Pod1,Pod3: Alle 3 Pods wieder aktiv
```

### Resource Limits & Monitoring

```mermaid
flowchart TB
    subgraph Limits["Resource Limits"]
        CPU["CPU Limit: 2 cores"]
        Mem["Memory Limit: 2Gi"]
        Req["Request: 500m / 512Mi"]
    end

    subgraph Monitoring["Proaktives Monitoring"]
        Alert1["Memory > 70%<br/>→ Warning"]
        Alert2["Memory > 85%<br/>→ Critical"]
        Alert3["Memory > 95%<br/>→ Auto-Scale"]
    end

    subgraph Actions["Automatische Aktionen"]
        Scale[HPA: Neue Pods starten]
        Restart[Graceful Restart]
        Page[On-Call benachrichtigen]
    end

    Limits --> Monitoring
    Alert1 --> Page
    Alert2 --> Restart
    Alert3 --> Scale

    style Alert1 fill:#FFD700
    style Alert2 fill:#FF9900
    style Alert3 fill:#ff6666
```

---

## Szenario 8: SSL-Zertifikat läuft ab

**Situation:** Das SSL-Zertifikat läuft in 7 Tagen ab.

### Automatische Zertifikatserneuerung

```mermaid
sequenceDiagram
    participant Monitor as Cert Monitor
    participant Alert as Alerting
    participant LE as Let's Encrypt
    participant Proxy as Proxy/Ingress

    Note over Monitor: Täglicher Check

    Monitor->>Proxy: Prüfe Zertifikat
    Proxy-->>Monitor: Gültig bis: +7 Tage

    rect rgb(255, 255, 200)
        Monitor->>Alert: WARNING: Cert expires in 7 days
        Alert->>Alert: Slack Notification
    end

    Note over Monitor: Automatische Erneuerung (< 30 Tage)

    rect rgb(200, 255, 200)
        Monitor->>LE: Request new certificate
        LE->>LE: ACME Challenge
        LE-->>Monitor: New certificate
        Monitor->>Proxy: Install certificate
        Proxy->>Proxy: Hot-Reload (0 downtime)
    end

    Monitor->>Alert: INFO: Certificate renewed
```

### Zertifikats-Monitoring Dashboard

```mermaid
flowchart TB
    subgraph Certs["Zertifikate"]
        Main["app.example.com<br/>Ablauf: 89 Tage ✅"]
        API["api.example.com<br/>Ablauf: 89 Tage ✅"]
        Wild["*.example.com<br/>Ablauf: 25 Tage ⚠️"]
    end

    subgraph Thresholds["Schwellenwerte"]
        T1["< 30 Tage: Auto-Renew"]
        T2["< 14 Tage: Warning"]
        T3["< 7 Tage: Critical"]
        T4["< 1 Tag: EMERGENCY"]
    end

    Certs --> Thresholds

    style Main fill:#90EE90
    style API fill:#90EE90
    style Wild fill:#FFD700
```

---

## Zusammenfassung: Schutzmatrix

```mermaid
flowchart TB
    subgraph Threats["Bedrohungen"]
        T1[Server Crash]
        T2[Netzwerk-Partition]
        T3[Bad Deployment]
        T4[DDoS Angriff]
        T5[Datacenter-Ausfall]
        T6[Datenkorruption]
        T7[Memory Leak]
        T8[Cert Expiry]
    end

    subgraph Protection["Schutzmaßnahmen"]
        P1[Patroni Auto-Failover]
        P2[etcd Quorum + Fencing]
        P3[Health Checks + Rollback]
        P4[CDN + WAF + Rate Limit]
        P5[Multi-Region + GeoDNS]
        P6[PITR + WAL Archiving]
        P7[K8s Self-Healing + HPA]
        P8[Auto-Renewal + Monitoring]
    end

    subgraph Recovery["Recovery Zeit"]
        R1["< 30s"]
        R2["< 30s"]
        R3["0s (kein Ausfall)"]
        R4["0s (transparent)"]
        R5["< 5min"]
        R6["< 1h"]
        R7["< 1min"]
        R8["0s (auto)"]
    end

    T1 --> P1 --> R1
    T2 --> P2 --> R2
    T3 --> P3 --> R3
    T4 --> P4 --> R4
    T5 --> P5 --> R5
    T6 --> P6 --> R6
    T7 --> P7 --> R7
    T8 --> P8 --> R8

    style Threats fill:#ffcccc
    style Protection fill:#90EE90
    style Recovery fill:#87CEEB
```

| Szenario | Schutz | Recovery | Datenverlust |
|----------|--------|----------|--------------|
| Server Crash | Patroni Failover | < 30s | 0 |
| Netzwerk-Partition | etcd Quorum | < 30s | 0 |
| Bad Deployment | Health Checks | 0s | 0 |
| DDoS Angriff | CDN + WAF | 0s | 0 |
| Datacenter-Ausfall | Multi-Region | < 5min | < 500ms Lag |
| Datenkorruption | PITR | < 1h | Bis zum Bug |
| Memory Leak | K8s Self-Healing | < 1min | 0 |
| Cert Expiry | Auto-Renewal | 0s | 0 |
