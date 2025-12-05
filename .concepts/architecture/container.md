# Container Diagram (C4 Level 2)

## Technical Architecture

Detailed view of the containers (deployable units) of the system.

## Container Diagram

```mermaid
C4Container
    title Next-Go-PG - Container Diagram

    Person(user, "User", "Application user")
    Person(admin, "Admin", "Via VPN")

    System_Boundary(nextgopg, "Next-Go-PG System") {
        Container(proxy, "kamal-proxy", "Go", "Reverse Proxy, SSL Termination, Zero-Downtime Deploys")
        Container(frontend, "Frontend", "Next.js 16, TypeScript", "React SPA mit App Router, TanStack Query")
        Container(backend, "Backend", "Go, Gorilla Mux", "REST API, Clean Architecture, SSE")
        ContainerDb(db, "Database", "PostgreSQL 16", "User Data, Sessions, Business Data")
        Container(loki, "Loki", "Grafana Loki", "Log Aggregation (VPN only)")
        Container(grafana, "Grafana", "Grafana 10", "Log Visualization (VPN only)")
    }

    System_Ext(email, "Email Service", "SendGrid/SMTP")
    System_Ext(vpn, "VPN", "Private Network Access")

    Rel(user, proxy, "HTTPS :443")
    Rel(proxy, frontend, "HTTP :3000")
    Rel(proxy, backend, "HTTP :8080")
    Rel(frontend, backend, "REST API", "JSON")
    Rel(backend, db, "SQL", "TCP :5432")
    Rel(backend, email, "SMTP/API")
    Rel(backend, loki, "HTTP POST", "Logs")
    Rel(frontend, loki, "HTTP POST", "Logs")
    Rel(loki, grafana, "LogQL")
    Rel(admin, vpn, "VPN")
    Rel(vpn, grafana, "HTTPS :3001")

    UpdateLayoutConfig($c4ShapeInRow="3", $c4BoundaryInRow="1")
```

## Container Details

### Frontend (Next.js + FSD)

```mermaid
flowchart TB
    subgraph Frontend["Frontend Container - Feature-Sliced Design"]
        direction TB
        subgraph App["app/"]
            AR[App Router]
        end

        subgraph Widgets["widgets/"]
            Header[header]
        end

        subgraph Features["features/"]
            Auth[auth]
            Stats[stats]
        end

        subgraph Entities["entities/"]
            User[user]
        end

        subgraph Shared["shared/"]
            UI[ui - shadcn]
            API[api - Orval]
            Lib[lib - Utils]
        end

        AR --> Widgets
        AR --> Features
        Widgets --> Features
        Widgets --> Shared
        Features --> Entities
        Features --> Shared
        Entities --> Shared
    end
```

| Layer | Technology | Responsibility |
|-------|------------|----------------|
| app/ | Next.js 16 App Router | Routing, SSR, Pages |
| widgets/ | React Components | Composite UI (Header) |
| features/ | React + Hooks | User Interactions (Auth, Stats) |
| entities/ | TypeScript | Business Objects (User) |
| shared/ui | shadcn/ui | UI Components |
| shared/api | Orval + TanStack Query | Type-safe API Calls |
| shared/lib | Better Auth, Utils | Auth, Helpers |

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

| Layer | Responsibility | Goca Command |
|-------|----------------|--------------|
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

## Communication

### Synchronous Communication

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

### Asynchronous Communication (SSE)

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
