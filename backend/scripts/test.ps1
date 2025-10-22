Write-Host "Running tests..."
go test ./... -race -coverprofile=coverage.out -covermode=atomic
