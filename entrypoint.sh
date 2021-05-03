#!/usr/bin/env bash
set -e

echo "Running Database migrations"
migrate -version
migrate -path=migrations -database postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_DATABASE?sslmode=disable up

exec "$@"
