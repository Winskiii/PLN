# SSDLC Checklist

## Requirements
- [ ] Security non-functionals captured
- [ ] Threat model created/updated

## Design
- [ ] Secure-by-default patterns
- [ ] Data classification and handling

## Implementation
- [ ] Input validation
- [ ] Error handling avoids sensitive info
- [ ] Logging is structured with correlation IDs
- [ ] Secrets not in source control

## Verification
- [ ] Unit tests
- [ ] Lint (golangci-lint)
- [ ] Static analysis (gosec)
- [ ] Vulnerability scan (govulncheck)

## Release
- [ ] CI enforces checks
- [ ] Container runs as non-root
- [ ] Minimal images

## Operations
- [ ] Logging and monitoring hooks
- [ ] Incident response plan
- [ ] Regular dependency updates
