param(
  [switch]$InstallTools
)

Write-Host "Setting up Go backend..."
if ($InstallTools) {
  Write-Host "Installing tools (golangci-lint, gosec, govulncheck)..."
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  go install github.com/securego/gosec/v2/cmd/gosec@latest
  go install golang.org/x/vuln/cmd/govulncheck@latest
}
Write-Host "Done."
