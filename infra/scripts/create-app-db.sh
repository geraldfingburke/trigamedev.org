#!/usr/bin/env bash
set -euo pipefail

# Helper script to create a new isolated database + user inside the shared Postgres container
if [ $# -lt 2 ]; then
    echo "Usage: $0 <db_name> <db_user> [db_password]"
    exit 1
fi

DB_NAME="$1"
DB_USER="$2"
DB_PASS="${3:-$(openssl rand -hex 16)}"

docker exec -i infra-postgres psql -U postgres <<EOF
DO \$\$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = '${DB_USER}') THEN
      CREATE ROLE ${DB_USER} WITH LOGIN PASSWORD '${DB_PASS}';
   END IF;
END
\$\$;

SELECT 'CREATE DATABASE ${DB_NAME} OWNER ${DB_USER}'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '${DB_NAME}')\gexec

GRANT ALL PRIVILEGES ON DATABASE ${DB_NAME} TO ${DB_USER};
EOF

echo "=========================================="
echo "Database Created!"
echo "DB_NAME:     ${DB_NAME}"
echo "DB_USER:     ${DB_USER}"
echo "DB_PASSWORD: ${DB_PASS}"
echo "Internal URL: postgres://${DB_USER}:${DB_PASS}@postgres:5432/${DB_NAME}?sslmode=disable"
echo "=========================================="
