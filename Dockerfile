# =============================================================================
# Sub2API Multi-Stage Dockerfile (Railway-ready with config template rendering)
# =============================================================================
# Stage 1: Build frontend
# Stage 2: Build Go backend with embedded frontend
# Stage 3: Final minimal image (+ entrypoint renders /app/config.yaml from template)
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

# Generate ent ORM code and Wire DI
RUN go generate ./ent && \
    go generate ./cmd/server

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

# Runtime deps:
# - curl: healthcheck
# - tzdata, ca-certificates: tls/time
# - gettext: provides envsubst for rendering config template at startup
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    curl \
    gettext \
    && rm -rf /var/cache/apk/*

# Create non-root user
RUN addgroup -g 1000 sub2api && \
    adduser -u 1000 -G sub2api -s /bin/sh -D sub2api

WORKDIR /app

# Copy binary from builder
COPY --from=backend-builder /app/sub2api /app/sub2api

# -----------------------------------------------------------------------------
# Config template (recommended): put a file at deploy/config.railway.tpl.yaml
# If you prefer a different path, change the COPY below and the entrypoint.
# -----------------------------------------------------------------------------
COPY deploy/config.railway.tpl.yaml /app/config.tpl.yaml

# -----------------------------------------------------------------------------
# Entrypoint:
# 1) Align Railway PORT -> SERVER_PORT
# 2) Map SUB2API_* envs to runtime envs (DB_*, REDIS_*, etc.) for backward compat
# 3) Render /app/config.yaml from /app/config.tpl.yaml using envsubst
# 4) Start /app/sub2api
# -----------------------------------------------------------------------------
COPY <<'EOF' /app/entrypoint.sh
#!/bin/sh
set -eu

log(){ echo "[entrypoint] $*"; }

# -----------------------------
# 1) Port alignment
# -----------------------------
PORT_VAL="${PORT:-${SERVER_PORT:-8080}}"
export PORT="$PORT_VAL"
export SERVER_PORT="$PORT_VAL"

# -----------------------------
# 2) Env mapping (SUB2API_* -> runtime vars)
#    Do NOT overwrite if target already set.
# -----------------------------
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

# -----------------------------
# 3) Railway / proxy defaults (can be overridden by Railway Variables)
# -----------------------------
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

: "${SECURITY_URL_ALLOWLIST_ENABLED:=false}"
export SECURITY_URL_ALLOWLIST_ENABLED

: "${SECURITY_URL_ALLOWLIST_ALLOW_INSECURE_HTTP:=true}"
export SECURITY_URL_ALLOWLIST_ALLOW_INSECURE_HTTP

# -----------------------------
# 4) Render config.yaml from template (if present)
#    This matches README's config.yaml approach.
# -----------------------------
CONFIG_TPL="${CONFIG_TPL_PATH:-/app/config.tpl.yaml}"
CONFIG_OUT="${CONFIG_PATH:-/app/config.yaml}"

if [ -f "$CONFIG_TPL" ]; then
  log "rendering config: $CONFIG_OUT (from $CONFIG_TPL)"
  # envsubst replaces ${VAR} placeholders using current env
  envsubst < "$CONFIG_TPL" > "$CONFIG_OUT"
  # best-effort: prevent accidental world-readable secrets (still inside container)
  chmod 600 "$CONFIG_OUT" 2>/dev/null || true
else
  log "config template not found at $CONFIG_TPL (skip rendering)"
fi

log "SERVER_PORT=$SERVER_PORT"
log "DB=${DB_HOST:-?}:${DB_PORT:-?}/${DB_NAME:-?} sslmode=${DB_SSLMODE:-}"
log "REDIS=${REDIS_HOST:-?}:${REDIS_PORT:-?} db=${REDIS_DB:-}"
log "GIN_MODE=$GIN_MODE"
log "CONFIG_OUT=$CONFIG_OUT"

# If Sub2API supports CONFIG_PATH, keep it exported; otherwise harmless.
export CONFIG_PATH="$CONFIG_OUT"

exec /app/sub2api
EOF

# Create data directory and fix permissions
RUN chmod +x /app/entrypoint.sh && \
    mkdir -p /app/data && \
    chown -R sub2api:sub2api /app

USER sub2api

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
    CMD sh -lc 'curl -fsS http://localhost:${SERVER_PORT:-8080}/health >/dev/null'

ENTRYPOINT ["/app/entrypoint.sh"]
