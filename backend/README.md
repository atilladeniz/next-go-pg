# Backend

Generated with Goca - Go Clean Architecture Code Generator

## Architecture

This project follows Clean Architecture principles:

- **Domain**: Entities and business rules  
- **Use Cases**: Application logic
- **Repository**: Data abstraction
- **Handler**: Delivery adapters

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
│   ├── domain/           # Entities and business rules
│   ├── usecase/          # Application logic
│   ├── repository/       # Persistence implementations
│   ├── handler/          # HTTP/gRPC adapters
│   │   └── http/
│   └── messages/         # Error and response messages
├── pkg/
│   ├── config/           # Application configuration
│   └── logger/           # Logging system
├── migrations/           # Database migrations
├── .env                  # Environment variables
├── .env.example          # Configuration example
├── docker-compose.yml    # Docker services
├── Makefile              # Useful commands
├── go.mod
└── README.md
```


## Useful Commands

### Generate new features:
```bash
# Complete feature with all layers
goca feature User --fields "name:string,email:string"

# Feature with validations
goca feature Product --fields "name:string,price:float64" --validation

# Integrate existing features
goca integrate --all
```


### Development commands:
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
