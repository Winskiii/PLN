Write-Host "Running linters..."
# Ensure Go toolchain supports 1.24 for linters
$env:GOTOOLCHAIN = "go1.24.9"
$golangci = (Get-Command golangci-lint -ErrorAction SilentlyContinue)
if (-not $golangci) { 
	Write-Host "golangci-lint not found; installing with toolchain $env:GOTOOLCHAIN..."; 
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest 
}

golangci-lint run
