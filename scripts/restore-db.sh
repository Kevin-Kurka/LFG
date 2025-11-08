#!/bin/bash

# Database restore script

set -e

# Configuration
BACKUP_FILE="${1}"
DB_HOST="${POSTGRES_HOST:-localhost}"
DB_PORT="${POSTGRES_PORT:-5432}"
DB_USER="${POSTGRES_USER:-lfg_user}"
DB_NAME="${POSTGRES_DB:-lfg_db}"

# Check if backup file is provided
if [ -z "${BACKUP_FILE}" ]; then
    echo "Usage: $0 <backup_file>"
    echo "Example: $0 ./backups/lfg_backup_20231108_120000.sql.gz"
    exit 1
fi

# Check if backup file exists
if [ ! -f "${BACKUP_FILE}" ]; then
    echo "Error: Backup file not found: ${BACKUP_FILE}"
    exit 1
fi

echo "WARNING: This will overwrite the current database!"
echo "Database: ${DB_NAME}"
echo "Host: ${DB_HOST}:${DB_PORT}"
echo "Backup file: ${BACKUP_FILE}"
echo ""
read -p "Are you sure you want to continue? (yes/no): " confirm

if [ "${confirm}" != "yes" ]; then
    echo "Restore cancelled."
    exit 0
fi

# Create temporary directory for decompression
TMP_DIR=$(mktemp -d)
TMP_FILE="${TMP_DIR}/backup.sql"

# Decompress if needed
if [[ "${BACKUP_FILE}" == *.gz ]]; then
    echo "Decompressing backup..."
    gunzip -c "${BACKUP_FILE}" > "${TMP_FILE}"
else
    cp "${BACKUP_FILE}" "${TMP_FILE}"
fi

# Restore database
echo "Restoring database..."
PGPASSWORD="${POSTGRES_PASSWORD}" pg_restore \
    -h "${DB_HOST}" \
    -p "${DB_PORT}" \
    -U "${DB_USER}" \
    -d "${DB_NAME}" \
    -c \
    "${TMP_FILE}"

# Clean up
rm -rf "${TMP_DIR}"

echo "Database restore completed successfully!"
