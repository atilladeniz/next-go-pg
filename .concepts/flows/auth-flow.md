# Authentication Flow

## Business Logic

Die Authentifizierung basiert auf Better Auth mit:
- Email/Passwort Login
- Session-basierte Auth (Cookies)
- Cross-Tab Synchronisation
- Server-side Session Validation

## Login Flow

```mermaid
sequenceDiagram
    participant U as User
    participant F as Frontend
    participant BA as Better Auth
    participant DB as Database

    U->>F: Enter credentials
    F->>BA: POST /api/auth/sign-in/email
    BA->>DB: SELECT user WHERE email
    DB-->>BA: User record

    alt Password valid
        BA->>DB: INSERT session
        BA-->>F: Set-Cookie: session
        F-->>U: Redirect to /dashboard
    else Password invalid
        BA-->>F: 401 Unauthorized
        F-->>U: Show error
    end
```

## Registration Flow

```mermaid
sequenceDiagram
    participant U as User
    participant F as Frontend
    participant BA as Better Auth
    participant DB as Database
    participant E as Email Service

    U->>F: Fill registration form
    F->>BA: POST /api/auth/sign-up/email

    BA->>DB: Check email exists
    DB-->>BA: Not found

    BA->>DB: INSERT user
    BA->>DB: INSERT session
    BA->>E: Send welcome email

    BA-->>F: Set-Cookie: session
    F-->>U: Redirect to /dashboard
```

## Session Validation (Server Component)

```mermaid
flowchart TD
    A[Request to Protected Page] --> B{Has Session Cookie?}
    B -->|No| C[Redirect to /login]
    B -->|Yes| D[getSession on Server]
    D --> E{Session Valid?}
    E -->|No| C
    E -->|Yes| F[Render Page with User Data]
    F --> G[prefetchQuery with Cookie Header]
    G --> H[HydrationBoundary]
    H --> I[Client Component receives hydrated data]
```

## Cross-Tab Logout

```mermaid
sequenceDiagram
    participant T1 as Tab 1
    participant T2 as Tab 2
    participant BC as BroadcastChannel
    participant BA as Better Auth

    T1->>BA: signOut()
    BA-->>T1: Session cleared
    T1->>BC: broadcast('auth-logout')

    BC-->>T2: Receive 'auth-logout'
    T2->>T2: Clear local state
    T2->>T2: router.push('/')

    Note over T1,T2: Beide Tabs sind ausgeloggt
```

## State Diagram: User Session

```mermaid
stateDiagram-v2
    [*] --> Anonymous: Visit Site

    Anonymous --> Authenticating: Click Login
    Authenticating --> Authenticated: Success
    Authenticating --> Anonymous: Failed

    Anonymous --> Registering: Click Register
    Registering --> Authenticated: Success
    Registering --> Anonymous: Failed

    Authenticated --> Anonymous: Logout
    Authenticated --> Anonymous: Session Expired
    Authenticated --> Anonymous: Cross-Tab Logout

    note right of Authenticated
        Session Cookie set
        User data available
        Protected routes accessible
    end note
```

## Token/Cookie Flow

```mermaid
flowchart LR
    subgraph Client
        B[Browser]
        C[Cookie Storage]
    end

    subgraph Server
        BA[Better Auth]
        DB[(Sessions Table)]
    end

    B -->|1. Login Request| BA
    BA -->|2. Create Session| DB
    BA -->|3. Set-Cookie| C
    B -->|4. Subsequent Requests + Cookie| BA
    BA -->|5. Validate Session| DB
    DB -->|6. Session Data| BA
    BA -->|7. Response + User Data| B
```

## Security Considerations

| Aspekt | Implementation |
|--------|----------------|
| **Cookie Flags** | HttpOnly, Secure, SameSite=Lax |
| **Session Storage** | Database (PostgreSQL) |
| **Password Hashing** | bcrypt (Better Auth default) |
| **CSRF Protection** | SameSite Cookie + Origin Check |
| **Rate Limiting** | Middleware (zu implementieren) |
