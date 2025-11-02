# API Platform - Start Script
# Starts all services in the correct order

Write-Host "======================================" -ForegroundColor Cyan
Write-Host "  API Platform - Starting Services" -ForegroundColor Cyan  
Write-Host "======================================" -ForegroundColor Cyan
Write-Host ""

$ErrorActionPreference = "Stop"
$ProjectRoot = Split-Path -Parent $PSScriptRoot

# Function to check if a port is in use
function Test-Port {
    param($Port)
    try {
        $connection = Test-NetConnection -ComputerName localhost -Port $Port -WarningAction SilentlyContinue -ErrorAction SilentlyContinue
        return $connection.TcpTestSucceeded
    } catch {
        return $false
    }
}

# Function to wait for service
function Wait-ForService {
    param(
        [string]$Name,
        [int]$Port,
        [int]$MaxAttempts = 30
    )
    
    Write-Host "Waiting for $Name to be ready..." -ForegroundColor Yellow
    $attempt = 0
    while ($attempt -lt $MaxAttempts) {
        if (Test-Port -Port $Port) {
            Write-Host "✓ $Name is ready!" -ForegroundColor Green
            return $true
        }
        Start-Sleep -Seconds 1
        $attempt++
    }
    Write-Host "✗ $Name failed to start" -ForegroundColor Red
    return $false
}

try {
    # Step 1: Start Docker services
    Write-Host "Step 1: Starting Docker services..." -ForegroundColor Cyan
    cd $ProjectRoot
    docker-compose up -d
    if ($LASTEXITCODE -ne 0) {
        throw "Failed to start Docker services"
    }
    
    # Wait for PostgreSQL
    if (-not (Wait-ForService -Name "PostgreSQL" -Port 5432)) {
        throw "PostgreSQL failed to start"
    }
    
    # Wait for Redis  
    if (-not (Wait-ForService -Name "Redis" -Port 6379)) {
        throw "Redis failed to start"
    }
    
    Write-Host ""
    
    # Step 2: Run migrations
    Write-Host "Step 2: Running database migrations..." -ForegroundColor Cyan
    & "$ProjectRoot\scripts\run-migrations.ps1"
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Warning: Migrations may have already been applied" -ForegroundColor Yellow
    }
    Write-Host ""
    
    # Step 3: Build backend services
    Write-Host "Step 3: Building backend services..." -ForegroundColor Cyan
    
    Write-Host "  Building API Gateway..." -ForegroundColor White
    cd "$ProjectRoot\backend\api-gateway"
    go build -o api-gateway.exe main.go
    if ($LASTEXITCODE -ne 0) {
        throw "Failed to build API Gateway"
    }
    
    Write-Host "  Building Executor..." -ForegroundColor White
    cd "$ProjectRoot\backend\executor"
    go build -o executor.exe -mod=mod main_simple.go
    if ($LASTEXITCODE -ne 0) {
        throw "Failed to build Executor"
    }
    
    Write-Host "  Building Analytics..." -ForegroundColor White
    cd "$ProjectRoot\backend\analytics"
    go build -o analytics.exe main.go
    if ($LASTEXITCODE -ne 0) {
        throw "Failed to build Analytics"
    }
    Write-Host ""
    
    # Step 4: Start backend services
    Write-Host "Step 4: Starting backend services..." -ForegroundColor Cyan
    
    cd "$ProjectRoot"
    
    Write-Host "  Starting API Gateway (Port 8080)..." -ForegroundColor White
    Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '$ProjectRoot\backend\api-gateway'; .\api-gateway.exe" -WindowStyle Normal
    
    Write-Host "  Starting Executor (Port 8081)..." -ForegroundColor White
    Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '$ProjectRoot\backend\executor'; .\executor.exe" -WindowStyle Normal
    
    Write-Host "  Starting Analytics (Port 8082)..." -ForegroundColor White
    Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '$ProjectRoot\backend\analytics'; .\analytics.exe" -WindowStyle Normal
    
    # Wait for backends to start
    Start-Sleep -Seconds 3
    
    if (-not (Wait-ForService -Name "API Gateway" -Port 8080)) {
        throw "API Gateway failed to start"
    }
    if (-not (Wait-ForService -Name "Executor" -Port 8081)) {
        throw "Executor failed to start"
    }
    if (-not (Wait-ForService -Name "Analytics" -Port 8082)) {
        throw "Analytics failed to start"
    }
    Write-Host ""
    
    # Step 5: Start frontend
    Write-Host "Step 5: Starting frontend..." -ForegroundColor Cyan
    cd "$ProjectRoot\frontend"
    
    # Install dependencies if needed
    if (-not (Test-Path "node_modules")) {
        Write-Host "  Installing frontend dependencies..." -ForegroundColor White
        npm install
    }
    
    Write-Host "  Starting Next.js dev server (Port 3000)..." -ForegroundColor White
    Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '$ProjectRoot\frontend'; npm run dev" -WindowStyle Normal
    
    Start-Sleep -Seconds 5
    Write-Host ""
    
    # Success message
    Write-Host "======================================" -ForegroundColor Green
    Write-Host "  ✓ All Services Started Successfully!" -ForegroundColor Green
    Write-Host "======================================" -ForegroundColor Green
    Write-Host ""
    Write-Host "Services:" -ForegroundColor White
    Write-Host "  • Frontend:    http://localhost:3000" -ForegroundColor Cyan
    Write-Host "  • API Gateway: http://localhost:8080" -ForegroundColor Cyan
    Write-Host "  • Executor:    http://localhost:8081" -ForegroundColor Cyan
    Write-Host "  • Analytics:   http://localhost:8082" -ForegroundColor Cyan
    Write-Host "  • PostgreSQL:  localhost:5432" -ForegroundColor Cyan
    Write-Host "  • Redis:       localhost:6379" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "To stop all services, run: .\scripts\stop.ps1" -ForegroundColor Yellow
    Write-Host ""
    
} catch {
    Write-Host ""
    Write-Host "======================================" -ForegroundColor Red
    Write-Host "  ✗ Error Starting Services" -ForegroundColor Red
    Write-Host "======================================" -ForegroundColor Red
    Write-Host ""
    Write-Host "Error: $_" -ForegroundColor Red
    Write-Host ""
    Write-Host "To clean up, run: .\scripts\stop.ps1" -ForegroundColor Yellow
    exit 1
}
