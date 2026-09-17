# Docker Architecture Diagram

## Production Mode (`docker-compose.yml`)

```
┌─────────────────────────────────────────────────────────────────┐
│                         Docker Host                              │
│                                                                  │
│  ┌────────────────────────────────────────────────────────┐    │
│  │             pos-network (Bridge)                       │    │
│  │                                                        │    │
│  │  ┌─────────────────┐  ┌─────────────────┐  ┌────────┴────┐ │
│  │  │   PostgreSQL    │  │      Redis      │  │     API      │ │
│  │  │  (postgres:16)  │  │   (redis:7)     │  │   (custom)   │ │
│  │  ├─────────────────┤  ├─────────────────┤  ├──────────────┤ │
│  │  │ Port: 5432      │  │ Port: 6379      │  │ Port: 8080   │ │
│  │  │ Health: ✓       │  │ Health: ✓       │  │ Health: ✓    │ │
│  │  │ User: postgres  │  │ AOF: enabled    │  │ User: 1000   │ │
│  │  └────────┬────────┘  └────────┬────────┘  └──────┬───────┘ │
│  │           │                    │                   │         │
│  │           │                    │                   │         │
│  │           ▼                    ▼                   ▼         │
│  │     ┌──────────┐         ┌──────────┐      ┌──────────┐    │
│  │     │  pgdata  │         │redisdata │      │ uploads/ │    │
│  │     │ (volume) │         │ (volume) │      │  (bind)  │    │
│  │     └──────────┘         └──────────┘      └──────────┘    │
│  └────────────────────────────────────────────────────────┘    │
│                                                                 │
│  Exposed Ports:                                                 │
│  ├─ localhost:5432  → postgres:5432                            │
│  ├─ localhost:6379  → redis:6379                               │
│  └─ localhost:8080  → api:8080                                 │
└─────────────────────────────────────────────────────────────────┘
```

## Development Mode (`docker-compose.dev.yml`)

```
┌─────────────────────────────────────────────────────────────────┐
│                         Docker Host                              │
│                                                                  │
│  ┌────────────────────────────────────────────────────────┐    │
│  │           pos-network-dev (Bridge)                     │    │
│  │                                                        │    │
│  │  ┌─────────────────┐  ┌─────────────────┐  ┌────────┴────┐ │
│  │  │   PostgreSQL    │  │      Redis      │  │   API-DEV    │ │
│  │  │  (postgres:16)  │  │   (redis:7)     │  │  + Air 🔥    │ │
│  │  ├─────────────────┤  ├─────────────────┤  ├──────────────┤ │
│  │  │ DB: pos_fnb_dev │  │ Port: 6379      │  │ Port: 8080   │ │
│  │  │ Health: ✓       │  │ Health: ✓       │  │ Hot Reload ✓ │ │
│  │  └────────┬────────┘  └────────┬────────┘  └──────┬───────┘ │
│  │           │                    │                   │         │
│  │           │                    │                   ▼         │
│  │           ▼                    ▼           ┌────────────────┐│
│  │     ┌──────────┐         ┌──────────┐     │  Source Code   ││
│  │     │ pgdata_  │         │redisdata_│     │   (mounted)    ││
│  │     │   dev    │         │   dev    │     │  + uploads/    ││
│  │     │(volume)  │         │ (volume) │     │  + go_modules  ││
│  │     └──────────┘         └──────────┘     └────────────────┘│
│  └────────────────────────────────────────────────────────┘    │
│                                                                 │
│  Features:                                                      │
│  ✅ Hot Reload dengan Air                                      │
│  ✅ Source code changes = auto rebuild                         │
│  ✅ Go modules cached                                          │
│  ✅ Separate dev database                                      │
└─────────────────────────────────────────────────────────────────┘
```

## Dockerfile Build Process (Multi-Stage)

```
┌───────────────────────────────────────────────────────────┐
│                 STAGE 1: Builder                           │
│                                                            │
│  Base: golang:1.22-alpine                                  │
│                                                            │
│  Steps:                                                    │
│  1. COPY go.mod go.sum → /app/                            │
│  2. RUN go mod download  (cached if no changes)           │
│  3. COPY source code → /app/                              │
│  4. RUN go build (CGO_ENABLED=0, optimized)               │
│                                                            │
│  Output: /pos-fnb-api (static binary ~15MB)               │
└───────────────────┬───────────────────────────────────────┘
                    │
                    │ COPY binary only
                    ▼
┌───────────────────────────────────────────────────────────┐
│                 STAGE 2: Runtime                           │
│                                                            │
│  Base: alpine:3.19 (minimal ~5MB)                          │
│                                                            │
│  Steps:                                                    │
│  1. Install ca-certificates, tzdata                        │
│  2. Create non-root user (appuser:1000)                    │
│  3. COPY --from=builder /pos-fnb-api                       │
│  4. Set USER appuser                                       │
│                                                            │
│  Output: Final image ~20MB                                 │
│          (vs ~1GB with full golang image)                  │
└───────────────────────────────────────────────────────────┘
```

## Application Flow

```
┌─────────────┐
│   Client    │
│  (Browser)  │
└──────┬──────┘
       │ HTTP Request
       │ :8080
       ▼
┌──────────────────┐
│   API Container  │
│                  │
│  ┌────────────┐  │
│  │  Router    │  │
│  │  (chi)     │  │
│  └─────┬──────┘  │
│        │         │
│  ┌─────▼──────┐  │
│  │ Middleware │  │
│  │  - CORS    │  │
│  │  - JWT     │  │
│  └─────┬──────┘  │
│        │         │
│  ┌─────▼──────┐  │
│  │  Handler   │  │
│  └─────┬──────┘  │
│        │         │
│  ┌─────▼──────┐  │        ┌──────────────┐
│  │   GORM     │──┼───────▶│  PostgreSQL  │
│  └────────────┘  │        │   :5432      │
│                  │        └──────────────┘
│  ┌────────────┐  │        ┌──────────────┐
│  │   Redis    │──┼───────▶│    Redis     │
│  │   Client   │  │        │   :6379      │
│  └────────────┘  │        └──────────────┘
└──────────────────┘
```

## Health Check Flow

```
                    Every 30s
                        │
                        ▼
         ┌─────────────────────────┐
         │  Docker Health Check    │
         │  wget /health           │
         └────────┬────────────────┘
                  │
         ┌────────▼────────┐
         │  Response OK?   │
         └────────┬────────┘
                  │
        ┌─────────┴─────────┐
        │                   │
    ✅ YES               ❌ NO
        │                   │
        ▼                   ▼
   ┌─────────┐        ┌──────────┐
   │ healthy │        │unhealthy │
   └─────────┘        └────┬─────┘
                           │
                      After 3 fails
                           │
                           ▼
                     ┌──────────┐
                     │ Restart  │
                     │Container │
                     └──────────┘
```

## Network Communication

```
┌────────────────────────────────────────────────────┐
│              pos-network (Bridge)                   │
│                                                     │
│  Container DNS:                                     │
│  ├─ postgres  → 172.20.0.2:5432                    │
│  ├─ redis     → 172.20.0.3:6379                    │
│  └─ api       → 172.20.0.4:8080                    │
│                                                     │
│  Internal Communication:                            │
│  ┌─────┐                                            │
│  │ API │ --DB_HOST=postgres--> ┌──────────┐        │
│  │     │                        │PostgreSQL│        │
│  │     │ <--5432-------------- └──────────┘        │
│  │     │                                            │
│  │     │ --REDIS_HOST=redis--> ┌──────────┐        │
│  │     │                        │  Redis   │        │
│  │     │ <--6379-------------- └──────────┘        │
│  └─────┘                                            │
└────────────────────────────────────────────────────┘
         │
         │ Port Mapping
         ▼
┌────────────────────────────────────────────────────┐
│              Host Network                           │
│                                                     │
│  localhost:8080  → api:8080                        │
│  localhost:5432  → postgres:5432 (dev only)        │
│  localhost:6379  → redis:6379    (dev only)        │
└────────────────────────────────────────────────────┘
```

## Volume Persistence

```
Docker Volumes (Managed by Docker)
├─ pgdata/               # PostgreSQL data
│  ├─ base/              # Database files
│  ├─ global/            # Cluster data
│  └─ pg_wal/            # Write-ahead logs
│
└─ redisdata/            # Redis persistence
   ├─ appendonly.aof     # AOF file
   └─ dump.rdb           # RDB snapshot

Host Bind Mounts (Direct mapping)
└─ ./uploads/            # Uploaded files
   └─ (user files)
```

## Helper Scripts Flow

```
┌──────────────┐
│ User Command │
└──────┬───────┘
       │
  ┌────▼─────┐
  │Windows?  │
  └────┬─────┘
       │
   ┌───┴───┐
   │       │
 YES      NO
   │       │
   ▼       ▼
┌─────────────┐   ┌─────────────┐
│ docker.ps1  │   │  Makefile   │
│ (PowerShell)│   │   (Make)    │
└──────┬──────┘   └──────┬──────┘
       │                 │
       └────────┬────────┘
                │
                ▼
       ┌─────────────────┐
       │ Docker Compose  │
       └────────┬────────┘
                │
       ┌────────▼─────────┐
       │  Start Services  │
       │  - PostgreSQL    │
       │  - Redis         │
       │  - API           │
       └──────────────────┘
```

## CI/CD Flow (Future)

```
┌──────────┐
│   Git    │
│   Push   │
└────┬─────┘
     │
     ▼
┌──────────────┐
│  CI Pipeline │
│  (GitHub     │
│   Actions)   │
└─────┬────────┘
      │
      ├─► Run Tests
      │   └─> go test ./...
      │
      ├─► Build Image
      │   └─> docker build
      │
      ├─► Push to Registry
      │   └─> Docker Hub / ECR
      │
      └─► Deploy
          └─> AWS ECS / K8s
```
