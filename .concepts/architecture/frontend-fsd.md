# Frontend Architecture - Feature-Sliced Design (C4 Level 3)

## Architecture Overview

The frontend uses **Feature-Sliced Design (FSD)** - an architectural methodology for scalable frontend applications.

## FSD Layer Diagram

```mermaid
flowchart TB
    subgraph App["app/ - Next.js App Router"]
        direction LR
        Auth["(auth)/"]
        Protected["(protected)/"]
        API["api/auth/"]
    end

    subgraph Widgets["widgets/ - Composite UI"]
        Header["header/"]
    end

    subgraph Features["features/ - User Interactions"]
        AuthFeature["auth/"]
        Stats["stats/"]
    end

    subgraph Entities["entities/ - Business Objects"]
        User["user/"]
    end

    subgraph Shared["shared/ - Reusable Code"]
        UI["ui/"]
        SharedAPI["api/"]
        Lib["lib/"]
        Config["config/"]
    end

    App --> Widgets
    App --> Features
    App --> Entities
    App --> Shared

    Widgets --> Features
    Widgets --> Entities
    Widgets --> Shared

    Features --> Entities
    Features --> Shared

    Entities --> Shared

    style App fill:#e1f5fe
    style Widgets fill:#fff3e0
    style Features fill:#e8f5e9
    style Entities fill:#fce4ec
    style Shared fill:#f3e5f5
```

## Import Rules (Dependency Flow)

```mermaid
flowchart TD
    subgraph "Layer Hierarchy (only downward!)"
        A[app/] -->|"can import"| W[widgets/]
        A --> F[features/]
        A --> E[entities/]
        A --> S[shared/]

        W --> F
        W --> E
        W --> S

        F --> E
        F --> S

        E --> S
    end

    subgraph "FORBIDDEN"
        X1[shared/ → features/]
        X2[entities/ → features/]
        X3[features/ → features/]
    end

    style X1 fill:#ffcdd2,stroke:#c62828
    style X2 fill:#ffcdd2,stroke:#c62828
    style X3 fill:#ffcdd2,stroke:#c62828
```

## Slice Structure (Component Diagram)

```mermaid
flowchart TB
    subgraph "features/auth/"
        direction TB
        subgraph UI1["ui/"]
            LF[login-form.tsx]
            RF[register-form.tsx]
        end
        subgraph Model1["model/"]
            AS[use-auth-sync.ts]
        end
        Index1[index.ts]

        UI1 --> Index1
        Model1 --> Index1
    end

    subgraph "features/stats/"
        direction TB
        subgraph UI2["ui/"]
            SG[stats-grid.tsx]
        end
        subgraph Model2["model/"]
            SSE[use-sse.ts]
        end
        Index2[index.ts]

        UI2 --> Index2
        Model2 --> Index2
    end

    subgraph "entities/user/"
        direction TB
        subgraph UI3["ui/"]
            UI4[user-info.tsx]
        end
        subgraph Model3["model/"]
            Types[types.ts]
        end
        Index3[index.ts]

        UI3 --> Index3
        Model3 --> Index3
    end
```

## Shared Layer Details

```mermaid
flowchart TB
    subgraph Shared["shared/"]
        subgraph UI["ui/ (shadcn)"]
            Button[button.tsx]
            Card[card.tsx]
            Input[input.tsx]
            Dialog[dialog.tsx]
            More[...60+ components]
        end

        subgraph API["api/ (Orval)"]
            Endpoints[endpoints/]
            Models[models/]
            Fetch[custom-fetch.ts]
        end

        subgraph Lib["lib/"]
            AuthClient["auth-client/"]
            AuthServer["auth-server/"]
            QueryClient[query-client.ts]
            Utils[utils.ts]
        end

        subgraph Config["config/"]
            Providers[providers.tsx]
            Theme[theme-provider.tsx]
        end
    end

    Endpoints -->|"generated from"| Swagger[(Swagger JSON)]
    AuthClient -->|"Client-safe"| BetterAuth[Better Auth]
    AuthServer -->|"Server-only"| BetterAuth
```

## Data Flow Diagram

```mermaid
sequenceDiagram
    participant App as app/(protected)/
    participant Widget as widgets/header
    participant Feature as features/stats
    participant Entity as entities/user
    participant Shared as shared/api

    App->>Shared: getSession()
    Shared-->>App: SessionUser
    App->>App: prefetchQuery()
    App->>Widget: <Header user={session.user} />
    Widget->>Feature: useAuthSync()
    App->>Feature: <StatsGrid />
    Feature->>Shared: useGetStats()
    Shared-->>Feature: Stats Data
    Feature->>Entity: <UserInfo user={user} />
```

## Request/Response Flow

```mermaid
flowchart LR
    subgraph "Server Component"
        SC[page.tsx]
        Session[getSession]
        Prefetch[prefetchQuery]
        Hydration[HydrationBoundary]
    end

    subgraph "Client Component"
        CC[stats-grid.tsx]
        Query[useGetStats]
        SSE[useSSE]
    end

    subgraph "Shared Layer"
        API[Orval Hooks]
        Auth[auth-server]
        QC[QueryClient]
    end

    SC --> Session
    Session --> Auth
    SC --> Prefetch
    Prefetch --> API
    Prefetch --> QC
    SC --> Hydration
    Hydration --> CC
    CC --> Query
    Query --> QC
    CC --> SSE
```

## TypeScript Path Aliases

```mermaid
flowchart LR
    subgraph "tsconfig.json paths"
        A1["@shared/*"] --> S1["src/shared/*"]
        A2["@entities/*"] --> S2["src/entities/*"]
        A3["@features/*"] --> S3["src/features/*"]
        A4["@widgets/*"] --> S4["src/widgets/*"]
    end
```

## Linting Architecture

```mermaid
flowchart TB
    subgraph "Quality Gates"
        Biome[Biome Check]
        Steiger[Steiger FSD]
        TypeCheck[TypeScript]
    end

    subgraph "Steiger Rules"
        R1[no-public-api-sidestep]
        R2[forbidden-imports]
        R3[insignificant-slice]
        R4[no-layer-public-api]
    end

    Lint[bun run lint] --> Biome
    Lint --> Steiger
    Steiger --> R1
    Steiger --> R2
    Steiger --> R3
    Steiger --> R4
```

## Comparison: Before vs. After FSD

### Before (Flat Structure)

```
src/
├── api/                # Mixed concerns
├── components/
│   ├── ui/            # shadcn
│   ├── header.tsx     # App-specific
│   ├── providers.tsx  # Config
│   └── ...
├── hooks/             # All hooks mixed
└── lib/               # All utilities mixed
```

### After (FSD)

```
src/
├── widgets/header/     # Composite, uses features
├── features/
│   ├── auth/          # Login, Register, Sync
│   └── stats/         # Stats, SSE
├── entities/user/      # User types, display
└── shared/
    ├── ui/            # shadcn (no business logic)
    ├── api/           # Orval (generated)
    ├── lib/           # Utilities
    └── config/        # App config
```

## Benefits of FSD Architecture

| Aspect | Description |
|--------|-------------|
| **Scalability** | New features = new slices, no changes to existing code |
| **Isolation** | Features are independent, no cross-imports |
| **Discoverability** | Clear structure, easy onboarding |
| **Testability** | Slices can be tested in isolation |
| **Refactoring** | Change one slice without affecting others |
| **Linting** | Automatic architecture rule validation |
