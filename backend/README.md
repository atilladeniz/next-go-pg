# Backend

Go Backend mit Clean Architecture, generiert mit [Goca CLI](https://github.com/sazardev/goca).

## Tech Stack

- **Go 1.23** - Programmiersprache
- **Gorilla Mux** - HTTP Router
- **GORM** - ORM für PostgreSQL
- **Swagger/swag** - API Dokumentation
- **Goca CLI** - Code Generator für Clean Architecture

## Architecture

Clean Architecture mit strikter Layer-Trennung:

```
internal/
├── domain/           # Entities, Business Rules
├── usecase/          # Application Logic
├── repository/       # Data Access Layer
├── handler/          # HTTP Handler
├── middleware/       # Auth, CORS
└── sse/              # Server-Sent Events
```

| Layer | Beschreibung | Goca Befehl |
|-------|--------------|-------------|
| Domain | Entities, Value Objects | `goca make entity` |
| UseCase | Business Logic | `goca make usecase` |
| Repository | Database Operations | `goca make repository` |
| Handler | HTTP Endpoints | `goca make handler` |

## Quick Start

### 1. Install dependencies:
```bash
go mod tidy
```


### 2. Configure database (PostgreSQL):

#### Option A: Using Docker (Recommended)
```bash
# Run PostgreSQL
docker run --name postgres-dev \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=backend \
  -p 5432:5432 \
  -d postgres:15

# Or using docker-compose
docker-compose up -d postgres
```


#### Option B: Local PostgreSQL
```bash
# Create database
createdb backend
```


### 3. Configure environment variables:
```bash
# Copy example file
cp .env.example .env

# Edit with your credentials
# DB_PASSWORD=password
# DB_NAME=backend
```


### 4. Run the application:
```bash
go run cmd/server/main.go
```


### 5. Test endpoints:
```bash
# Health check
curl http://localhost:8080/health

# Create user (if you have the User feature)
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}'
```


## Project Structure

```
backend/
├── cmd/
│   └── server/           # Application entry point
│       └── main.go
├── internal/
│   ├── domain/           # Entities (goca make entity)
│   ├── usecase/          # Business Logic (goca make usecase)
│   ├── repository/       # Data Access (goca make repository)
│   ├── handler/          # HTTP Handler (goca make handler)
│   ├── middleware/       # Auth, CORS
│   └── sse/              # Server-Sent Events
├── pkg/
│   ├── config/           # Application configuration
│   └── logger/           # Logging system
├── docs/                 # Swagger documentation
├── .goca.yaml            # Goca configuration
├── .env                  # Environment variables
├── .env.example          # Configuration example
├── Makefile              # Build commands
└── go.mod
```


## Goca Befehle

### Neues Feature generieren

```bash
# Komplettes Feature mit allen Layers
goca feature Product --fields "name:string,price:float64,stock:int"

# Feature mit Validierung
goca feature Order --fields "userId:string,total:float64" --validation

# Alle Features integrieren (Routes, DI)
goca integrate --all
```

### Einzelne Layer generieren

```bash
# Nur Entity (Domain Layer)
goca make entity Product

# Nur Repository (Data Layer)
goca make repository Product

# Nur UseCase (Business Logic)
goca make usecase Product

# Nur Handler (HTTP Layer)
goca make handler Product
```

### Nach Goca/API-Änderungen

```bash
# Swagger + Orval in einem Befehl (vom Root-Verzeichnis)
cd ..
make api

# Das macht automatisch:
# 1. swag init → backend/docs/swagger.json
# 2. orval → frontend/src/api/endpoints/
```

## Development Commands

```bash
# Run application
make run

# Run tests
make test

# Build for production
make build

# Linting and formatting
make lint
make fmt

# Swagger generieren
make swagger
```


## Troubleshooting

### Error: "dial tcp [::1]:5432: connection refused"
PostgreSQL database is not running. 

**Solution:**
```bash
# With Docker
docker run --name postgres-dev \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=backend \
  -p 5432:5432 \
  -d postgres:15

# Verify it's running
docker ps
```


### Error: "database not configured"
Database environment variables are not configured.

**Solution:**
```bash
# Configure in .env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=backend
```


### Error: "command not found: goca"
Goca CLI is not installed or not in PATH.

**Solution:**
```bash
# Reinstall Goca
go install github.com/sazardev/goca@latest

# Verify installation
goca version
```


### Health Check shows "degraded"
Application runs but cannot connect to database.

**Solution:**
1. Verify PostgreSQL is running
2. Verify environment variables in .env
3. Test connection manually: `psql -h localhost -U postgres -d backend`

## Additional Resources

- [Goca Documentation](https://github.com/sazardev/goca)
- [Clean Architecture Principles](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Complete Tutorial](https://github.com/sazardev/goca/wiki/Complete-Tutorial)

## Contributing

This project was generated with Goca. To contribute:

1. Add new features with `goca feature`
2. Maintain layer separation
3. Write tests for new functionality
4. Follow Clean Architecture conventions

---

Generated with [Goca](https://github.com/sazardev/goca)
