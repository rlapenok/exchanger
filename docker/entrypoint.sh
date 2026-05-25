#!/bin/sh
set -eu

echo "Waiting for postgres at ${DB_HOST}:${DB_PORT}..."
until nc -z "$DB_HOST" "$DB_PORT"; do
  sleep 1
done

DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE:-disable}&search_path=${DB_SCHEMA:-public}&connect_timeout=${DB_CONNECTION_TIMEOUT:-5}"

echo "Running migrations..."
migrate -path /app/migrations -database "$DATABASE_URL" up

echo "Starting exchanger..."
exec /app/exchanger
