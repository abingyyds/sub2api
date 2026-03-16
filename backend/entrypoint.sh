#!/bin/sh
set -e

echo "Running database migrations..."

# Build DATABASE_URL from individual env vars if not already set
if [ -z "$DATABASE_URL" ]; then
    DATABASE_URL="postgres://${DATABASE_USER}:${DATABASE_PASSWORD}@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_DBNAME}?sslmode=${DATABASE_SSLMODE:-disable}"
    export DATABASE_URL
    echo "Constructed DATABASE_URL from environment variables"
fi

# 使用 goose 运行迁移
if command -v goose >/dev/null 2>&1; then
    echo "Running goose migrations..."
    goose -dir ./migrations postgres "${DATABASE_URL}" up
    echo "Migrations completed successfully"
else
    echo "Warning: goose not found, skipping migrations"
fi

echo "Starting application..."
exec ./main
