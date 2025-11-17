# Work Management System - Complete Prototype

Prototype CRUD sederhana dengan security untuk manajemen tugas dan proyek.

## 🎯 Features

### Core CRUD Operations
- ✅ **Tasks Management** - Create, Read, Update, Delete tasks
- ✅ **Projects Management** - Create, Read, Update, Delete projects  
- ✅ **Users Management** - Create, Read, Update, Delete users (ADMIN only)

### Security Features
- ✅ JWT Authentication (Access Token 30min, Refresh Token 24h)
- ✅ Password Hashing (bcrypt cost 12)
- ✅ Role-Based Access Control (4 roles: SYSADMIN, ADMIN, MANAGER, EMPLOYEE)
- ✅ Account Lockout (5 failed attempts, 30min lock)
- ✅ Security Headers (HSTS, CSP, X-Frame-Options, etc.)
- ✅ Rate Limiting (100 req/min global, 5 login attempts/15min)
- ✅ Input Validation & XSS Prevention
- ✅ Audit Logging (all CRUD operations)

### Dashboard
- Task statistics (total, completed, in progress)
- Recent tasks view
- Simple overview cards

## 🛠 Tech Stack

### Backend
- **Go 1.21+** dengan Chi Router
- **MySQL 8.0+** (go-sql-driver/mysql)
- JWT (golang-jwt/jwt/v5)
- bcrypt password hashing
- Structured logging (zap)

### Frontend
- **React 18** + Vite
- **Material-UI (MUI)** untuk UI components
- Axios untuk HTTP client
- React Router untuk routing

## 📋 Prerequisites

- Go 1.21 or higher
- Node.js 18 or higher
- MySQL 8.0 or higher
- Git

## 🚀 Quick Start

### 1. Clone Repository

```bash
git clone <repository-url>
cd PLN
```

### 2. Database Setup

```bash
# Login to MySQL
mysql -u root -p

# Create database
CREATE DATABASE workmanagement;

# Exit MySQL
exit
```

### 3. Run Migrations

```bash
cd backend/migrations

# Run migrations in order
mysql -u root -p workmanagement < 001_initial_schema.sql
mysql -u root -p workmanagement < 002_seed_roles.sql
mysql -u root -p workmanagement < 003_seed_admin.sql
```

### 4. Backend Setup

```bash
cd backend

# Copy environment file
copy .env.example .env

# Edit .env and configure:
# - DATABASE_PASSWORD=your_mysql_password
# - JWT_SECRET=your-secret-key

# Install dependencies
go mod download

# Run server
go run cmd/server/main.go
```

Backend akan jalan di `http://localhost:8080`

### 5. Frontend Setup

```bash
cd frontend

# Copy environment file
copy .env.example .env

# Install dependencies
npm install

# Run dev server
npm run dev
```

Frontend akan jalan di `http://localhost:5173`

### 6. Login

Buka browser ke `http://localhost:5173`

**Default Admin Login:**
- Email: `admin@pln.co.id`
- Password: `Admin@12345`

⚠️ **PENTING**: Ganti password setelah login pertama kali!

## 📁 Project Structure

```
PLN/
├── backend/
│   ├── cmd/server/main.go           # Entry point
│   ├── internal/
│   │   ├── config/                  # Configuration
│   │   ├── database/                # MySQL connection
│   │   ├── models/                  # Data models
│   │   ├── handlers/                # HTTP handlers
│   │   ├── middleware/              # Middleware (auth, security, etc.)
│   │   ├── server/                  # Server setup & routes
│   │   └── utils/                   # Utilities (JWT, crypto, validator)
│   ├── migrations/                  # SQL migrations
│   ├── .env.example                 # Environment template
│   └── go.mod                       # Go dependencies
│
└── frontend/
    ├── src/
    │   ├── components/
    │   │   └── layout/              # Layout components
    │   ├── contexts/                # Auth context
    │   ├── pages/                   # Page components
    │   ├── services/                # API services
    │   ├── App.jsx                  # Main app
    │   └── main.jsx                 # Entry point
    ├── package.json                 # NPM dependencies
    └── vite.config.js               # Vite config
```

## 🔐 Roles & Permissions

| Role | Permissions |
|------|------------|
| **SYSADMIN** | Full access to everything |
| **ADMIN** | Manage users, all projects, all tasks |
| **MANAGER** | Create projects/tasks, manage assigned projects |
| **EMPLOYEE** | View and update assigned tasks |

## 📡 API Endpoints

### Authentication
```
POST   /api/v1/auth/login      - Login
GET    /api/v1/auth/me         - Get current user
```

### Users (ADMIN+ only)
```
GET    /api/v1/users           - List all users
POST   /api/v1/users           - Create user
GET    /api/v1/users/:id       - Get user
PUT    /api/v1/users/:id       - Update user
DELETE /api/v1/users/:id       - Delete user
```

### Projects
```
GET    /api/v1/projects        - List projects
POST   /api/v1/projects        - Create project (MANAGER+)
GET    /api/v1/projects/:id    - Get project
PUT    /api/v1/projects/:id    - Update project
DELETE /api/v1/projects/:id    - Delete project (ADMIN+)
```

### Tasks
```
GET    /api/v1/tasks           - List all tasks
GET    /api/v1/tasks/my        - My assigned tasks
GET    /api/v1/tasks/stats     - Task statistics
POST   /api/v1/tasks           - Create task (MANAGER+)
GET    /api/v1/tasks/:id       - Get task
PUT    /api/v1/tasks/:id       - Update task
DELETE /api/v1/tasks/:id       - Delete task (MANAGER+)
```

## 🔧 Development

### Backend

```bash
# Run with auto-reload (install air first)
go install github.com/cosmtrek/air@latest
air

# Run tests
go test ./...

# Build
go build -o bin/server cmd/server/main.go
```

### Frontend

```bash
# Development
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

## 🐛 Troubleshooting

### Database Connection Error

```bash
# Check MySQL is running
mysql -u root -p

# Check database exists
SHOW DATABASES;

# Verify .env configuration
DATABASE_HOST=localhost
DATABASE_PORT=3306
DATABASE_USERNAME=root
DATABASE_PASSWORD=your_password
DATABASE_DATABASE=workmanagement
```

### JWT Token Invalid

Pastikan `JWT_SECRET` di `.env` backend sama dengan yang digunakan saat generate token.

### CORS Error

Update `SECURITY_ALLOW_ORIGINS` di `.env` backend:
```
SECURITY_ALLOW_ORIGINS=http://localhost:3000,http://localhost:5173
```

### Role Get UUID Error

Jalankan migrations dalam urutan yang benar:
1. `001_initial_schema.sql` - Create tables
2. `002_seed_roles.sql` - Insert roles
3. `003_seed_admin.sql` - Create admin user

## 📝 Notes

- Ini adalah **PROTOTYPE** untuk development/testing
- Tidak termasuk: MFA, advanced reporting, file upload, email notifications
- Password policy simplified: min 8 chars, uppercase, lowercase, number
- JWT tokens disimpan di localStorage (production pakai httpOnly cookies)
- Soft delete untuk users, projects, tasks

## 🤝 Contributing

Ini adalah prototype PLN project. Untuk kontribusi:
1. Fork repository
2. Create feature branch
3. Commit changes
4. Push to branch
5. Create Pull Request

## 📄 License

Internal PLN Project

## 🆘 Support

Untuk bantuan atau pertanyaan, hubungi team development PLN.

---

**Created with ❤️ for PLN**
