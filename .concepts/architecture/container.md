# Container Diagram (C4 Level 2)

## Technische Architektur

Detaillierte Sicht auf die Container (deploybare Einheiten) des Systems.

## Container Diagram

```mermaid
C4Container
    title GocaTest - Container Diagram

    Person(user, "User", "Benutzer der Anwendung")

    System_Boundary(gocatest, "GocaTest System") {
        Container(proxy, "kamal-proxy", "Go", "Reverse Proxy, SSL Termination, Zero-Downtime Deploys")
        Container(frontend, "Frontend", "Next.js 16, TypeScript", "React SPA mit App Router, TanStack Query")
        Container(backend, "Backend", "Go, Gorilla Mux", "REST API, Clean Architecture, SSE")
        ContainerDb(db, "Database", "PostgreSQL 16", "User Data, Sessions, Business Data")
    }

    System_Ext(email, "Email Service", "SendGrid/SMTP")

    Rel(user, proxy, "HTTPS :443")
    Rel(proxy, frontend, "HTTP :3000")
    Rel(proxy, backend, "HTTP :8080")
    Rel(frontend, backend, "REST API", "JSON")
    Rel(backend, db, "SQL", "TCP :5432")
    Rel(backend, email, "SMTP/API")

    UpdateLayoutConfig($c4ShapeInRow="3", $c4BoundaryInRow="1")
```

## Container Details

### Frontend (Next.js)

```mermaid
flowchart TB
    subgraph Frontend["Frontend Container"]
        direction TB
        AR[App Router]
        RC[React Components]
        TQ[TanStack Query]
        BA[Better Auth Client]
        OR[Orval API Client]

        AR --> RC
        RC --> TQ
        RC --> BA
        TQ --> OR
    end
```

| Komponente | Technologie | Verantwortung |
|------------|-------------|---------------|
| App Router | Next.js 16 | Routing, SSR, Server Components |
| Components | React + shadcn/ui | UI Rendering |
| TanStack Query | React Query | Server State, Caching |
| Better Auth | Auth Client | Session Management |
| Orval | Generated Hooks | Type-safe API Calls |

### Backend (Go)

```mermaid
flowchart TB
    subgraph Backend["Backend Container"]
        direction TB
        H[Handler Layer]
        U[UseCase Layer]
        R[Repository Layer]
        D[Domain Layer]
        M[Middleware]
        SSE[SSE Broker]

        M --> H
        H --> U
        U --> R
        R --> D
        H --> SSE
    end
```

| Layer | Verantwortung | Goca Command |
|-------|---------------|--------------|
| Handler | HTTP Endpoints, Swagger | `goca make handler` |
| UseCase | Business Logic | `goca make usecase` |
| Repository | Data Access | `goca make repository` |
| Domain | Entities, Rules | `goca make entity` |

### Database Schema

```mermaid
erDiagram
    users {
        uuid id PK
        string email UK
        string name
        timestamp created_at
    }

    sessions {
        uuid id PK
        uuid user_id FK
        string token
        timestamp expires_at
    }

    user_stats {
        uuid id PK
        uuid user_id FK
        int projects
        int tasks
        int completed
    }

    users ||--o{ sessions : has
    users ||--o| user_stats : has
```

## Kommunikation

### Synchrone Kommunikation

```mermaid
sequenceDiagram
    participant U as User
    participant P as Proxy
    participant F as Frontend
    participant B as Backend
    participant D as Database

    U->>P: GET /dashboard
    P->>F: Forward Request
    F->>F: Server Component
    F->>B: GET /api/v1/stats
    B->>D: SELECT stats
    D-->>B: Stats Data
    B-->>F: JSON Response
    F-->>P: HTML + Hydration
    P-->>U: Complete Page
```

### Asynchrone Kommunikation (SSE)

```mermaid
sequenceDiagram
    participant U as User
    participant F as Frontend
    participant B as Backend

    U->>F: Open Dashboard
    F->>B: GET /api/v1/events (SSE)
    Note over F,B: Connection kept open

    loop Real-time Updates
        B-->>F: event: stats-updated
        F->>F: Invalidate Query
        F-->>U: UI Update
    end
```
