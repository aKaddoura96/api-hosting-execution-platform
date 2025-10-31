# Load environment variables from .env if it exists
if (Test-Path .env) {
    Get-Content .env | ForEach-Object {
        if ($_ -match '^([^#].+?)=(.+)$') {
            [Environment]::SetEnvironmentVariable($matches[1], $matches[2])
        }
    }
}

# Default to local database URL if not set
$DATABASE_URL = $env:DATABASE_URL
if (-not $DATABASE_URL) {
    $DATABASE_URL = "postgres://apiplatform:dev_password@localhost:5432/apiplatform?sslmode=disable"
}

Write-Host "Running database migrations..." -ForegroundColor Green

# Get all migration files
$migrationFiles = Get-ChildItem -Path "scripts\migrations\*.sql" | Sort-Object Name

foreach ($file in $migrationFiles) {
    Write-Host "Applying migration: $($file.Name)" -ForegroundColor Cyan
    
    # Use psql to run the migration
    psql $DATABASE_URL -f $file.FullName
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✓ Migration $($file.Name) completed successfully" -ForegroundColor Green
    } else {
        Write-Host "✗ Migration $($file.Name) failed" -ForegroundColor Red
        exit 1
    }
}

Write-Host "All migrations completed successfully!" -ForegroundColor Green
