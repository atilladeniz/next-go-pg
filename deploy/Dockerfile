# Multi-stage Dockerfile for GocaTest
# Builds both Frontend (Next.js) and Backend (Go) in a single container
# Supervised by a lightweight process manager
#
# Build: docker build -f deploy/Dockerfile -t gocatest .
# Run:   docker run -p 80:80 gocatest

# ============================================
# Stage 1: Build Frontend (Next.js)
# ============================================
FROM node:22-alpine AS frontend-builder

WORKDIR /app/frontend

# Install bun
RUN npm install -g bun

# Copy package files
COPY frontend/package.json frontend/bun.lock* ./

# Install dependencies
RUN bun install --frozen-lockfile

# Copy frontend source
COPY frontend/ ./

# Build arguments for environment
ARG NEXT_PUBLIC_APP_URL
ARG NEXT_PUBLIC_API_URL

ENV NEXT_PUBLIC_APP_URL=$NEXT_PUBLIC_APP_URL
ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL

# Build Next.js (standalone output)
RUN bun run build

# ============================================
# Stage 2: Build Backend (Go)
# ============================================
FROM golang:1.23-alpine AS backend-builder

WORKDIR /app/backend

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY backend/go.mod backend/go.sum ./

# Download dependencies
RUN go mod download

# Copy backend source
COPY backend/ ./

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server ./cmd/server

# ============================================
# Stage 3: Production Runtime
# ============================================
FROM alpine:3.20 AS runtime

# Install runtime dependencies
RUN apk add --no-cache \
    nodejs \
    npm \
    ca-certificates \
    tzdata \
    curl \
    supervisor

WORKDIR /app

# Copy frontend build (standalone)
COPY --from=frontend-builder /app/frontend/.next/standalone ./frontend/
COPY --from=frontend-builder /app/frontend/.next/static ./frontend/.next/static
COPY --from=frontend-builder /app/frontend/public ./frontend/public

# Copy backend binary
COPY --from=backend-builder /server ./backend/server

# Copy supervisor config
COPY deploy/supervisord.conf /etc/supervisor/conf.d/supervisord.conf

# Create non-root user
RUN addgroup -g 1001 -S app && \
    adduser -S -D -H -u 1001 -s /sbin/nologin -G app app

# Set ownership
RUN chown -R app:app /app

# Health check endpoint
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Expose port (kamal-proxy expects port 80)
EXPOSE 80

# Run supervisor (manages both frontend and backend)
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]
