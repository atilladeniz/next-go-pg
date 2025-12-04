# Data Flow Architecture

## Overview

Datenfluss durch die Anwendung mit TanStack Query, HydrationBoundary und SSE.

## Request/Response Flow

```mermaid
flowchart TB
    subgraph Client["Client (Browser)"]
        UC[User Action]
        TQ[TanStack Query]
        OR[Orval Hook]
        CA[Cache]
    end

    subgraph Server["Server"]
        SC[Server Component]
        QC[QueryClient]
        HB[HydrationBoundary]
    end

    subgraph Backend["Backend API"]
        H[Handler]
        U[UseCase]
        R[Repository]
        DB[(PostgreSQL)]
    end

    UC --> TQ
    TQ --> OR
    OR -->|Fetch| H
    H --> U
    U --> R
    R --> DB
    DB --> R
    R --> U
    U --> H
    H -->|JSON| OR
    OR --> CA
    CA --> TQ

    SC --> QC
    QC -->|prefetchQuery| H
    QC --> HB
    HB -->|Hydrate| CA
```

## Server-Side Prefetching Pattern

```mermaid
sequenceDiagram
    participant B as Browser
    participant SC as Server Component
    participant QC as QueryClient
    participant API as Backend API
    participant CC as Client Component

    B->>SC: Request Page
    SC->>SC: getSession()
    SC->>QC: getQueryClient()
    SC->>API: prefetchQuery (with Cookie)
    API-->>QC: Data cached
    SC->>SC: dehydrate(queryClient)
    SC-->>B: HTML + Dehydrated State

    B->>CC: Hydrate Client
    CC->>CC: useQuery() - Cache Hit!
    Note over CC: Kein Loading State!
```

## Real-Time Updates (SSE)

```mermaid
sequenceDiagram
    participant U as User Action
    participant B1 as Browser 1
    participant B2 as Browser 2
    participant API as Backend
    participant SSE as SSE Broker
    participant DB as Database

    U->>B1: Create Item
    B1->>API: POST /api/v1/items
    API->>DB: INSERT item
    API->>SSE: Broadcast "items-updated"
    API-->>B1: 201 Created

    par Parallel Notification
        SSE-->>B1: event: items-updated
        SSE-->>B2: event: items-updated
    end

    B1->>B1: invalidateQueries(['items'])
    B2->>B2: invalidateQueries(['items'])

    Note over B1,B2: Beide sehen Update sofort
```

## Cache Invalidation Strategy

```mermaid
flowchart TD
    A[Data Mutation] --> B{Mutation Type}

    B -->|Create| C[Invalidate List Query]
    B -->|Update| D[Invalidate List + Detail Query]
    B -->|Delete| E[Invalidate List Query + Remove from Cache]

    C --> F[SSE Broadcast]
    D --> F
    E --> F

    F --> G[All Clients Invalidate]
    G --> H[Refetch Fresh Data]
```

## Query Key Structure

```mermaid
mindmap
    root((Query Keys))
        users
            getUsers
            getUser
                userId
            getUserStats
        items
            getItems
                filters
            getItem
                itemId
        auth
            getSession
```

## Optimistic Updates

```mermaid
sequenceDiagram
    participant U as User
    participant TQ as TanStack Query
    participant CA as Cache
    participant API as Backend

    U->>TQ: Update Item
    TQ->>CA: Optimistic Update
    CA-->>U: UI Updates Immediately

    TQ->>API: PUT /api/v1/items/:id

    alt Success
        API-->>TQ: 200 OK
        TQ->>CA: Confirm Update
    else Error
        API-->>TQ: Error
        TQ->>CA: Rollback to Previous
        CA-->>U: Show Error + Reverted UI
    end
```

## Data Transformation Pipeline

```mermaid
flowchart LR
    subgraph Backend
        DB[(PostgreSQL)]
        R[Repository]
        U[UseCase]
        H[Handler]
    end

    subgraph Frontend
        OR[Orval Client]
        TQ[TanStack Query]
        C[Component]
    end

    DB -->|SQL Result| R
    R -->|Domain Entity| U
    U -->|DTO| H
    H -->|JSON| OR
    OR -->|TypeScript Type| TQ
    TQ -->|Cached Data| C
```

## Error Handling Flow

```mermaid
flowchart TD
    A[API Request] --> B{Response Status}

    B -->|2xx| C[Success Handler]
    B -->|401| D[Redirect to Login]
    B -->|403| E[Show Forbidden]
    B -->|404| F[Show Not Found]
    B -->|422| G[Show Validation Errors]
    B -->|5xx| H[Show Server Error]

    C --> I[Update Cache]
    D --> J[Clear Session]
    G --> K[Highlight Form Fields]
    H --> L[Retry Option]
```
