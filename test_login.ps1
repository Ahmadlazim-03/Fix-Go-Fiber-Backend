#!/usr/bin/env pwsh

# Start server in background
Write-Host "Starting server..." -ForegroundColor Green
$job = Start-Job -ScriptBlock { Set-Location "C:\PROJECT_GITHUB\Fix-Go-Fiber-Backend"; ./server.exe }

# Wait for server to start
Start-Sleep -Seconds 5

# Test admin login
Write-Host "Testing admin login..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/admin/login" -Method POST -ContentType "application/json" -Body '{"username":"admin","password":"admin123"}'
    Write-Host "✅ Login successful!" -ForegroundColor Green
    Write-Host "Response: $($response | ConvertTo-Json -Depth 3)" -ForegroundColor Cyan
} catch {
    Write-Host "❌ Login failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Clean up
Write-Host "Stopping server..." -ForegroundColor Yellow
Stop-Job -Job $job
Remove-Job -Job $job