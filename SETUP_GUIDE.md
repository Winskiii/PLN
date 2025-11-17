# Setup Instructions - Work Management System

## Step-by-Step Setup Guide

### 1. Install Prerequisites

#### Windows:
- **Go**: Download from https://go.dev/dl/ (v1.21+)
- **Node.js**: Download from https://nodejs.org/ (v18+)
- **MySQL**: Download from https://dev.mysql.com/downloads/mysql/ (v8.0+)

Verify installations:
```powershell
go version
node --version
npm --version
mysql --version
```

### 2. Database Setup

```powershell
# Start MySQL service (if not running)
net start MySQL80

# Login to MySQL
mysql -u root -p
# Enter your MySQL root password

# In MySQL prompt:
CREATE DATABASE workmanagement CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
SHOW DATABASES;
EXIT;
```

### 3. Run Database Migrations

```powershell
cd backend\migrations

# Run each migration file
mysql -u root -p workmanagement < 001_initial_schema.sql
mysql -u root -p workmanagement < 002_seed_roles.sql
mysql -u root -p workmanagement < 003_seed_admin.sql

# Verify tables created
mysql -u root -p workmanagement -e "SHOW TABLES;"
```

Expected output:
```
+---------------------------+
| Tables_in_workmanagement  |
+---------------------------+
| audit_logs                |
| projects                  |
| roles                     |
| tasks                     |
| users                     |
+---------------------------+
```

### 4. Backend Configuration

```powershell
cd ..\  # Back to backend folder

# Copy environment file
copy .env.example .env

# Edit .env file (use notepad or any editor)
notepad .env
```

Update these values in `.env`:
```ini
DATABASE_PASSWORD=your_mysql_password
JWT_SECRET=your-super-secret-jwt-key-min-32-characters-long
```

### 5. Install Backend Dependencies

```powershell
# In backend folder
go mod download

# Verify dependencies
go list -m all
```

### 6. Run Backend Server

```powershell
# In backend folder
go run cmd/server/main.go
```

Expected output:
```
{"level":"info","ts":...,"msg":"database connected","host":"localhost","port":3306,"database":"workmanagement"}
{"level":"info","ts":...,"msg":"server starting","address":"0.0.0.0:8080"}
```

Keep this terminal running!

### 7. Frontend Configuration

Open **NEW** PowerShell window:

```powershell
cd PLN\frontend

# Copy environment file
copy .env.example .env

# Edit if needed (default is fine)
notepad .env
```

### 8. Install Frontend Dependencies

```powershell
# In frontend folder
npm install
```

Wait for installation to complete...

### 9. Run Frontend Dev Server

```powershell
# In frontend folder
npm run dev
```

Expected output:
```
  VITE v5.0.11  ready in 500 ms

  ➜  Local:   http://localhost:5173/
  ➜  Network: use --host to expose
```

### 10. Test the Application

1. Open browser to `http://localhost:5173`
2. You should see login page
3. Login with:
   - **Email**: `admin@pln.co.id`
   - **Password**: `Admin@12345`
4. After successful login, you should see Dashboard

### 11. Verify Functionality

Test each feature:

#### Dashboard
- ✅ See task statistics
- ✅ View recent tasks

#### Tasks Page
- ✅ Click "Add Task" button (as ADMIN/MANAGER)
- ✅ Fill form and create task
- ✅ Edit existing task
- ✅ Delete task

#### Projects Page
- ✅ Click "Add Project" button
- ✅ Create new project
- ✅ Edit project
- ✅ Delete project

#### Users Page
- ✅ View all users
- ✅ Create new user
- ✅ Edit user
- ✅ Delete user

## Common Issues & Solutions

### Issue 1: "go: command not found"
**Solution**: Go not installed or not in PATH
- Reinstall Go
- Add Go to PATH: `C:\Go\bin`
- Restart PowerShell

### Issue 2: "Access denied for user 'root'@'localhost'"
**Solution**: Wrong MySQL password
- Check password in `.env`
- Try connecting manually: `mysql -u root -p`

### Issue 3: "Error 1049: Unknown database 'workmanagement'"
**Solution**: Database not created
- Run: `mysql -u root -p -e "CREATE DATABASE workmanagement;"`

### Issue 4: Port 8080 already in use
**Solution**: Another process using port
- Stop other services using port 8080
- OR change port in `.env`: `SERVER_PORT=8081`

### Issue 5: CORS error in browser
**Solution**: Frontend can't connect to backend
- Check backend is running on port 8080
- Check `.env` has correct CORS origins
- Restart both frontend and backend

### Issue 6: "Cannot find module" in frontend
**Solution**: Dependencies not installed
- Delete `node_modules` folder
- Run `npm install` again

### Issue 7: "Invalid token" after login
**Solution**: JWT secret mismatch
- Make sure `JWT_SECRET` in `.env` is at least 32 characters
- Restart backend server

## Development Tips

### Hot Reload Backend (Optional)
```powershell
# Install air
go install github.com/cosmtrek/air@latest

# Run with hot reload
air
```

### View Backend Logs
Logs will show all HTTP requests:
```
{"level":"info","msg":"http_request","method":"POST","path":"/api/v1/auth/login","status":200}
```

### View Database Data
```powershell
mysql -u root -p workmanagement

# In MySQL:
SELECT * FROM users;
SELECT * FROM projects;
SELECT * FROM tasks;
SELECT * FROM audit_logs ORDER BY created_at DESC LIMIT 10;
```

### Reset Database
```powershell
mysql -u root -p -e "DROP DATABASE workmanagement; CREATE DATABASE workmanagement;"

# Then run migrations again
cd backend\migrations
mysql -u root -p workmanagement < 001_initial_schema.sql
mysql -u root -p workmanagement < 002_seed_roles.sql
mysql -u root -p workmanagement < 003_seed_admin.sql
```

## Next Steps

1. ✅ Change admin password after first login
2. ✅ Create test users with different roles
3. ✅ Create sample projects
4. ✅ Create sample tasks
5. ✅ Test role permissions

## Production Deployment (Future)

This is a PROTOTYPE. For production:
- [ ] Use environment variables for secrets
- [ ] Enable HTTPS/TLS
- [ ] Use httpOnly cookies for JWT
- [ ] Add rate limiting per user
- [ ] Enable audit log monitoring
- [ ] Set up database backups
- [ ] Add health check endpoints
- [ ] Configure proper CORS
- [ ] Add request validation middleware
- [ ] Set up monitoring/alerting

---

Need help? Check `README_PROTOTYPE.md` or contact development team.
