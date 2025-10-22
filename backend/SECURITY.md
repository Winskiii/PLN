# Security Policy

## Reporting a Vulnerability
If you discover a security issue, please do not file a public issue. Instead, email [REPLACE_WITH_CONTACT@example.com] with details. We will acknowledge and work to resolve the issue promptly.

Provide:
- Affected version/commit
- Reproduction steps, PoC
- Impact assessment

## Supported Versions
We support the latest main branch and the most recent tagged release.

## Hardening
- Non-root containers
- Minimal base images
- HTTP security headers
- Input validation and rate limiting
- Dependencies scanned with govulncheck and gosec

## Secrets
Do not commit secrets. Use environment variables or a secret manager. `.env` is for local development only.
