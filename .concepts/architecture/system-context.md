# System Context Diagram (C4 Level 1)

## Business Context

Next-Go-PG ist eine Full-Stack Webanwendung mit:
- **Frontend**: Next.js 16 SPA für User Interface
- **Backend**: Go API für Business Logic
- **Database**: PostgreSQL für Persistenz
- **Auth**: Better Auth für Authentifizierung

## System Context Diagram

```mermaid
C4Context
    title Next-Go-PG - System Context Diagram

    Person(user, "User", "Registrierter Benutzer der Anwendung")
    Person(admin, "Admin", "Administrator mit erweiterten Rechten")

    System(nextgopg, "Next-Go-PG", "Full-Stack Webanwendung mit Next.js Frontend und Go Backend")

    System_Ext(email, "Email Service", "SMTP/SendGrid für Transaktions-Emails")
    System_Ext(github, "GitHub", "OAuth Provider & Container Registry")

    Rel(user, nextgopg, "Nutzt", "HTTPS")
    Rel(admin, nextgopg, "Verwaltet", "HTTPS")
    Rel(nextgopg, email, "Sendet Emails", "SMTP/API")
    Rel(nextgopg, github, "Auth & Registry", "OAuth/HTTPS")

    UpdateLayoutConfig($c4ShapeInRow="3", $c4BoundaryInRow="1")
```

## Akteure

| Akteur | Beschreibung | Interaktionen |
|--------|--------------|---------------|
| **User** | Registrierter Benutzer | Login, CRUD Operationen, Dashboard |
| **Admin** | Administrator | User Management, System Settings |

## Externe Systeme

| System | Zweck | Protokoll |
|--------|-------|-----------|
| **Email Service** | Passwort Reset, Notifications | SMTP / REST API |
| **GitHub** | OAuth Login, Container Registry | OAuth 2.0 / HTTPS |

## Datenflüsse

```mermaid
flowchart LR
    subgraph External
        U[User Browser]
        E[Email Service]
        G[GitHub]
    end

    subgraph Next-Go-PG
        F[Frontend :3000]
        B[Backend :8080]
        D[(PostgreSQL)]
    end

    U -->|HTTPS| F
    F -->|API Calls| B
    B -->|SQL| D
    B -->|SMTP| E
    B -->|OAuth| G
```

## Sicherheitsgrenzen

- **Internet → Frontend**: TLS 1.3, kamal-proxy
- **Frontend → Backend**: Interne Kommunikation, Session Cookies
- **Backend → Database**: Verschlüsselte Verbindung (sslmode=require)
