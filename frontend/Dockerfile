# syntax=docker/dockerfile:1

# Multi-arch Docker image for Next.js frontend
# Base stage with Bun
FROM --platform=$BUILDPLATFORM oven/bun:1-debian AS base
WORKDIR /app

# Dependencies stage
FROM base AS deps
COPY package.json bun.lock ./
# Cache Bun's package store across builds so repeated installs only
# re-resolve, not re-download.
RUN --mount=type=cache,id=bun-store,target=/root/.bun/install/cache \
    bun install --frozen-lockfile

# Builder stage
FROM base AS builder
ARG SKIP_ENV_VALIDATION=true
ENV SKIP_ENV_VALIDATION=${SKIP_ENV_VALIDATION}
ENV NEXT_TELEMETRY_DISABLED=1

COPY --from=deps /app/node_modules ./node_modules
COPY . .

# Persist Next.js's incremental build cache so unchanged routes/pages
# don't re-transpile on every image build.
RUN --mount=type=cache,id=next-build,target=/app/.next/cache \
    bun run build

# Production stage - use distroless for security
FROM gcr.io/distroless/nodejs22-debian12:nonroot AS runner
WORKDIR /app

ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1
ENV PORT=3000
ENV HOSTNAME="0.0.0.0"

# Copy standalone output
COPY --from=builder /app/public ./public
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static

EXPOSE 3000

# Run as non-root user (distroless default)
CMD ["server.js"]
