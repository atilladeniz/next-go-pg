# Entity Relationship Diagram

## Database Schema

The ER diagram shows all database entities and their relationships.

## Core Entities

```mermaid
erDiagram
    users {
        uuid id PK "Primary Key"
        string email UK "Unique, Indexed"
        string name "Display Name"
        string email_verified "Verification Status"
        string image "Avatar URL"
        timestamp created_at "Auto"
        timestamp updated_at "Auto"
    }

    sessions {
        uuid id PK
        uuid user_id FK "References users.id"
        string token UK "Session Token"
        timestamp expires_at "Expiration"
        string ip_address "Client IP"
        string user_agent "Browser Info"
        timestamp created_at
        timestamp updated_at
    }

    accounts {
        uuid id PK
        uuid user_id FK "References users.id"
        string account_id "Provider Account ID"
        string provider_id "OAuth Provider"
        string access_token "OAuth Token"
        string refresh_token "Refresh Token"
        timestamp access_token_expires_at
        string scope "OAuth Scope"
        timestamp created_at
        timestamp updated_at
    }

    verifications {
        uuid id PK
        string identifier "Email/Phone"
        string value "Verification Code"
        timestamp expires_at
        timestamp created_at
        timestamp updated_at
    }

    users ||--o{ sessions : "has"
    users ||--o{ accounts : "has"
```

## Application Entities

```mermaid
erDiagram
    users ||--o| user_stats : "has"

    user_stats {
        uuid id PK
        uuid user_id FK UK "One per User"
        int projects "Project Count"
        int tasks "Task Count"
        int completed "Completed Tasks"
        timestamp created_at
        timestamp updated_at
    }

    %% Example: If you have more entities
    %% users ||--o{ projects : "owns"
    %% projects ||--o{ tasks : "contains"

    %% projects {
    %%     uuid id PK
    %%     uuid user_id FK
    %%     string name
    %%     string description
    %%     string status
    %%     timestamp created_at
    %%     timestamp updated_at
    %% }

    %% tasks {
    %%     uuid id PK
    %%     uuid project_id FK
    %%     string title
    %%     string status
    %%     timestamp due_date
    %%     timestamp created_at
    %%     timestamp updated_at
    %% }
```

## Better Auth Schema

Better Auth automatically creates these tables:

```mermaid
erDiagram
    user {
        text id PK
        text name
        text email UK
        text emailVerified
        text image
        timestamp createdAt
        timestamp updatedAt
    }

    session {
        text id PK
        text userId FK
        text token UK
        timestamp expiresAt
        text ipAddress
        text userAgent
        timestamp createdAt
        timestamp updatedAt
    }

    account {
        text id PK
        text userId FK
        text accountId
        text providerId
        text accessToken
        text refreshToken
        timestamp accessTokenExpiresAt
        text scope
        timestamp createdAt
        timestamp updatedAt
    }

    verification {
        text id PK
        text identifier
        text value
        timestamp expiresAt
        timestamp createdAt
        timestamp updatedAt
    }

    user ||--o{ session : "has"
    user ||--o{ account : "has"
```

## Indexes

```mermaid
flowchart LR
    subgraph Primary Keys
        U_PK[users.id]
        S_PK[sessions.id]
        US_PK[user_stats.id]
    end

    subgraph Unique Indexes
        U_EMAIL[users.email]
        S_TOKEN[sessions.token]
        US_USER[user_stats.user_id]
    end

    subgraph Foreign Keys
        S_FK[sessions.user_id → users.id]
        US_FK[user_stats.user_id → users.id]
    end

    subgraph Query Indexes
        S_EXP[sessions.expires_at]
    end
```

## GORM Model Mapping

```go
// Domain Entity → GORM Model → Database Table

type User struct {
    ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
    Email     string    `gorm:"uniqueIndex;not null"`
    Name      string
    CreatedAt time.Time
    UpdatedAt time.Time

    // Relations
    Sessions  []Session
    UserStats *UserStats
}

type UserStats struct {
    ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
    UserID    uuid.UUID `gorm:"uniqueIndex;not null"`
    Projects  int       `gorm:"default:0"`
    Tasks     int       `gorm:"default:0"`
    Completed int       `gorm:"default:0"`
    CreatedAt time.Time
    UpdatedAt time.Time

    // Relations
    User User `gorm:"foreignKey:UserID"`
}
```

## Migration Strategy

```mermaid
flowchart TD
    A[New Entity] --> B[goca make entity]
    B --> C[Add to registry.go]
    C --> D[Restart Backend]
    D --> E[GORM AutoMigrate]
    E --> F{Table exists?}
    F -->|No| G[CREATE TABLE]
    F -->|Yes| H[ALTER TABLE if needed]
    G --> I[Done]
    H --> I
```

## Relationship Types

| Relationship | Example | GORM |
|--------------|---------|------|
| **1:1** | User ↔ UserStats | `hasOne` / `belongsTo` |
| **1:N** | User ↔ Sessions | `hasMany` / `belongsTo` |
| **N:M** | User ↔ Roles | `many2many` (Join Table) |
