@echo off
echo Connecting to PostgreSQL and dropping tables...

psql -h localhost -p 5432 -U postgres -d advanced_backend -c "DROP TABLE IF EXISTS pekerjaan_alumni CASCADE;"
psql -h localhost -p 5432 -U postgres -d advanced_backend -c "DROP TABLE IF EXISTS alumni CASCADE;"
psql -h localhost -p 5432 -U postgres -d advanced_backend -c "DROP TABLE IF EXISTS admin_users CASCADE;"
psql -h localhost -p 5432 -U postgres -d advanced_backend -c "DROP TABLE IF EXISTS mahasiswa CASCADE;"

echo Tables dropped successfully!
echo.
echo Now you can run the application again with: go run cmd/server/main.go
pause