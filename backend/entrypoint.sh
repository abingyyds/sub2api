#!/bin/sh
set -e

echo "Running database migrations..."

# 使用 goose 运行迁移
if command -v goose >/dev/null 2>&1; then
    goose -dir ./migrations postgres "${DATABASE_URL}" up
else
    echo "Warning: goose not found, skipping migrations"
fi

echo "Starting application..."
exec ./main
