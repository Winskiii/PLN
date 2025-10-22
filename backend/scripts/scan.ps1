Write-Host "Running security scans..."
$ErrorActionPreference = 'Continue'
$gosec = (Get-Command gosec -ErrorAction SilentlyContinue)
if (-not $gosec) { Write-Host "gosec not found; installing..."; go install github.com/securego/gosec/v2/cmd/gosec@latest }

gosec ./...
$gosecExit = $LASTEXITCODE

$govulncheck = (Get-Command govulncheck -ErrorAction SilentlyContinue)
if (-not $govulncheck) {
	Write-Host "govulncheck not found; attempting install..."
	try {
		go install golang.org/x/vuln/cmd/govulncheck@latest
	} catch {
		Write-Warning "Failed to install govulncheck (possibly offline). Skipping govulncheck."
	}
	$govulncheck = (Get-Command govulncheck -ErrorAction SilentlyContinue)
}

if ($govulncheck) {
	govulncheck ./...
	$vulnExit = $LASTEXITCODE
} else {
	Write-Warning "govulncheck unavailable; skipped."
	$vulnExit = 0
}

if ($gosecExit -ne 0 -or $vulnExit -ne 0) { exit 1 } else { exit 0 }
