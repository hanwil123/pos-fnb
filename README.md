# POS FnB System — Backend (Golang)

Backend REST API untuk sistem POS FnB berbasis QR code per meja.
Router: [chi](https://github.com/go-chi/chi) · ORM: [GORM](https://gorm.io) · DB: PostgreSQL.

## Menjalankan dengan Docker Compose (Rekomendasi! ✅)

**Cara termudah untuk menjalankan aplikasi:**

```bash
# 1. Copy environment variables
cp .env.example .env

# 2. Start semua services (PostgreSQL + Redis + API)
docker compose up --build
```

**Atau gunakan helper script (Windows):**
```powershell
.\docker.ps1 up
```

**Atau gunakan Makefile (Linux/Mac/Git Bash):**
```bash
make up
```

Server akan berjalan di `http://localhost:8080`. 

📖 **Dokumentasi lengkap Docker**: Lihat [DOCKER.md](./DOCKER.md)

📚 **Semua dokumentasi**: Lihat [INDEX.md](./INDEX.md) untuk navigasi lengkap

---

## Menjalankan secara lokal (tanpa Docker)

1. Pastikan PostgreSQL sudah berjalan dan buat database `pos_fnb`.
2. Copy environment variable:
   ```bash
   cp .env.example .env
   ```
   Sesuaikan `DB_HOST`, `DB_USER`, `DB_PASSWORD`, dll dengan setup lokal kamu.
3. Download dependency & jalankan:
   ```bash
   go mod tidy
   go run ./cmd/api
   ```
4. Server berjalan di `http://localhost:8080`. Auto-migration akan otomatis membuat semua tabel saat pertama kali start.
5. Cek health check:
   ```bash
   curl http://localhost:8080/health
   ```

## Menjalankan dengan Docker Compose (lebih mudah)

```bash
docker compose up --build
```
Ini akan menjalankan PostgreSQL, Redis, dan API sekaligus.

## Contoh Testing Manual Flow

**1. Buat restoran & meja dulu (lewat DB langsung atau tambahkan endpoint seed nanti):**
```sql
INSERT INTO restaurants (id, name, created_at, updated_at)
VALUES ('11111111-1111-1111-1111-111111111111', 'Kopi Senja', now(), now());
```

**2. Buat meja via API (butuh JWT staff — untuk testing awal bisa generate token manual dengan jwt.io menggunakan JWT_SECRET yang sama):**
```bash
curl -X POST http://localhost:8080/api/v1/staff/tables \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"restaurant_id":"11111111-1111-1111-1111-111111111111","table_number":"A1"}'
```

**3. Simulasikan customer scan QR (pakai qr_token dari response di atas):**
```bash
curl http://localhost:8080/api/v1/table/<qr_token>
```

**4. Lihat menu:**
```bash
curl "http://localhost:8080/api/v1/menu?restaurant_id=11111111-1111-1111-1111-111111111111"
```

**5. Buat order (payment_method: "cashier" untuk dapat QR bayar di kasir):**
```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "<session_id dari langkah 3>",
    "table_id": "<table id>",
    "customer_name": "Budi",
    "payment_method": "cashier",
    "items": [{"menu_item_id": "<id menu>", "quantity": 2}]
  }'
```
Response akan berisi `order_qr_payload` — ini yang di-encode jadi QR code di frontend, dan yang akan discan kasir lewat endpoint `/api/v1/staff/orders/scan-qr`.

## Struktur Folder

```
cmd/api/            entrypoint aplikasi
internal/
  config/            load environment variables
  database/          koneksi & auto-migration GORM
  models/            struct model = tabel database
  handler/           HTTP handler per resource (table, menu, order)
  middleware/         CORS & JWT auth staff
  router/            definisi semua route
  utils/             helper response JSON & signed QR token (HMAC)
```

## Yang Masih Perlu Ditambahkan (Next Steps)

- [ ] Endpoint login staff (`POST /api/v1/staff/login`) — saat ini JWT harus digenerate manual untuk testing.
- [ ] Integrasi payment gateway online (Midtrans/Xendit) di `payment_method: "online"`.
- [ ] WebSocket untuk update status order realtime ke customer & Kitchen Display System.
- [ ] Materialized view `menu_trending` + endpoint recommendations yang benar-benar dinamis.
- [ ] AI chat endpoint (`POST /api/v1/ai/chat`) dengan RAG memakai `pgvector`.
- [ ] Rate limiting & idempotency key pada endpoint order/payment.
- [ ] Unit test & integration test.
