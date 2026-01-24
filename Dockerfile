# =============================================================================
# Sub2API Multi-Stage Dockerfile
# =============================================================================
# Stage 1: Build frontend
# Stage 2: Build Go backend with embedded frontend
# Stage 3: Final minimal image
# =============================================================================

ARG NODE_IMAGE=node:24-alpine
ARG GOLANG_IMAGE=golang:1.25.5-alpine
ARG ALPINE_IMAGE=alpine:3.20
ARG GOPROXY=https://goproxy.cn,direct
ARG GOSUMDB=sum.golang.google.cn

# -----------------------------------------------------------------------------
# Stage 1: Frontend Builder
# -----------------------------------------------------------------------------
FROM ${NODE_IMAGE} AS frontend-builder

WORKDIR /app/frontend

# Install pnpm
RUN corepack enable && corepack prepare pnpm@latest --activate

# Install dependencies first (better caching)
COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

# Copy frontend source and build
COPY frontend/ ./
RUN pnpm run build

# -----------------------------------------------------------------------------
# Stage 2: Backend Builder
# -----------------------------------------------------------------------------
FROM ${GOLANG_IMAGE} AS backend-builder

# Build arguments for version info (set by CI)
ARG VERSION=docker
ARG COMMIT=docker
ARG DATE
ARG GOPROXY
ARG GOSUMDB

ENV GOPROXY=${GOPROXY}
ENV GOSUMDB=${GOSUMDB}

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app/backend

# Copy go mod files first (better caching)
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy backend source first
COPY backend/ ./

# Copy frontend dist from previous stage (must be after backend copy to avoid being overwritten)
COPY --from=frontend-builder /app/backend/internal/web/dist ./internal/web/dist

# Build the binary (BuildType=release for CI builds, embed frontend)
RUN CGO_ENABLED=0 GOOS=linux go build \
    -tags embed \
    -ldflags="-s -w -X main.Commit=${COMMIT} -X main.Date=${DATE:-$(date -u +%Y-%m-%dT%H:%M:%SZ)} -X main.BuildType=release" \
    -o /app/sub2api \
    ./cmd/server

# -----------------------------------------------------------------------------
# Stage 3: Final Runtime Image
# -----------------------------------------------------------------------------
FROM ${ALPINE_IMAGE}

RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    curl \
    && rm -rf /var/cache/apk/*

RUN addgroup -g 1000 sub2api && \
    adduser -u 1000 -G sub2api -s /bin/sh -D sub2api

WORKDIR /app

COPY --from=backend-builder /app/sub2api /app/sub2api

# ---- Add entrypoint that maps Railway/Sub2API env vars to runtime env vars ----
# This makes Railway Variables effective without startCommand hacks.
COPY <<'EOF' /app/entrypoint.sh
#!/bin/sh
set -eu

log(){ echo "[entrypoint] $*"; }

# Align ports: Railway uses PORT; Sub2API often uses SERVER_PORT
PORT_VAL="${PORT:-${SERVER_PORT:-8080}}"
export PORT="$PORT_VAL"
export SERVER_PORT="$PORT_VAL"

# Map your existing Railway vars (SUB2API_*) -> runtime vars (DB_*, REDIS_*, etc.)
# Do NOT overwrite if user already set the target vars explicitly.
set_if_empty() {
  key="$1"; val="$2"
  eval "cur=\${$key:-}"
  if [ -z "$cur" ] && [ -n "$val" ]; then
    export "$key=$val"
  fi
}

set_if_empty DB_HOST "${SUB2API_DB_HOST:-}"
set_if_empty DB_PORT "${SUB2API_DB_PORT:-}"
set_if_empty DB_USER "${SUB2API_DB_USER:-}"
set_if_empty DB_PASSWORD "${SUB2API_DB_PASSWORD:-}"
set_if_empty DB_NAME "${SUB2API_DB_NAME:-}"
set_if_empty DB_SSLMODE "${SUB2API_DB_SSLMODE:-}"

set_if_empty REDIS_HOST "${SUB2API_REDIS_HOST:-}"
set_if_empty REDIS_PORT "${SUB2API_REDIS_PORT:-}"
set_if_empty REDIS_PASSWORD "${SUB2API_REDIS_PASSWORD:-}"
set_if_empty REDIS_DB "${SUB2API_REDIS_DB:-}"

set_if_empty ADMIN_EMAIL "${SUB2API_ADMIN_EMAIL:-}"
set_if_empty ADMIN_PASSWORD "${SUB2API_ADMIN_PASSWORD:-}"

# JWT: set both names to be safe across versions
if [ -n "${SUB2API_JWT_SECRET:-}" ]; then
  set_if_empty JWT_SECRET "${SUB2API_JWT_SECRET}"
  export SUB2API_JWT_SECRET="${SUB2API_JWT_SECRET}"
fi

# Railway / reverse-proxy hardening for 403(OPTIONS) / scheme / host issues
# If you want strict CORS, replace "*" with your domains.
: "${GIN_MODE:=release}"
export GIN_MODE

: "${TRUSTED_PROXIES:=0.0.0.0/0}"
export TRUSTED_PROXIES

: "${GIN_TRUSTED_PROXIES:=0.0.0.0/0}"
export GIN_TRUSTED_PROXIES

: "${SUB2API_CORS_ALLOWED_ORIGINS:=*}"
export SUB2API_CORS_ALLOWED_ORIGINS

: "${SUB2API_CORS_ALLOWED_ORIGINS_JSON:=[\"*\"]}"
export SUB2API_CORS_ALLOWED_ORIGINS_JSON

# URL allowlist: Railway is behind https proxy; many deployments need relaxed checks
: "${SECURITY_URL_ALLOWLIST_ENABLED:=false}"
export SECURITY_URL_ALLOWLIST_ENABLED

: "${SECURITY_URL_ALLOWLIST_ALLOW_INSECURE_HTTP:=true}"
export SECURITY_URL_ALLOWLIST_ALLOW_INSECURE_HTTP

log "SERVER_PORT=$SERVER_PORT"
log "DB=${DB_HOST:-?}:${DB_PORT:-?}/${DB_NAME:-?} sslmode=${DB_SSLMODE:-}"
log "REDIS=${REDIS_HOST:-?}:${REDIS_PORT:-?} db=${REDIS_DB:-}"
log "GIN_MODE=$GIN_MODE"

exec /app/sub2api
EOF

RUN chmod +x /app/entrypoint.sh && \
    mkdir -p /app/data && chown -R sub2api:sub2api /app

USER sub2api

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
    CMD sh -lc 'curl -fsS http://localhost:${SERVER_PORT:-8080}/health >/dev/null'

ENTRYPOINT ["/app/entrypoint.sh"]
