#!/bin/bash

# Database migration script using golang-migrate

set -e

# Configuration
DB_HOST=${POSTGRES_HOST:-localhost}
DB_PORT=${POSTGRES_PORT:-5432}
DB_USER=${POSTGRES_USER:-lfg_user}
DB_NAME=${POSTGRES_DB:-lfg_db}
DB_PASSWORD=${POSTGRES_PASSWORD}
MIGRATIONS_DIR="./database/migrations"

# Build connection string
DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

# Check if golang-migrate is installed
if ! command -v migrate &> /dev/null; then
    echo "golang-migrate not found. Installing..."
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
    sudo mv migrate /usr/local/bin/migrate
    chmod +x /usr/local/bin/migrate
fi

# Function to display usage
usage() {
    echo "Usage: $0 {up|down|force|version|create}"
    echo "  up         - Apply all pending migrations"
    echo "  down       - Rollback the last migration"
    echo "  force N    - Force migration to version N"
    echo "  version    - Show current migration version"
    echo "  create NAME - Create new migration files"
    exit 1
}

# Execute migration command
case "$1" in
    up)
        echo "Applying migrations..."
        migrate -path ${MIGRATIONS_DIR} -database "${DATABASE_URL}" up
        echo "Migrations applied successfully!"
        ;;
    down)
        echo "Rolling back last migration..."
        migrate -path ${MIGRATIONS_DIR} -database "${DATABASE_URL}" down 1
        echo "Rollback completed!"
        ;;
    force)
        if [ -z "$2" ]; then
            echo "Error: Version number required"
            usage
        fi
        echo "Forcing migration to version $2..."
        migrate -path ${MIGRATIONS_DIR} -database "${DATABASE_URL}" force $2
        echo "Migration forced to version $2"
        ;;
    version)
        migrate -path ${MIGRATIONS_DIR} -database "${DATABASE_URL}" version
        ;;
    create)
        if [ -z "$2" ]; then
            echo "Error: Migration name required"
            usage
        fi
        migrate create -ext sql -dir ${MIGRATIONS_DIR} -seq $2
        echo "Migration files created for: $2"
        ;;
    *)
        usage
        ;;
esac
