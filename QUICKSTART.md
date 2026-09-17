# 🚀 Quick Start Guide

Panduan cepat untuk menjalankan POS FnB Backend dalam 3 langkah.

## ⚡ Super Quick (Docker)

```bash
# 1. Copy environment file
cp .env.example .env

# 2. Start semua services
docker compose up --build

# 3. Test
curl http://localhost:8080/health
```

Selesai! API berjalan di **http://localhost:8080**

---

## 📖 Detailed Steps

### Windows (PowerShell)

```powershell
# Setup
cp .env.example .env

# Production Mode
.\docker.ps1 up

# Development Mode (dengan hot reload)
.\docker.ps1 dev-up

# Lihat logs
.\docker.ps1 logs

# Stop services
.\docker.ps1 down
```

### Linux / Mac / Git Bash

```bash
# Setup
cp .env.example .env

# Production Mode
make up

# Development Mode (dengan hot reload)
make dev-up

# Lihat logs
make logs

# Stop services
make down
```

---

## 🧪 Testing API

### Health Check

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "ok",
  "timestamp": "2024-..."
}
```

### Cek Services

```bash
# Windows
.\docker.ps1 shell
docker compose ps

# Linux/Mac
make shell
docker compose ps
```

Semua services harus berstatus **healthy**.

---

## 🔧 Development vs Production

| Feature | Production (`docker compose up`) | Development (`dev-up`) |
|---------|----------------------------------|----------------------|
| **Hot Reload** | ❌ Tidak | ✅ Ya (Air) |
| **Source Mounted** | ❌ Tidak | ✅ Ya |
| **Build Time** | Cepat (binary) | Lebih lambat |
| **Use Case** | Testing production build | Active development |

### Kapan Menggunakan Apa?

**Production Mode** (`docker compose up`):
- Testing final build
- Deployment simulation
- Performance testing
- Production-like environment

**Development Mode** (`dev-up`):
- Active coding
- Hot reload saat save file
- Quick iteration
- Debugging

---

## 🗄️ Database Access

### PostgreSQL

```bash
# Windows
.\docker.ps1 db-shell

# Linux/Mac
make db-shell

# Atau langsung
docker compose exec postgres psql -U postgres -d pos_fnb
```

### Redis

```bash
docker compose exec redis redis-cli
```

---

## 📝 Environment Variables

File `.env` sudah berisi nilai default untuk development. 

**Wajib diganti di production:**
- `JWT_SECRET` - Secret untuk JWT token
- `QR_HMAC_SECRET` - Secret untuk QR code signing
- `DB_PASS` - Password PostgreSQL
- `SMTP_USER` & `SMTP_PASS` - Jika menggunakan email

---

## 🐛 Troubleshooting

### Port sudah digunakan

```bash
# Check port usage
netstat -ano | findstr :8080    # Windows
lsof -i :8080                   # Linux/Mac

# Solusi: Edit .env
APP_PORT=8081
```

### Container tidak start

```bash
# Lihat logs
docker compose logs

# Restart dari awal
docker compose down -v
docker compose up --build
```

### Hot reload tidak bekerja

Pastikan menggunakan development mode:
```bash
# Bukan ini
docker compose up

# Tapi ini
docker compose -f docker-compose.dev.yml up
# Atau
.\docker.ps1 dev-up
```

---

## 📚 Next Steps

1. **Baca [DOCKER.md](./DOCKER.md)** untuk dokumentasi lengkap
2. **Baca [README.md](./README.md)** untuk API endpoints
3. **Baca [ARCHITECTURE.md](./ARCHITECTURE.md)** untuk struktur code

---

## 💡 Tips

### Rebuild setelah perubahan dependencies

```bash
docker compose down
docker compose build --no-cache
docker compose up
```

### Cleanup volumes (hapus data)

```bash
# Windows
.\docker.ps1 clean

# Linux/Mac
make clean
```

### Development workflow yang efisien

```bash
# Terminal 1: Run development mode
.\docker.ps1 dev-up

# Terminal 2: Watch logs
.\docker.ps1 dev-logs

# Edit code → Auto reload → Check logs → Repeat
```

---

## 🆘 Butuh Bantuan?

- **Docker issues**: Lihat [DOCKER.md](./DOCKER.md) troubleshooting section
- **API issues**: Lihat [README.md](./README.md) testing section
- **Architecture questions**: Lihat [ARCHITECTURE.md](./ARCHITECTURE.md)

---

**Happy Coding! 🎉**
