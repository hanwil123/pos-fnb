# 🐳 Docker Cheat Sheet - POS FnB Backend

Quick reference untuk commands yang sering digunakan.

## 🚀 Quick Start

```bash
# Setup & Start (First Time)
cp .env.example .env
docker compose up --build

# Test
curl http://localhost:8080/health
```

---

## 📋 Production Mode

### Windows (PowerShell)
```powershell
.\docker.ps1 up        # Start
.\docker.ps1 down      # Stop
.\docker.ps1 restart   # Restart
.\docker.ps1 logs      # View logs
.\docker.ps1 build     # Rebuild
.\docker.ps1 clean     # Remove all (with volumes)
```

### Linux/Mac/Git Bash
```bash
make up        # Start
make down      # Stop
make restart   # Restart
make logs      # View logs
make build     # Rebuild
make clean     # Remove all (with volumes)
```

### Docker Compose Native
```bash
docker compose up -d              # Start (detached)
docker compose down               # Stop
docker compose restart            # Restart
docker compose logs -f            # Follow logs
docker compose build              # Build images
docker compose down -v            # Remove with volumes
```

---

## 👨‍💻 Development Mode (Hot Reload)

### Windows
```powershell
.\docker.ps1 dev-up      # Start dev mode
.\docker.ps1 dev-down    # Stop dev mode
.\docker.ps1 dev-logs    # View dev logs
```

### Linux/Mac
```bash
make dev-up      # Start dev mode
make dev-down    # Stop dev mode
make dev-logs    # View dev logs
```

### Native
```bash
docker compose -f docker-compose.dev.yml up
docker compose -f docker-compose.dev.yml down
docker compose -f docker-compose.dev.yml logs -f
```

---

## 🔍 Inspection & Debugging

### Check Status
```bash
docker compose ps                    # List containers
docker compose ps --format json      # JSON format
docker stats                         # Resource usage
```

### View Logs
```bash
docker compose logs                  # All logs
docker compose logs -f               # Follow logs
docker compose logs -f api           # API logs only
docker compose logs --tail=100 api   # Last 100 lines
```

### Container Shell Access
```bash
# API container
docker compose exec api sh
.\docker.ps1 shell

# PostgreSQL
docker compose exec postgres psql -U postgres -d pos_fnb
.\docker.ps1 db-shell

# Redis
docker compose exec redis redis-cli
```

---

## 🗄️ Database Operations

### PostgreSQL Commands
```bash
# Access psql
docker compose exec postgres psql -U postgres -d pos_fnb

# Run SQL file
docker compose exec -T postgres psql -U postgres -d pos_fnb < script.sql

# Backup
docker compose exec postgres pg_dump -U postgres pos_fnb > backup.sql
docker compose exec postgres pg_dump -U postgres pos_fnb | gzip > backup.sql.gz

# Restore
docker compose exec -T postgres psql -U postgres -d pos_fnb < backup.sql
gunzip -c backup.sql.gz | docker compose exec -T postgres psql -U postgres -d pos_fnb

# List databases
docker compose exec postgres psql -U postgres -c "\l"

# List tables
docker compose exec postgres psql -U postgres -d pos_fnb -c "\dt"

# Check database size
docker compose exec postgres psql -U postgres -d pos_fnb -c "SELECT pg_size_pretty(pg_database_size('pos_fnb'));"
```

### Redis Commands
```bash
# Access redis-cli
docker compose exec redis redis-cli

# Inside redis-cli:
PING                    # Test connection
KEYS *                  # List all keys
GET key                 # Get value
SET key value           # Set value
DEL key                 # Delete key
FLUSHALL                # Clear all data (⚠️ careful!)
INFO                    # Server info
DBSIZE                  # Number of keys
```

---

## 🔨 Build & Rebuild

### Full Rebuild
```bash
# Rebuild everything from scratch
docker compose down -v
docker compose build --no-cache
docker compose up -d

# Or using helper
.\docker.ps1 clean
.\docker.ps1 build
.\docker.ps1 up
```

### Rebuild Specific Service
```bash
docker compose build api
docker compose up -d --no-deps api
```

### Pull Latest Images
```bash
docker compose pull
docker compose up -d
```

---

## 🧪 Testing

### Run Tests
```bash
# In container
docker compose exec api go test -v ./...

# Specific package
docker compose exec api go test -v ./internal/handler

# With coverage
docker compose exec api go test -v -cover ./...

# Or using helper
.\docker.ps1 test
```

### Manual API Testing
```bash
# Health check
curl http://localhost:8080/health

# Menu endpoint (example)
curl http://localhost:8080/api/v1/menu?restaurant_id=xxx

# With authentication
curl -H "Authorization: Bearer <token>" http://localhost:8080/api/v1/staff/tables
```

---

## 🧹 Cleanup

### Remove Containers
```bash
docker compose down                # Stop & remove containers
docker compose down -v             # Also remove volumes
docker compose rm -f               # Force remove stopped containers
```

### System Cleanup
```bash
# Remove unused images
docker image prune -a

# Remove unused volumes
docker volume prune

# Remove unused networks
docker network prune

# Clean everything (⚠️ nuclear option)
docker system prune -a --volumes
```

### Project Cleanup
```bash
# Helper scripts
.\docker.ps1 clean         # Windows
make clean                 # Linux/Mac
```

---

## 📊 Monitoring

### Resource Usage
```bash
docker stats                       # All containers
docker stats pos-fnb-api           # Specific container
docker stats --no-stream           # One-time snapshot
```

### Disk Usage
```bash
docker system df                   # Overview
docker system df -v                # Verbose
```

### Container Processes
```bash
docker compose top                 # All services
docker compose top api             # API service only
```

---

## 🔧 Environment Variables

### List Current Config
```bash
docker compose config              # Show full config
docker compose config --services   # List services
docker compose config --volumes    # List volumes
```

### Override Variables
```bash
# Temporary override
APP_PORT=8081 docker compose up

# Use different .env file
docker compose --env-file .env.prod up
```

---

## 🌐 Network Operations

### Inspect Network
```bash
docker network ls
docker network inspect pos-fnb-backend_pos-network
```

### Test Connectivity
```bash
# From API to PostgreSQL
docker compose exec api wget -O- http://postgres:5432

# Ping test (if ping installed)
docker compose exec api ping postgres
```

---

## 🔄 Update & Maintenance

### Update Application
```bash
# Pull latest code
git pull

# Rebuild and restart
docker compose down
docker compose build
docker compose up -d
```

### Update Base Images
```bash
# Pull latest base images
docker compose pull

# Rebuild with new base
docker compose build --pull
docker compose up -d
```

### Check for Updates
```bash
# Check image digests
docker compose images
docker image inspect postgres:16-alpine
```

---

## 🚨 Troubleshooting

### Container Won't Start
```bash
# Check logs
docker compose logs api

# Check last exit code
docker compose ps

# Try running interactively
docker compose run --rm api sh
```

### Port Already in Use
```bash
# Find process using port (Windows)
netstat -ano | findstr :8080

# Find process using port (Linux/Mac)
lsof -i :8080

# Kill process (Windows)
taskkill /PID <pid> /F

# Kill process (Linux/Mac)
kill -9 <pid>
```

### Database Connection Issues
```bash
# Check postgres is healthy
docker compose ps postgres

# Check postgres logs
docker compose logs postgres

# Test connection
docker compose exec api nc -zv postgres 5432
```

### Reset Everything
```bash
# Nuclear option - start fresh
docker compose down -v
docker system prune -a
docker volume prune
cp .env.example .env
docker compose up --build
```

---

## 📝 Useful Aliases (Optional)

Add to your shell profile:

### Bash/Zsh (~/.bashrc or ~/.zshrc)
```bash
alias dcu='docker compose up -d'
alias dcd='docker compose down'
alias dcl='docker compose logs -f'
alias dcr='docker compose restart'
alias dcb='docker compose build'
alias dce='docker compose exec'
```

### PowerShell ($PROFILE)
```powershell
function dcu { docker compose up -d }
function dcd { docker compose down }
function dcl { docker compose logs -f }
function dcr { docker compose restart }
function dcb { docker compose build }
```

---

## 🎯 Common Workflows

### Morning Startup
```bash
docker compose up -d
docker compose logs -f api
# Ctrl+C when ready
```

### After Code Changes (Dev)
```bash
# Dev mode - auto reload, no action needed!
# Just save file and watch logs

# Production mode - manual rebuild
docker compose restart api
```

### Before Committing
```bash
docker compose exec api go test ./...
docker compose exec api go fmt ./...
```

### End of Day
```bash
docker compose down
# Or keep running:
# docker compose stop
```

---

## 📚 References

- Full Documentation: [DOCKER.md](./DOCKER.md)
- Quick Start: [QUICKSTART.md](./QUICKSTART.md)
- Architecture: [DOCKER_DIAGRAMS.md](./DOCKER_DIAGRAMS.md)
- Setup Summary: [DOCKER_SETUP_SUMMARY.md](./DOCKER_SETUP_SUMMARY.md)

---

**Pro Tip:** Bookmark this file for quick reference! 📌
