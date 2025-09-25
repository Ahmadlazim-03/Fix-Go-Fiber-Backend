# PowerShell script to drop all tables and recreate them
Write-Host "Connecting to PostgreSQL and dropping tables..." -ForegroundColor Yellow

$env:PGPASSWORD = ""
$dbHost = "localhost"
$dbPort = "5432"
$dbUser = "postgres"
$dbName = "advanced_backend"

try {
    Write-Host "Dropping tables..." -ForegroundColor Green
    
    psql -h $dbHost -p $dbPort -U $dbUser -d $dbName -c "DROP TABLE IF EXISTS pekerjaan_alumni CASCADE;"
    psql -h $dbHost -p $dbPort -U $dbUser -d $dbName -c "DROP TABLE IF EXISTS alumni CASCADE;"
    psql -h $dbHost -p $dbPort -U $dbUser -d $dbName -c "DROP TABLE IF EXISTS admin_users CASCADE;"
    psql -h $dbHost -p $dbPort -U $dbUser -d $dbName -c "DROP TABLE IF EXISTS mahasiswa CASCADE;"
    
    Write-Host "Tables dropped successfully!" -ForegroundColor Green
    Write-Host ""
    Write-Host "Now you can run the application again with: go run cmd/server/main.go" -ForegroundColor Cyan
}
catch {
    Write-Host "Error: $_" -ForegroundColor Red
}

# Keep window open
Read-Host "Press Enter to continue..."