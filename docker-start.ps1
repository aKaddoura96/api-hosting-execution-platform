# API Platform - Docker Start Script
# Starts all services with Docker Compose

Write-Host "======================================" -ForegroundColor Cyan
Write-Host "  Starting API Platform (Docker)" -ForegroundColor Cyan
Write-Host "======================================" -ForegroundColor Cyan
Write-Host ""

try {
    # Build and start all services
    Write-Host "Building and starting all services..." -ForegroundColor Yellow
    docker-compose up --build -d
    
    if ($LASTEXITCODE -ne 0) {
        throw "Failed to start services"
    }
    
    Write-Host ""
    Write-Host "Waiting for services to be ready..." -ForegroundColor Yellow
    Start-Sleep -Seconds 10
    
    # Wait for database to be ready
    Write-Host "Waiting for database..." -ForegroundColor Yellow
    $maxAttempts = 30
    $attempt = 0
    
    while ($attempt -lt $maxAttempts) {
        $result = docker exec api-platform-postgres pg_isready -U apiplatform 2>&1
        if ($LASTEXITCODE -eq 0) {
            Write-Host "✓ Database is ready!" -ForegroundColor Green
            break
        }
        Start-Sleep -Seconds 1
        $attempt++
    }
    
    # Run migrations
    Write-Host ""
    Write-Host "Running database migrations..." -ForegroundColor Yellow
    Get-ChildItem -Path "scripts/migrations/*.sql" | Sort-Object Name | ForEach-Object {
        Write-Host "  Applying $($_.Name)..." -ForegroundColor White
        Get-Content $_.FullName | docker exec -i api-platform-postgres psql -U apiplatform -d apiplatform
    }
    
    Write-Host ""
    Write-Host "======================================" -ForegroundColor Green
    Write-Host "  ✓ API Platform Started!" -ForegroundColor Green
    Write-Host "======================================" -ForegroundColor Green
    Write-Host ""
    Write-Host "Services:" -ForegroundColor White
    Write-Host "  • Frontend:    http://localhost:3000" -ForegroundColor Cyan
    Write-Host "  • API Gateway: http://localhost:8080" -ForegroundColor Cyan
    Write-Host "  • Executor:    http://localhost:8081" -ForegroundColor Cyan
    Write-Host "  • Analytics:   http://localhost:8082" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Commands:" -ForegroundColor White
    Write-Host "  • View logs:   docker-compose logs -f" -ForegroundColor Gray
    Write-Host "  • Stop all:    docker-compose down" -ForegroundColor Gray
    Write-Host "  • Restart:     docker-compose restart" -ForegroundColor Gray
    Write-Host ""
    
} catch {
    Write-Host ""
    Write-Host "======================================" -ForegroundColor Red
    Write-Host "  ✗ Error Starting Services" -ForegroundColor Red
    Write-Host "======================================" -ForegroundColor Red
    Write-Host ""
    Write-Host "Error: $_" -ForegroundColor Red
    Write-Host ""
    Write-Host "Try: docker-compose down && docker-compose up --build" -ForegroundColor Yellow
    exit 1
}
