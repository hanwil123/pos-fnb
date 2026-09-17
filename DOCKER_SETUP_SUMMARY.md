# 📦 Docker Setup Summary

Dokumentasi lengkap untuk setup Docker yang telah dibuat untuk POS FnB Backend.

## ✅ File yang Dibuat/Diupdate

### Core Docker Files

1. **`Dockerfile`** - Production optimized multi-stage build
   - Stage 1: Builder (compile Go binary)
   - Stage 2: Runtime (Alpine dengan binary)
   - Security: Non-root user
   - Health check included
   - Size: ~20MB

2. **`docker-compose.yml`** - Production setup
   - PostgreSQL 16
   - Redis 7
   - API service
   - Health checks pada semua services
   - Volume persistence
   - Network isolation

3. **`Dockerfile.dev`** - Development build dengan Air hot reload
   - Langsung compile dan reload saat file berubah
   - Ideal untuk development

4. **`docker-compose.dev.yml`** - Development setup
   - Hot reload dengan Air
   - Source code mounted
   - Development database terpisah
   - Go modules cache

### Configuration Files

5. **`.dockerignore`** - Mengoptimalkan Docker build
   - Exclude file yang tidak perlu
   - Mempercepat build time

6. **`.air.toml`** - Konfigurasi Air hot reload
   - Auto rebuild on `.go` file changes
   - Exclude test files
   - Custom build commands

7. **`.gitignore`** - Git ignore rules
   - Ignore binaries
   - Ignore .env files
   - Ignore uploads

8. **`.env.example`** - Template environment variables
   - Semua variable yang diperlukan
   - Nilai default untuk development

### Helper Scripts

9. **`docker.ps1`** - PowerShell helper untuk Windows
   - Commands: up, down, restart, logs, clean, test, shell, db-shell
   - Development mode: dev-up, dev-down, dev-logs
   - Colored output & user friendly

10. **`Makefile`** - Make commands untuk Linux/Mac/Git Bash
    - Same commands as docker.ps1
    - POSIX compatible

### Documentation

11. **`DOCKER.md`** - Comprehensive Docker documentation
    - Setup guide
    - Commands reference
    - Troubleshooting
    - Best practices
    - Security guidelines

12. **`QUICKSTART.md`** - Quick start guide
    - 3-step quick start
    - Platform-specific commands
    - Testing instructions
    - Common troubleshooting

13. **`README.md`** - Updated dengan Docker info

14. **`uploads/`** - Directory untuk file uploads
    - Mounted sebagai volume
    - `.gitkeep` untuk Git tracking

---

## 🎯 Features

### Production Mode

**Command:**
```bash
docker compose up --build
```

**Features:**
- ✅ Multi-stage optimized build
- ✅ Small image size (~20MB)
- ✅ Security hardened (non-root user)
- ✅ Health checks
- ✅ Auto migration
- ✅ Volume persistence
- ✅ Network isolation

**Services:**
- PostgreSQL 16 (port 5432)
- Redis 7 (port 6379)
- API (port 8080)

### Development Mode

**Command:**
```bash
docker compose -f docker-compose.dev.yml up
# Atau
.\docker.ps1 dev-up
```

**Features:**
- ✅ Hot reload dengan Air
- ✅ Source code mounted (perubahan langsung terdeteksi)
- ✅ Development database terpisah
- ✅ Go modules cache (faster rebuild)
- ✅ Debug logging enabled
- ✅ Auto migration & seeding

---

## 📂 Project Structure

```
pos-fnb-backend/
├── cmd/
│   └── api/              # Application entrypoint
├── internal/             # Internal packages
├── uploads/              # Upload directory (mounted)
│   └── .gitkeep
├── tmp/                  # Air build directory (auto-created)
│
├── Dockerfile            # Production build
├── Dockerfile.dev        # Development build
├── docker-compose.yml    # Production compose
├── docker-compose.dev.yml # Development compose
├── .dockerignore         # Docker ignore rules
├── .air.toml             # Air hot reload config
├── .gitignore            # Git ignore rules
├── .env                  # Environment variables (gitignored)
├── .env.example          # Environment template
│
├── docker.ps1            # Windows helper script
├── Makefile              # Linux/Mac helper
│
├── DOCKER.md             # Docker documentation
├── QUICKSTART.md         # Quick start guide
├── README.md             # Project readme
└── ARCHITECTURE.md       # Architecture docs
```

---

## 🚀 Usage Examples

### Scenario 1: First Time Setup

```bash
# 1. Clone & setup
git clone <repo>
cd pos-fnb-backend
cp .env.example .env

# 2. Start services
docker compose up --build

# 3. Test
curl http://localhost:8080/health
```

### Scenario 2: Development Workflow

```bash
# Terminal 1: Start dev mode
.\docker.ps1 dev-up

# Terminal 2: Watch logs
.\docker.ps1 dev-logs

# Edit code di editor favorit
# → Air akan auto-reload
# → Cek logs di Terminal 2
```

### Scenario 3: Testing Production Build

```bash
# Build production image
docker compose build

# Run production mode
docker compose up

# Test endpoints
curl http://localhost:8080/api/v1/menu
```

### Scenario 4: Database Management

```bash
# Access PostgreSQL
.\docker.ps1 db-shell

# Backup database
docker compose exec postgres pg_dump -U postgres pos_fnb > backup.sql

# Restore database
docker compose exec -T postgres psql -U postgres -d pos_fnb < backup.sql
```

### Scenario 5: Cleanup

```bash
# Soft cleanup (stop services)
docker compose down

# Hard cleanup (remove volumes)
.\docker.ps1 clean
```

---

## 🔧 Configuration

### Environment Variables

**Production:**
```env
APP_ENV=production
DB_SSLMODE=require
JWT_SECRET=<strong-random-secret>
QR_HMAC_SECRET=<strong-random-secret>
DB_PASS=<strong-password>
```

**Development:**
```env
APP_ENV=development
DB_SSLMODE=disable
DB_AUTO_MIGRATE=true
DB_SEED=true
LOG_LEVEL=debug
```

### Port Configuration

Default ports (dapat diubah di `.env`):
- API: 8080
- PostgreSQL: 5432
- Redis: 6379

---

## 🔒 Security Best Practices

1. **Never commit `.env`** - Already in `.gitignore`
2. **Change default secrets** - Especially `JWT_SECRET`, `QR_HMAC_SECRET`
3. **Use strong passwords** - For database in production
4. **Enable SSL** - Set `DB_SSLMODE=require` in production
5. **Limit exposed ports** - Don't expose DB ports in production
6. **Regular updates** - Keep base images updated

---

## 📊 Performance Optimizations

### Docker Build

- ✅ Multi-stage build (smaller image)
- ✅ Layer caching (faster rebuilds)
- ✅ `.dockerignore` (exclude unnecessary files)
- ✅ CGO disabled (static binary)
- ✅ Build flags: `-ldflags="-w -s"` (smaller binary)

### Runtime

- ✅ Alpine base image (minimal footprint)
- ✅ Health checks (automatic recovery)
- ✅ Resource limits (dapat dikonfigurasi)
- ✅ Volume mounts (persistence)

### Development

- ✅ Go modules cache (faster dev builds)
- ✅ Air hot reload (instant feedback)
- ✅ Source mounted (no rebuild needed)

---

## 🧪 Testing

### Manual Testing

```bash
# Health check
curl http://localhost:8080/health

# Menu endpoint
curl http://localhost:8080/api/v1/menu?restaurant_id=xxx

# Check database
docker compose exec postgres psql -U postgres -d pos_fnb -c "SELECT * FROM restaurants;"
```

### Automated Testing

```bash
# Run tests in container
docker compose exec api go test -v ./...

# Or using helper
.\docker.ps1 test
```

---

## 🐛 Common Issues & Solutions

### Issue: Port already in use

**Solution:**
```bash
# Edit .env
APP_PORT=8081
DB_PORT=5433
```

### Issue: Permission denied (Windows)

**Solution:**
```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### Issue: Container keeps restarting

**Solution:**
```bash
# Check logs
docker compose logs api

# Common causes:
# - Database not ready (health check will handle this)
# - Wrong environment variables
# - Port conflicts
```

### Issue: Changes not reflected (dev mode)

**Solution:**
```bash
# Make sure using dev compose
docker compose -f docker-compose.dev.yml up

# Not regular compose
docker compose up  # ❌ Wrong for dev
```

---

## 🎓 Learning Resources

### Docker
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Dockerfile Best Practices](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/)
- [Multi-stage Builds](https://docs.docker.com/build/building/multi-stage/)

### Go + Docker
- [Go Docker Hub](https://hub.docker.com/_/golang)
- [Containerizing Go Applications](https://docs.docker.com/language/golang/)

### Hot Reload
- [Air - Hot reload for Go](https://github.com/cosmtrek/air)

---

## 📝 Maintenance

### Regular Tasks

**Weekly:**
- Check logs: `docker compose logs`
- Monitor disk usage: `docker system df`

**Monthly:**
- Update base images: `docker compose pull`
- Rebuild: `docker compose build --no-cache`
- Cleanup: `docker system prune`

**Before Production:**
- Review `.env` (change all secrets)
- Enable SSL for database
- Remove dev ports exposure
- Test backup/restore procedure

---

## 🎉 Summary

Setup Docker yang telah dibuat mencakup:

✅ **Production-ready** dengan optimized build
✅ **Development-friendly** dengan hot reload
✅ **Well-documented** dengan guides lengkap
✅ **Cross-platform** (Windows, Linux, Mac)
✅ **Helper scripts** untuk kemudahan penggunaan
✅ **Best practices** untuk security & performance
✅ **Health checks** untuk reliability
✅ **Volume persistence** untuk data
✅ **Network isolation** untuk security

**Next Steps:**
1. Review `.env` dan sesuaikan dengan kebutuhan
2. Jalankan `docker compose up --build`
3. Test API endpoints
4. Mulai development dengan `dev-up`

Happy coding! 🚀
