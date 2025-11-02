# API Platform - Stop Script
# Stops all running services

Write-Host "======================================" -ForegroundColor Cyan
Write-Host "  API Platform - Stopping Services" -ForegroundColor Cyan
Write-Host "======================================" -ForegroundColor Cyan
Write-Host ""

$ProjectRoot = Split-Path -Parent $PSScriptRoot

# Function to kill process on port
function Stop-ProcessOnPort {
    param([int]$Port)
    
    try {
        $process = Get-NetTCPConnection -LocalPort $Port -ErrorAction SilentlyContinue | Select-Object -ExpandProperty OwningProcess -Unique
        if ($process) {
            foreach ($pid in $process) {
                $proc = Get-Process -Id $pid -ErrorAction SilentlyContinue
                if ($proc) {
                    Write-Host "  Stopping process on port $Port (PID: $pid - $($proc.Name))..." -ForegroundColor Yellow
                    Stop-Process -Id $pid -Force -ErrorAction SilentlyContinue
                }
            }
            Write-Host "  ✓ Port $Port freed" -ForegroundColor Green
        }
    } catch {
        # Port might not be in use
    }
}

try {
    # Stop frontend (port 3000)
    Write-Host "Stopping frontend..." -ForegroundColor Cyan
    Stop-ProcessOnPort -Port 3000
    
    # Stop backend services
    Write-Host "Stopping backend services..." -ForegroundColor Cyan
    Stop-ProcessOnPort -Port 8080  # API Gateway
    Stop-ProcessOnPort -Port 8081  # Executor
    Stop-ProcessOnPort -Port 8082  # Analytics
    
    # Kill any remaining node/go processes
    Write-Host "Cleaning up processes..." -ForegroundColor Cyan
    Get-Process -Name node, go -ErrorAction SilentlyContinue | Where-Object { 
        $_.Path -like "*api-hosting-execution-platform*" 
    } | Stop-Process -Force -ErrorAction SilentlyContinue
    
    # Stop Docker services
    Write-Host "Stopping Docker services..." -ForegroundColor Cyan
    cd $ProjectRoot
    docker-compose down
    
    Write-Host ""
    Write-Host "======================================" -ForegroundColor Green
    Write-Host "  ✓ All Services Stopped" -ForegroundColor Green
    Write-Host "======================================" -ForegroundColor Green
    Write-Host ""
    
} catch {
    Write-Host ""
    Write-Host "Error stopping services: $_" -ForegroundColor Red
    Write-Host ""
    Write-Host "You may need to manually stop processes:" -ForegroundColor Yellow
    Write-Host "  Get-Process -Name node, go | Stop-Process -Force" -ForegroundColor Yellow
    Write-Host "  docker-compose down" -ForegroundColor Yellow
    exit 1
}
