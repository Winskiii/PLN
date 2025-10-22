if (-not (Test-Path .env)) { Write-Host "No .env found. Copying from .env.example"; Copy-Item .env.example .env }
Write-Host "Starting server..."
go run ./cmd/server
