# syntax=docker/dockerfile:1

# Multi-arch Docker image for Go backend
# Build stage
FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

# Install dependencies (mounted module cache survives layer invalidation)
COPY go.mod go.sum ./
RUN --mount=type=cache,id=gomod,target=/go/pkg/mod \
    go mod download

# Copy source code
COPY . .

# Build with persistent module + build caches so incremental builds
# only recompile changed packages.
RUN --mount=type=cache,id=gomod,target=/go/pkg/mod \
    --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-s -w" -o bin/backend cmd/server/main.go && \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-s -w" -o bin/migrate cmd/migrate/main.go

# Final stage
FROM alpine:3.20

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binaries from builder stage
COPY --from=builder /app/bin/backend .
COPY --from=builder /app/bin/migrate .

# Copy migrations (callable on-demand via `kamal app exec /app/migrate -up`)
COPY --from=builder /app/migrations ./migrations

# Create non-root user
RUN adduser -D -u 1000 appuser
USER appuser

# Expose port
EXPOSE 8080

# Server runs AutoMigrate on startup; SQL migrations are invoked
# explicitly (Kamal pre-deploy hook or `just prod-migrate-up`), not here.
ENTRYPOINT ["/app/backend"]
