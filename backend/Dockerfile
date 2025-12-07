# syntax=docker/dockerfile:1

# Multi-arch Docker image for Go backend
# Build stage
FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the server binary
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-s -w" -o bin/backend cmd/server/main.go

# Build the migrate binary
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-s -w" -o bin/migrate cmd/migrate/main.go

# Final stage - Alpine for shell support (needed for entrypoint)
FROM alpine:3.20

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binaries from builder stage
COPY --from=builder /app/bin/backend .
COPY --from=builder /app/bin/migrate .

# Copy migrations
COPY --from=builder /app/migrations ./migrations

# Copy entrypoint script
COPY scripts/docker-entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Create non-root user
RUN adduser -D -u 1000 appuser
USER appuser

# Expose port
EXPOSE 8080

# Run with entrypoint (migrations + server)
ENTRYPOINT ["/entrypoint.sh"]
