#!/bin/bash

# Load environment variables if .env exists
if [ -f .env ]; then
    export $(cat .env | grep -v '#' | xargs)
fi

# Default to local database URL if not set
DATABASE_URL="${DATABASE_URL:-postgres://apiplatform:dev_password@localhost:5432/apiplatform?sslmode=disable}"

echo "Running database migrations..."

# Run all migration files in order
for file in scripts/migrations/*.sql; do
    echo "Applying migration: $file"
    psql "$DATABASE_URL" -f "$file"
    
    if [ $? -eq 0 ]; then
        echo "✓ Migration $file completed successfully"
    else
        echo "✗ Migration $file failed"
        exit 1
    fi
done

echo "All migrations completed successfully!"
