# Logging - Next-Go-PG

## Overview

This project uses structured JSON logging for both backend and frontend, designed for production observability and easy integration with log aggregation services (ELK, Datadog, CloudWatch, Sentry, PostHog).

## Backend (Go) - zerolog

### Location
- Logger: `backend/pkg/logger/logger.go`
- Middleware: `backend/internal/middleware/logging.go`

### Features
- Structured JSON output (production)
- Pretty colored output (development)
- Request ID tracing across all logs
- Log levels: debug, info, warn, error, fatal
- Email masking for privacy
- Business event logging for audit trails

### Configuration

Environment variables:
```bash
LOG_LEVEL=info      # debug, info, warn, error
ENVIRONMENT=production  # development = pretty output
```

### Usage

```go
import "github.com/atilladeniz/next-go-pg/backend/pkg/logger"

// Simple logging
logger.Info().Msg("Server started")
logger.Error().Err(err).Msg("Database connection failed")

// Structured logging
logger.Info().
    Str("user_id", "123").
    Str("action", "login").
    Int("duration_ms", 42).
    Msg("User logged in")

// With request context (includes request_id automatically)
logger.WithContext(ctx).Info().
    Str("endpoint", "/api/users").
    Msg("Request processed")

// Business events (for audit trails)
logger.BusinessEvent(ctx, "user.created", "user", userID, map[string]interface{}{
    "email": user.Email,
    "plan": "premium",
})

// Auth events
logger.AuthSuccess(ctx, userID, email)
logger.AuthFailure(ctx, email, "invalid_password")
```

### Output Examples

Development (pretty):
```
2024-01-15T10:30:45+01:00 INF Server starting port=8080 service=next-go-pg-api version=dev
2024-01-15T10:30:46+01:00 INF HTTP request client_ip=127.0.0.1 latency=1.234ms method=GET path=/api/v1/stats request_id=abc-123 status=200
```

Production (JSON):
```json
{"level":"info","service":"next-go-pg-api","version":"1.0.0","time":"2024-01-15T10:30:45Z","request_id":"abc-123","method":"GET","path":"/api/v1/stats","status":200,"latency":1234000,"client_ip":"127.0.0.1","message":"HTTP request"}
```

### Request ID Tracing

Every HTTP request gets a unique ID (`X-Request-ID` header):
1. If client sends `X-Request-ID`, it's used
2. Otherwise, a new UUID is generated
3. ID is returned in response headers
4. ID is automatically added to all logs via `logger.WithContext(ctx)`

## Frontend (Next.js) - Pino

### Location
`frontend/src/shared/lib/logger/index.ts`

### Features
- Structured JSON output (production)
- Pretty colored output (development)
- Works on server and client (browser)
- Sensitive data redaction
- Component-specific child loggers

### Usage

```typescript
import { log, createLogger } from "@shared/lib/logger"

// Simple logging
log.info("Page loaded")
log.error("API call failed", { endpoint: "/api/users", status: 500 })

// HTTP request logging
log.request("GET", "/api/stats", 200, 42)

// Auth events
log.authSuccess(userId)
log.authFailure("invalid_credentials")

// API call tracking
log.apiCall("/api/users", "GET", 150, true)

// Business events (for PostHog integration)
log.event("button_clicked", { buttonId: "submit", page: "checkout" })

// Component-specific logger
const authLogger = createLogger("auth")
authLogger.info("Login form submitted")
authLogger.error("Validation failed", { field: "email" })
```

### Configuration

Environment variables:
```bash
LOG_LEVEL=info      # debug, info, warn, error
NODE_ENV=production # development = pretty output
```

### Sensitive Data Redaction

These fields are automatically redacted:
- `password`
- `token`
- `authorization`
- `cookie`
- `*.password`
- `*.token`

## Integration with External Services

### Sentry Integration

Backend (Go):
```go
// In error handler or middleware
import "github.com/getsentry/sentry-go"

if err != nil {
    logger.Error().Err(err).Msg("Critical error")
    sentry.CaptureException(err)
}
```

Frontend (Next.js):
```typescript
// In error boundary or catch
import * as Sentry from "@sentry/nextjs"

log.error("Critical error", { error: err.message })
Sentry.captureException(err)
```

### PostHog Integration

```typescript
import { log } from "@shared/lib/logger"
import posthog from "posthog-js"

// Log + track
log.event("checkout_completed", { orderId, total })
posthog.capture("checkout_completed", { orderId, total })
```

### Log Aggregation (ELK/Datadog)

With Kamal deployment, logs are collected from stdout. Configure your log shipper to parse JSON:

```yaml
# filebeat.yml example
- type: container
  paths:
    - '/var/lib/docker/containers/*/*.log'
  json.keys_under_root: true
  json.add_error_key: true
```

## Best Practices

1. **Always use structured logging** - Add context as fields, not in the message
   ```go
   // Good
   logger.Info().Str("user_id", id).Msg("User created")

   // Bad
   logger.Info().Msgf("User %s created", id)
   ```

2. **Use appropriate log levels**
   - `debug`: Development details
   - `info`: Normal operations
   - `warn`: Recoverable issues
   - `error`: Failures requiring attention
   - `fatal`: Application cannot continue

3. **Include request IDs** - Use `logger.WithContext(ctx)` for request tracing

4. **Log business events** - Use `BusinessEvent()` for audit trails

5. **Don't log sensitive data** - Use redaction, mask emails
