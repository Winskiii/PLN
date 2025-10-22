# Threat Model Template

## Context
- System: Secure API backend
- Data: Authentication tokens (if added later), user data (if applicable)
- Actors: Users, attackers, internal services

## Assets
- Service availability
- Data integrity
- Logs and audit trails

## Trust boundaries
- Internet <-> API gateway/server
- Server <-> datastore (future)

## Threats (STRIDE)
- Spoofing: Missing auth, inadequate session handling
- Tampering: Input manipulation, dependency compromise
- Repudiation: Insufficient logging
- Information disclosure: Overly broad CORS, verbose errors
- Denial of Service: Flooding, large bodies
- Elevation of privilege: Unsafe defaults

## Mitigations
- CORS allowlist, security headers, structured logs with correlation IDs
- Validation and size limits, rate limiting, panic recovery
- Dependency scanning and updates
- Principle of least privilege in runtime and container

## Assumptions
- No database yet (add controls when introduced)

## Residual Risks
- Zero-day dependency vulns between scans
- Misconfiguration in deployment environments

## Validation/Testing
- Unit tests
- Lint + security scanners
- Manual abuse-case testing
