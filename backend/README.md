# Secure Go Backend (SSDLC)

This backend is a secure-by-default Go (1.22) API scaffold designed to integrate practical SSDLC (Secure Software Development Lifecycle) practices.

## Features
- HTTP server using chi router with:
  - `/healthz`, `/readyz`, and `/api/v1/ping` endpoints
  - CORS, security headers, panic recovery, request logging, correlation IDs
  - per-IP rate limiting, request validation helpers
- Config via env vars and `.env` (for local dev)
- Structured logging (zerolog)
- SSDLC tooling: golangci-lint, gosec, govulncheck, unit tests
- CI workflow (GitHub Actions) for build/test/scan
- Dockerfile (multi-stage) running as non-root

## Getting Started

1. Copy `.env.example` to `.env` and adjust values.
2. Install Go 1.22+.
3. Install tools (optional): golangci-lint, gosec, govulncheck.

### Run locally

PowerShell:

```powershell
$env:APP_ENV="development"; go run ./cmd/server
```

### Run tests

```powershell
pwsh ./scripts/test.ps1
```

### Lint

```powershell
pwsh ./scripts/lint.ps1
```

### Security scans

```powershell
pwsh ./scripts/scan.ps1
```

### Docker

```powershell
# Build
docker build -t secure-api:local .
# Run
docker run --rm -p 8080:8080 --env-file .env secure-api:local
```

## SSDLC Artifacts
- `docs/threat-model.md` – template to capture threats, mitigations
- `docs/ssdlc-checklist.md` – checklist across the lifecycle
- `SECURITY.md` – vulnerability disclosure policy

## API
- GET `/healthz` – liveness
- GET `/readyz` – readiness
- GET `/api/v1/ping` – returns `{ "message": "pong" }`
 - POST `/api/v1/auth/login` – body `{"username","password"}` returns `{ "token": "..." }`
 - GET `/api/v1/auth/me` – requires `Authorization: Bearer <token>`; returns `{ "user": "..." }`

## Notes
- Keep secrets out of git; use `.env` for local only and secret manager in prod.
- Update CORS origins via `ALLOW_ORIGINS`.

## Frontend wiring (React)
During development you can call the backend directly at `http://localhost:8080`. Example:

```ts
// fetch ping
const res = await fetch('http://localhost:8080/api/v1/ping');
const data = await res.json();

// login
const loginRes = await fetch('http://localhost:8080/api/v1/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ username: 'admin', password: 'password123' }), // replace
});
const { token } = await loginRes.json();

// me (authorized)
const meRes = await fetch('http://localhost:8080/api/v1/auth/me', {
  headers: { Authorization: `Bearer ${token}` },
});
```

If you prefer to avoid CORS preflights in dev, configure your React dev server proxy to the backend (create `src/setupProxy.js` in CRA or Vite proxy settings) and remove the domain from requests.
