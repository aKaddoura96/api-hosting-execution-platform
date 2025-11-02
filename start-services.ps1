# Start all backend services
Write-Host "Starting API Platform Services..." -ForegroundColor Green

# Start API Gateway
Write-Host "`nStarting API Gateway (Port 8080)..." -ForegroundColor Cyan
Start-Process pwsh -ArgumentList "-NoExit", "-Command", "cd '$PSScriptRoot\backend\api-gateway'; go run main.go"

# Wait a bit for Gateway to start
Start-Sleep -Seconds 2

# Start Executor Service
Write-Host "Starting Executor Service (Port 8081)..." -ForegroundColor Cyan
Start-Process pwsh -ArgumentList "-NoExit", "-Command", "cd '$PSScriptRoot\backend\executor'; go run main.go"

# Wait a bit
Start-Sleep -Seconds 2

# Start Analytics Service
Write-Host "Starting Analytics Service (Port 8082)..." -ForegroundColor Cyan
Start-Process pwsh -ArgumentList "-NoExit", "-Command", "cd '$PSScriptRoot\backend\analytics'; go run main.go"

Write-Host "`n✅ All services starting in separate windows!" -ForegroundColor Green
Write-Host "Check each window for service status`n" -ForegroundColor Yellow
