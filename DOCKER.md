# 🐳 Docker Setup Guide

Panduan lengkap untuk menjalankan POS FnB Backend menggunakan Docker.

## 📋 Prerequisites

Pastikan sudah terinstall:
- **Docker Desktop** (untuk Windows/Mac) atau **Docker Engine** (untuk Linux)
- **Docker Compose** v2.x atau lebih baru

Cek versi Docker:
```bash
docker --version
docker compose version
```

## 🚀 Quick Start

### 1. Setup Environment Variables

Copy file `.env.example` ke `.env`:
```bash
cp .env.example .env
```

Edit `.env` sesuai kebutuhan. Untuk development, nilai default sudah cukup.

### 2. Jalankan Aplikasi

**Menggunakan Docker Compose langsung:**
```bash
docker compose up --build
```

**Atau menggunakan helper script (Windows):**
```powershell
.\docker.ps1 up
```

**Atau menggunakan Makefile (Linux/Mac/Git Bash):**
```bash
make up
```

### 3. Cek Status

Setelah semua service berjalan, aplikasi dapat diakses di:
- **API**: http://localhost:8080
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379

Test health check:
```bash
curl http://localhost:8080/health
```

## 📦 Services

Docker Compose akan menjalankan 3 services:

| Service | Image | Port | Purpose |
|---------|-------|------|---------|
| **postgres** | postgres:16-alpine | 5432 | Database utama |
| **redis** | redis:7-alpine | 6379 | Caching & sessions |
| **api** | (custom build) | 8080 | REST API Backend |

## 🛠️ Command Reference

### Helper Script (Windows - PowerShell)

```powershell
# Lihat semua perintah
.\docker.ps1 help

# Build images
.\docker.ps1 build

# Start services (detached mode)
.\docker.ps1 up

# Stop services
.\docker.ps1 down

# Restart services
.\docker.ps1 restart

# Lihat logs
.\docker.ps1 logs

# Cleanup (hapus volumes)
.\docker.ps1 clean

# Run tests
.\docker.ps1 test

# Masuk ke API container shell
.\docker.ps1 shell

# Masuk ke PostgreSQL shell
.\docker.ps1 db-shell

# Development mode (rebuild + foreground)
.\docker.ps1 dev
```

### Makefile (Linux/Mac/Git Bash)

```bash
make help      # Lihat semua perintah
make build     # Build images
make up        # Start services
make down      # Stop services
make restart   # Restart services
make logs      # Lihat logs
make clean     # Cleanup (hapus volumes)
make test      # Run tests
make shell     # Masuk ke API container shell
make db-shell  # Masuk ke PostgreSQL shell
make dev       # Development mode
```

### Docker Compose Native

```bash
# Build images
docker compose build

# Start services (detached)
docker compose up -d

# Start services (foreground)
docker compose up

# Stop services
docker compose down

# Restart specific service
docker compose restart api

# Lihat logs
docker compose logs -f

# Lihat logs service tertentu
docker compose logs -f api

# Stop dan hapus volumes
docker compose down -v

# Exec command di container
docker compose exec api sh
docker compose exec postgres psql -U postgres -d pos_fnb
```

## 🔍 Troubleshooting

### Port sudah digunakan

Jika port 5432, 6379, atau 8080 sudah digunakan, edit file `.env`:
```env
DB_PORT=5433
APP_PORT=8081
```

Dan update `docker-compose.yml` untuk mapping port yang sesuai, atau stop service yang menggunakan port tersebut.

### Database connection error

Tunggu beberapa detik setelah `docker compose up` karena PostgreSQL perlu waktu untuk initialize. 

Cek health status:
```bash
docker compose ps
```

Semua service harus dalam status `healthy`.

### Rebuild setelah perubahan code

```bash
docker compose down
docker compose build --no-cache api
docker compose up -d
```

Atau gunakan:
```powershell
.\docker.ps1 dev
```

### Lihat logs untuk debugging

```bash
# Semua logs
docker compose logs -f

# Logs API saja
docker compose logs -f api

# Logs PostgreSQL saja
docker compose logs -f postgres
```

### Reset database

```bash
# Stop dan hapus volumes
docker compose down -v

# Start ulang
docker compose up -d
```

⚠️ **Warning**: Ini akan menghapus semua data di database!

### Container tidak bisa connect ke database

1. Pastikan container `postgres` dalam status `healthy`:
   ```bash
   docker compose ps
   ```

2. Cek logs PostgreSQL:
   ```bash
   docker compose logs postgres
   ```

3. Pastikan environment variable `DB_HOST` di service `api` adalah `postgres` (bukan `localhost`).

## 📝 Development Workflow

### 1. Local Development dengan Hot Reload

Untuk development dengan hot reload (memerlukan setup tambahan dengan `air` atau `reflex`), jalankan API di luar Docker:

```bash
# Start hanya database
docker compose up -d postgres redis

# Run API locally
go run ./cmd/api
```

Ubah `.env` menjadi:
```env
DB_HOST=localhost
REDIS_HOST=localhost:6379
```

### 2. Testing di Container

```bash
# Run tests
docker compose exec api go test -v ./...

# Run specific test
docker compose exec api go test -v ./internal/handler -run TestCreateOrder
```

### 3. Database Management

**Masuk ke PostgreSQL shell:**
```bash
docker compose exec postgres psql -U postgres -d pos_fnb
```

**Backup database:**
```bash
docker compose exec postgres pg_dump -U postgres pos_fnb > backup.sql
```

**Restore database:**
```bash
docker compose exec -T postgres psql -U postgres -d pos_fnb < backup.sql
```

## 🏗️ Build Process

Dockerfile menggunakan **multi-stage build**:

1. **Builder stage**: Compile Go application dengan CGO disabled
2. **Run stage**: Alpine Linux minimal dengan binary yang sudah di-build

Benefits:
- Image size lebih kecil (~20MB vs ~1GB)
- Build caching untuk dependencies
- Security: running as non-root user

## 🔐 Security Best Practices

1. **Jangan commit `.env`** - Sudah ditambahkan ke `.gitignore`
2. **Ganti secrets di production** - Terutama `JWT_SECRET` dan `QR_HMAC_SECRET`
3. **Use strong PostgreSQL password** di production
4. **Enable SSL** untuk PostgreSQL di production (`DB_SSLMODE=require`)
5. **Limit network exposure** - Jangan expose database port di production

## 📊 Monitoring

### Health Checks

Docker Compose sudah dikonfigurasi dengan health checks:
- PostgreSQL: `pg_isready`
- Redis: `redis-cli ping`
- API: `GET /health`

Cek status:
```bash
docker compose ps
```

### Resource Usage

```bash
# Lihat resource usage
docker stats

# Lihat disk usage
docker system df
```

## 🧹 Cleanup

### Cleanup project containers
```bash
docker compose down -v
```

### Cleanup all unused Docker resources
```bash
docker system prune -a
```

## 📚 Additional Resources

- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Go Docker Best Practices](https://docs.docker.com/language/golang/build-images/)
- [PostgreSQL Docker Hub](https://hub.docker.com/_/postgres)
- [Redis Docker Hub](https://hub.docker.com/_/redis)

## 🆘 Getting Help

Jika mengalami masalah:
1. Cek logs: `docker compose logs -f`
2. Cek container status: `docker compose ps`
3. Restart services: `docker compose restart`
4. Rebuild dari scratch: `docker compose down -v && docker compose up --build`
