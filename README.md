# PLN - Secure Web Application

Aplikasi web aman dengan penerapan **Secure Software Development Life Cycle (SSDLC)** menggunakan **Go (Backend)** dan **React + Vite (Frontend)**.

---

## 📋 Daftar Isi

- [Overview](#overview)
- [Tech Stack](#tech-stack)
- [Aspek SSDLC yang Diimplementasikan](#aspek-ssdlc-yang-diimplementasikan)
- [Struktur Project](#struktur-project)
- [Quick Start](#quick-start)
- [Security Features](#security-features)
- [Dokumentasi](#dokumentasi)
- [Testing](#testing)
- [Deployment](#deployment)
- [Contributing](#contributing)

---

## 🎯 Overview

Project ini adalah implementasi lengkap **Secure Software Development Life Cycle (SSDLC)** untuk aplikasi web PLN. Setiap fase development mengintegrasikan praktik keamanan terbaik (security best practices) dari perencanaan hingga maintenance.

### Tujuan Project:
- ✅ Menerapkan OWASP Top 10 mitigation
- ✅ Implementasi secure coding standards
- ✅ Automated security testing
- ✅ Compliance-ready architecture
- ✅ Production-grade security controls

---

## 🛠️ Tech Stack

### Backend
- **Go 1.21+** - High-performance backend
- **Chi Router** - Lightweight HTTP router
- **JWT** - Token-based authentication
- **Viper** - Configuration management
- **Zap** - Structured logging

### Frontend
- **React 18** - Modern UI framework
- **Vite** - Lightning-fast build tool
- **Axios** - HTTP client
- **React Router** - Client-side routing

### Security Tools
- **golangci-lint** - Static code analysis
- **gosec** - Go security checker
- **govulncheck** - Vulnerability scanning
- **Trivy** - Container security scanner

---

## 🔒 Aspek SSDLC yang Diimplementasikan

### 📋 **Security Documentation**

| Dokumen | Lokasi | Deskripsi |
|---------|--------|-----------|
| **Threat Model** | [`backend/docs/threat-model.md`](backend/docs/threat-model.md) | Analisis ancaman menggunakan STRIDE methodology |
| **SSDLC Checklist** | [`backend/docs/ssdlc-checklist.md`](backend/docs/ssdlc-checklist.md) | Checklist lengkap setiap fase SSDLC |
| **Security Policy** | [`backend/SECURITY.md`](backend/SECURITY.md) | Kebijakan keamanan dan vulnerability reporting |
| **API Documentation** | [`backend/docs/openapi.yaml`](backend/docs/openapi.yaml) | OpenAPI 3.0 specification dengan security schemes |
| **Frontend Integration** | [`backend/docs/frontend-integration.md`](backend/docs/frontend-integration.md) | Panduan integrasi frontend dengan security considerations |
| **Pedoman SSDLC Lengkap** | [`docs/Pedoman_SSDLC_Lengkap.md`](docs/Pedoman_SSDLC_Lengkap.md) | Panduan lengkap SSDLC dalam Bahasa Indonesia |

### 🔒 **Security Middleware**

| Middleware | Lokasi | Fungsi |
|------------|--------|--------|
| **Security Headers** | [`backend/internal/middleware/security.go`](backend/internal/middleware/security.go) | X-Content-Type-Options, X-Frame-Options, CSP, HSTS, dll |
| **JWT Authentication** | [`backend/internal/middleware/auth.go`](backend/internal/middleware/auth.go) | Token validation & user authorization |
| **Panic Recovery** | [`backend/internal/middleware/recover.go`](backend/internal/middleware/recover.go) | Graceful error handling & logging |
| **Request Logging** | [`backend/internal/middleware/logging.go`](backend/internal/middleware/logging.go) | Structured logging dengan correlation IDs |

### 🧪 **Testing & Quality Assurance**

| Test Type | Lokasi | Coverage |
|-----------|--------|----------|
| **Middleware Tests** | [`backend/internal/middleware/security_test.go`](backend/internal/middleware/security_test.go) | 100% |
| **Handler Tests** | [`backend/internal/handlers/health_test.go`](backend/internal/handlers/health_test.go) | 95%+ |
| **Auth Tests** | [`backend/internal/handlers/auth_test.go`](backend/internal/handlers/auth_test.go) | 90%+ |
| **Coverage Report** | [`backend/coverage/`](backend/coverage/) | HTML & JSON reports |

**Menjalankan Tests:**
```powershell
cd backend
.\scripts\test.ps1
```

### 🔍 **Security Scanning Tools**

| Tool | Config File | Purpose |
|------|-------------|---------|
| **golangci-lint** | [`backend/.golangci.yml`](backend/.golangci.yml) | 50+ linters untuk code quality & security |
| **gosec** | [`backend/scripts/scan.ps1`](backend/scripts/scan.ps1) | Go security scanner (SAST) |
| **govulncheck** | [`backend/scripts/scan.ps1`](backend/scripts/scan.ps1) | Go vulnerability database scanner |
| **Pre-commit Hooks** | [`backend/.pre-commit-config.yaml`](backend/.pre-commit-config.yaml) | Automated checks sebelum commit |

**Menjalankan Security Scan:**
```powershell
cd backend
.\scripts\scan.ps1
```

### 🚀 **CI/CD Security Pipeline**

| Stage | File | Checks |
|-------|------|--------|
| **GitHub Actions** | [`backend/.github/workflows/ci.yml`](backend/.github/workflows/ci.yml) | Build, test, lint, security scan |
| **Pre-commit** | [`backend/.pre-commit-config.yaml`](backend/.pre-commit-config.yaml) | Local validation sebelum push |

**Pipeline meliputi:**
- ✅ Static Application Security Testing (SAST)
- ✅ Dependency vulnerability scanning
- ✅ Code quality checks
- ✅ Unit & integration tests
- ✅ Container security scanning

### 🐳 **Container Security**

| File | Security Features |
|------|-------------------|
| [`backend/Dockerfile`](backend/Dockerfile) | Multi-stage build, non-root user, minimal base image |
| [`backend/.dockerignore`](backend/.dockerignore) | Exclude sensitive files from image |

**Security Best Practices:**
- ✅ Distroless/Alpine base image
- ✅ Non-root user execution
- ✅ Read-only filesystem (where possible)
- ✅ Health checks implemented
- ✅ No hardcoded secrets

### ⚙️ **Configuration Management**

| File | Purpose |
|------|---------|
| [`backend/internal/config/config.go`](backend/internal/config/config.go) | Configuration struct dengan validation |
| [`backend/internal/config/load.go`](backend/internal/config/load.go) | Secure config loading dari env vars |
| [`backend/.env.example`](backend/.env.example) | Template environment variables |

**Security Features:**
- ✅ No hardcoded secrets
- ✅ Environment-based configuration
- ✅ Validation untuk semua config values
- ✅ Secure defaults

### 🌐 **Server Security**

| Komponen | Lokasi | Features |
|----------|--------|----------|
| **Main Server** | [`backend/internal/server/server.go`](backend/internal/server/server.go) | CORS, rate limiting, timeouts |
| **Entry Point** | [`backend/cmd/server/main.go`](backend/cmd/server/main.go) | Graceful shutdown, signal handling |

**Implemented Controls:**
- ✅ CORS configuration
- ✅ Rate limiting (100 req/min per IP)
- ✅ Request body size limits (10MB)
- ✅ Read/Write timeouts
- ✅ Graceful shutdown

---

## 📁 Struktur Project

```
PLN/
├── backend/                    # Go Backend
│   ├── cmd/
│   │   └── server/
│   │       └── main.go         # Entry point
│   ├── internal/
│   │   ├── config/             # Configuration management
│   │   ├── handlers/           # HTTP handlers
│   │   ├── middleware/         # Security middleware
│   │   ├── models/             # Data models
│   │   └── server/             # Server setup
│   ├── docs/                   # Documentation
│   │   ├── threat-model.md     # Threat modeling
│   │   ├── ssdlc-checklist.md  # SSDLC checklist
│   │   ├── openapi.yaml        # API spec
│   │   └── frontend-integration.md
│   ├── scripts/                # Utility scripts
│   │   ├── lint.ps1
│   │   ├── scan.ps1
│   │   └── test.ps1
│   ├── data/                   # Data files
│   │   └── dummy.csv           # Sample CSV data
│   ├── .golangci.yml           # Linter config
│   ├── .pre-commit-config.yaml # Pre-commit hooks
│   ├── Dockerfile              # Container image
│   └── SECURITY.md             # Security policy
│
├── frontend/                   # React Frontend
│   ├── src/
│   │   ├── App.js              # Main component
│   │   ├── CsvTable.js         # CSV data display
│   │   └── index.js            # Entry point
│   ├── public/                 # Static assets
│   ├── vite.config.js          # Vite configuration
│   ├── eslint.config.js        # ESLint rules
│   └── package.json
│
├── docs/                       # Project documentation
│   └── Pedoman_SSDLC_Lengkap.md # Comprehensive SSDLC guide (ID)
│
└── README.md                   # This file
```

---

## 🚀 Quick Start

### Prerequisites

- **Go 1.21+**
- **Node.js 18+**
- **npm atau yarn**

### 1. Clone Repository

```powershell
git clone https://github.com/Winskiii/PLN.git
cd PLN
```

### 2. Setup Backend

```powershell
cd backend

# Install dependencies
go mod download

# Copy environment template
copy .env.example .env

# Edit .env sesuai kebutuhan (optional)
notepad .env

# Run backend
go run cmd/server/main.go
```

Backend akan berjalan di **http://localhost:8080**

### 3. Setup Frontend

```powershell
cd frontend

# Install dependencies
npm install

# Run development server
npm run dev
```

Frontend akan berjalan di **http://localhost:3001**

### 4. Akses Aplikasi

Buka browser dan akses:
- **Frontend**: http://localhost:3001
- **Backend Health**: http://localhost:8080/healthz
- **API Health**: http://localhost:8080/api/health

---

## 🛡️ Security Features

### 1. Authentication & Authorization
- ✅ JWT-based authentication
- ✅ Secure password hashing (bcrypt)
- ✅ Token expiration & refresh
- ✅ Role-Based Access Control (RBAC)

### 2. Input Validation & Sanitization
- ✅ Request body size limits
- ✅ Content-Type validation
- ✅ SQL injection prevention (prepared statements)
- ✅ XSS prevention (output encoding)

### 3. Secure Headers
```
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Content-Security-Policy: default-src 'self'
Strict-Transport-Security: max-age=31536000
Referrer-Policy: strict-origin-when-cross-origin
Permissions-Policy: geolocation=(), microphone=(), camera=()
```

### 4. Rate Limiting
- 100 requests per minute per IP
- Configurable via environment variables
- Prevents brute force & DoS attacks

### 5. CORS Protection
- Whitelist-based origins
- Credentials support
- Preflight request handling

### 6. Logging & Monitoring
- Structured logging (JSON format)
- Correlation IDs untuk request tracking
- Security event logging
- Error logging tanpa expose sensitive data

### 7. Data Protection
- ✅ HTTPS enforcement (production)
- ✅ Secure cookie settings
- ✅ Environment-based secrets
- ✅ No hardcoded credentials

---

## 📚 Dokumentasi

### Backend Documentation

- **[Threat Model](backend/docs/threat-model.md)** - STRIDE analysis & security controls
- **[SSDLC Checklist](backend/docs/ssdlc-checklist.md)** - Development phase checklist
- **[API Specification](backend/docs/openapi.yaml)** - OpenAPI 3.0 spec
- **[Frontend Integration](backend/docs/frontend-integration.md)** - Integration guide
- **[Security Policy](backend/SECURITY.md)** - Vulnerability reporting

### Comprehensive Guide

- **[Pedoman SSDLC Lengkap](docs/Pedoman_SSDLC_Lengkap.md)** - Panduan lengkap SSDLC (Bahasa Indonesia)

### API Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/healthz` | Liveness probe | ❌ |
| GET | `/readyz` | Readiness probe | ❌ |
| GET | `/api/health` | Frontend health check | ❌ |
| GET | `/api/csv` | Get CSV data | ❌ |
| POST | `/api/v1/auth/login` | User login | ❌ |
| GET | `/api/v1/auth/me` | Get user info | ✅ |
| GET | `/api/v1/ping` | API ping | ❌ |

---

## 🧪 Testing

### Unit Tests

```powershell
cd backend

# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Generate coverage report
.\scripts\test.ps1
```

Coverage reports tersedia di: `backend/coverage/coverage.html`

### Security Scanning

```powershell
cd backend

# Run linter
.\scripts\lint.ps1

# Run security scan
.\scripts\scan.ps1
```

### Manual Testing

**Test Security Headers:**
```powershell
curl -I http://localhost:8080/healthz
```

**Test Rate Limiting:**
```powershell
# Kirim 101 requests (akan ke-block di request ke-101)
for ($i=1; $i -le 101; $i++) { 
    curl http://localhost:8080/api/health 
}
```

**Test CORS:**
```powershell
curl -H "Origin: http://malicious.com" http://localhost:8080/api/health
# Harusnya di-reject
```

---

## 🚢 Deployment

### Docker Deployment

**Build Image:**
```powershell
cd backend
docker build -t pln-backend:latest .
```

**Run Container:**
```powershell
docker run -d \
  -p 8080:8080 \
  -e JWT_SECRET=your-secret-here \
  -e DATABASE_URL=postgres://... \
  --name pln-backend \
  pln-backend:latest
```

### Production Checklist

```
□ Environment variables di-set dengan aman
□ HTTPS enabled (TLS certificate valid)
□ Database credentials menggunakan secrets management
□ Rate limiting configured sesuai traffic
□ Monitoring & alerting aktif
□ Backup strategy in place
□ Incident response plan ready
□ Security logging enabled
□ WAF deployed (recommended)
□ Regular security audits scheduled
```

---

## 🤝 Contributing

### Development Workflow

1. **Create Feature Branch**
   ```powershell
   git checkout -b feature/your-feature-name
   ```

2. **Make Changes & Test**
   ```powershell
   # Backend tests
   cd backend
   .\scripts\test.ps1
   .\scripts\lint.ps1
   .\scripts\scan.ps1
   ```

3. **Commit Changes**
   ```powershell
   git add .
   git commit -m "feat: add new feature"
   ```

4. **Push & Create PR**
   ```powershell
   git push origin feature/your-feature-name
   ```

### Commit Message Convention

```
feat: new feature
fix: bug fix
docs: documentation changes
style: code style changes (formatting)
refactor: code refactoring
test: add or update tests
chore: maintenance tasks
security: security improvements
```

### Code Review Checklist

```
□ Code follows project conventions
□ All tests passing
□ Security scan passing (no new vulnerabilities)
□ Documentation updated
□ No hardcoded secrets
□ Input validation implemented
□ Error handling proper
□ Logging added for important events
```

---

## 📊 Security Metrics

### Current Status

| Metric | Status | Target |
|--------|--------|--------|
| **Test Coverage** | 90%+ | 95% |
| **Security Scan** | ✅ Pass | 0 Critical/High |
| **Linter Issues** | ✅ Pass | 0 Errors |
| **Dependencies** | ✅ Up-to-date | No known CVEs |
| **OWASP Top 10** | ✅ Mitigated | All covered |

### Continuous Improvement

- 🎯 Weekly dependency updates
- 🎯 Monthly security audits
- 🎯 Quarterly penetration testing
- 🎯 Continuous monitoring & alerting

---

## 📞 Support & Contact

### Security Issues

Untuk melaporkan vulnerability, silakan baca [SECURITY.md](backend/SECURITY.md)

**DO NOT** open public issue untuk security vulnerabilities!

### General Issues

Untuk bug reports atau feature requests:
1. Check existing issues
2. Create new issue dengan template yang sesuai
3. Provide detailed information

---

## 📄 License

[Sesuaikan dengan license project Anda]

---

## 🙏 Acknowledgments

- OWASP untuk security guidelines
- Go community untuk excellent tools
- React team untuk modern frontend framework

---

## 📈 Roadmap

### Phase 1: Foundation ✅
- [x] Basic SSDLC implementation
- [x] Security middleware
- [x] Authentication & authorization
- [x] Security testing automation

### Phase 2: Enhancement 🚧
- [ ] Advanced threat monitoring
- [ ] Automated incident response
- [ ] Security dashboard
- [ ] Compliance automation

### Phase 3: Scale 📋
- [ ] Multi-region deployment
- [ ] Advanced caching
- [ ] Real-time analytics
- [ ] AI-powered security

---

**Built with ❤️ and 🔒 by following SSDLC best practices**
