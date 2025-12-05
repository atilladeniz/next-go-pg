# Logging - Next-Go-PG

## Overview

Structured JSON logging for production observability, designed for:
- **Log Aggregation**: ELK, Datadog, CloudWatch, Loki
- **Error Tracking**: Sentry, PostHog integration ready
- **Distributed Tracing**: Request ID, Trace ID, Span ID support
- **Audit Trails**: Business event logging for compliance

## Quick Reference

### Backend (Go)

```go
import "github.com/atilladeniz/next-go-pg/backend/pkg/logger"

// Simple
logger.Info().Msg("Server started")
logger.Error().Err(err).Msg("Failed")

// Structured
logger.Info().
    Str("user_id", "123").
    Str("action", "login").
    Msg("User logged in")

// With request context
logger.WithContext(ctx).Info().Msg("Request processed")
```

### Frontend (Next.js)

```typescript
import { log, createLogger, setUserContext } from "@shared/lib/logger"

// Simple
log.info("Page loaded")
log.error("API failed", { endpoint: "/api/users" })

// With user context
setUserContext({ userId: "123", userName: "Max" })
log.info("Action performed") // Includes user info

// Component-specific
const authLogger = createLogger("auth")
authLogger.info("Login submitted")
```

---

## Backend (Go) - zerolog

### Location
- Logger: `backend/pkg/logger/logger.go`
- Middleware: `backend/internal/middleware/logging.go`

### Features

| Feature | Description |
|---------|-------------|
| Categories | `http`, `auth`, `db`, `business`, `system`, `sse`, `cache` |
| User Context | Automatic user_id, user_name, session_id in logs |
| Request Tracing | request_id, trace_id, span_id |
| Caller Info | file:line in production logs |
| Slow Query Detection | Auto-warn for queries >100ms |
| Privacy | Email masking, sensitive data protection |

### Configuration

```bash
LOG_LEVEL=info          # debug, info, warn, error
ENVIRONMENT=production  # development = pretty colored output
```

### Log Categories

Filter logs by category in your log aggregation tool:

```go
// HTTP requests
logger.HTTPRequest(ctx, logger.HTTPRequestLog{...})

// Database queries
logger.DBQuery(ctx, "SELECT", "users", latency, rows, nil)

// Authentication
logger.Auth(ctx, logger.AuthEvent{
    Action:  "login",
    UserID:  "123",
    Email:   "user@example.com",
    Success: true,
})

// Business events (audit trail)
logger.Business(ctx, logger.BusinessEventLog{
    EventType: "order.placed",
    Entity:    "order",
    EntityID:  "ord-123",
    Action:    "create",
    Changes:   map[string]any{"total": 99.99},
})

// System events
logger.SystemEvent("startup", map[string]any{"version": "1.0.0"})
```

### User Context in Logs

Add user info to context after authentication:

```go
// In auth middleware after successful auth
ctx = logger.AddUserToContext(ctx, user.ID, user.Name, user.Email, session.ID)

// All subsequent logs include user info automatically
logger.WithContext(ctx).Info().Msg("Action performed")
// Output: {"user_id":"123","user_name":"Max","session_id":"sess-456",...}
```

### HTTP Request Logging

The middleware automatically logs:

```json
{
  "level": "info",
  "category": "http",
  "request_id": "abc-123",
  "method": "POST",
  "path": "/api/v1/orders",
  "status": 201,
  "latency": "15.2ms",
  "client_ip": "192.168.1.1",
  "user_agent": "Mozilla/5.0...",
  "request_size": 256,
  "response_size": 1024,
  "user_id": "123",
  "user_name": "Max",
  "message": "HTTP request"
}
```

### Performance Logging

```go
// Timed operations
defer logger.Timed(ctx, "process_order")()

// Warn only if slow
defer logger.TimedWithThreshold(ctx, "db_query", 100*time.Millisecond)()
```

### Distributed Tracing

Support for trace propagation:

```go
// Headers are automatically extracted by middleware:
// X-Request-ID, X-Trace-ID, X-Span-ID

// Or add manually:
ctx = logger.AddTraceToContext(ctx, traceID, spanID)
```

---

## Frontend (Next.js) - Pino

### Location
`frontend/src/shared/lib/logger/index.ts`

### Features

| Feature | Description |
|---------|-------------|
| Categories | `http`, `auth`, `api`, `ui`, `business`, `performance`, `error` |
| User Context | Global user context for all logs |
| Server + Client | Works in SSR and browser |
| Redaction | Auto-redact passwords, tokens, credit cards |
| Performance | Timing helpers, slow operation detection |

### Configuration

```bash
LOG_LEVEL=info      # debug, info, warn, error
NODE_ENV=production # development = pretty output
```

### Log Categories

```typescript
import { log, LogCategory } from "@shared/lib/logger"

// HTTP/API
log.request("GET", "/api/stats", 200, 42)
log.apiCall("/api/users", "POST", 150, true)

// Authentication
log.authSuccess("login", userId)
log.authFailure("login", "invalid_password")

// Business events
log.event("checkout_completed", { orderId: "123", total: 99.99 })

// UI/Component events
log.component("CheckoutForm", "submit", { items: 3 })

// Performance
log.navigation("/dashboard", 250)
log.slowOperation("renderList", 500, 200) // Warns if >200ms

// Errors
log.exception(error, { context: "payment" })
log.unhandled(error, "window.onerror")
```

### User Context

```typescript
import { setUserContext, clearUserContext } from "@shared/lib/logger"

// After login
setUserContext({
  userId: session.user.id,
  userName: session.user.name,
  sessionId: session.id,
})

// After logout
clearUserContext()

// All logs now include user info
log.info("Page viewed") // {"userId":"123","userName":"Max",...}
```

### Component Loggers

```typescript
import { createLogger } from "@shared/lib/logger"

const logger = createLogger("PaymentForm")

function PaymentForm() {
  const done = logger.timed("processPayment")

  try {
    // ... process
    logger.info("Payment successful", { amount: 99.99 })
  } catch (err) {
    logger.error("Payment failed", { error: err.message })
  } finally {
    done() // Logs duration
  }
}
```

### Async Timing Helper

```typescript
import { withTiming } from "@shared/lib/logger"

const result = await withTiming(
  "fetchUserData",
  () => fetchUser(id),
  200 // Warn if >200ms
)
```

---

## Output Examples

### Development (Pretty)

```
15:04:05 | INFO  | http | HTTP request method=GET path=/api/stats status=200 latency=15ms user_id=123
15:04:06 | WARN  | auth | Auth event action=login success=false reason=invalid_password email=us***@example.com
15:04:07 | INFO  | business | Business event event_type=order.placed entity=order entity_id=ord-123
```

### Production (JSON)

```json
{"level":"info","time":"2024-01-15T10:30:45Z","service":"next-go-pg-api","version":"1.0.0","category":"http","request_id":"abc-123","user_id":"123","user_name":"Max","method":"GET","path":"/api/stats","status":200,"latency":15000000,"message":"HTTP request"}
```

---

## Self-Hosted Log Aggregation (Grafana + Loki)

This project includes a pre-configured Grafana + Loki + Promtail stack for local/self-hosted log aggregation.

### Quick Start

```bash
# Start logging stack
make logs-up

# Open Grafana
make logs-open  # Opens http://localhost:3001

# Stop logging stack
make logs-down
```

### Components

| Service | Port | Description |
|---------|------|-------------|
| Grafana | 3001 | Dashboard & visualization |
| Loki | 3100 | Log aggregation backend |
| Promtail | - | Log shipper (collects container logs) |

### Default Credentials

- **Username**: admin
- **Password**: admin

### Querying Logs

#### In Grafana

1. Open http://localhost:3001
2. Go to Explore (compass icon)
3. Select "Loki" datasource
4. Use LogQL queries:

```logql
# All logs from backend
{service="next-go-pg-api"}

# Error logs only
{level="error"}

# Auth events
{category="auth"}

# HTTP requests with status 500+
{category="http"} |= "status\":5"

# Logs from specific user
{service="next-go-pg-api"} |= "user_id\":\"123\""

# Combined filters
{service="next-go-pg-api", level="error"} |= "database"
```

#### Via CLI

```bash
# Query logs from command line
make logs-query q='{service="next-go-pg-api"}'
make logs-query q='{level="error"}' limit=50
make logs-query q='{category="auth"}'
```

### Configuration Files

- `deploy/loki/loki-config.yml` - Loki configuration
- `deploy/loki/promtail-config.yml` - Promtail log collection rules
- `deploy/grafana/provisioning/datasources/datasources.yml` - Grafana datasource
- `docker-compose.logging.yml` - Docker Compose for logging stack

### Log Labels (for filtering)

Promtail extracts these labels from JSON logs:

| Label | Description |
|-------|-------------|
| `service` | Application name (next-go-pg-api, next-go-pg-frontend) |
| `level` | Log level (debug, info, warn, error) |
| `category` | Log category (http, auth, db, business, etc.) |

### Retention

Default: **31 days** (configurable in `loki-config.yml`)

### Production Deployment

For production with Kamal:

1. Deploy Loki + Grafana on your server
2. Update Promtail config to ship logs to Loki
3. Configure retention and storage based on volume

---

## Log Aggregation

### Filtering by Category

```
# Kibana/ELK
category: "auth" AND success: false

# Datadog
@category:db @slow_query:true

# CloudWatch Insights
fields @timestamp, @message
| filter category = 'business'
| filter event_type like /order/
```

### Alerts

Set up alerts for:
- `category:auth AND success:false` (count > 10/min)
- `category:db AND slow_query:true`
- `level:error`
- `category:http AND status:>=500`

---

## Integration

### Sentry

```go
// Backend
import "github.com/getsentry/sentry-go"

logger.Error().Err(err).Msg("Critical error")
sentry.CaptureException(err)
```

```typescript
// Frontend
import * as Sentry from "@sentry/nextjs"

log.exception(error, { context: "payment" })
Sentry.captureException(error)
```

### PostHog

```typescript
import posthog from "posthog-js"

// Log + Track together
log.event("checkout_completed", { orderId, total })
posthog.capture("checkout_completed", { orderId, total })
```

---

## Best Practices

1. **Use structured logging** - Add context as fields, not in message
   ```go
   // ✅ Good
   logger.Info().Str("user_id", id).Msg("User created")

   // ❌ Bad
   logger.Info().Msgf("User %s created", id)
   ```

2. **Use categories** - Makes filtering easy
   ```go
   logger.WithCategoryCtx(ctx, logger.CategoryAuth).Info()...
   ```

3. **Include request context** - Always use `WithContext(ctx)`

4. **Log business events** - For audit trails and analytics

5. **Don't log sensitive data** - Use redaction, mask emails

6. **Use appropriate levels**
   - `debug`: Development details
   - `info`: Normal operations
   - `warn`: Recoverable issues
   - `error`: Failures
   - `fatal`: Application cannot continue
