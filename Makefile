.PHONY: help build up down restart logs clean test dev-up dev-down dev-logs

help: ## Tampilkan perintah yang tersedia
	@echo "Perintah Docker yang tersedia:"
	@echo ""
	@echo "Production:"
	@echo "  make build      - Build semua Docker images"
	@echo "  make up         - Jalankan semua services"
	@echo "  make down       - Stop semua services"
	@echo "  make restart    - Restart semua services"
	@echo "  make logs       - Lihat logs semua services"
	@echo ""
	@echo "Development (hot reload):"
	@echo "  make dev-up     - Start development mode dengan hot reload"
	@echo "  make dev-down   - Stop development mode"
	@echo "  make dev-logs   - Lihat logs development mode"
	@echo ""
	@echo "Utilities:"
	@echo "  make clean      - Stop services dan hapus volumes"
	@echo "  make test       - Jalankan test di dalam container"
	@echo "  make shell      - Masuk ke shell API container"
	@echo "  make db-shell   - Masuk ke PostgreSQL shell"

build: ## Build Docker images
	docker compose build

up: ## Start semua services
	docker compose up -d
	@echo "✅ Services started. API: http://localhost:8080"
	@echo "✅ PostgreSQL: localhost:5432"
	@echo "✅ Redis: localhost:6379"

down: ## Stop semua services
	docker compose down

restart: ## Restart semua services
	docker compose restart

logs: ## Lihat logs
	docker compose logs -f

clean: ## Stop dan hapus volumes
	docker compose down -v
	@echo "⚠️  Database dan Redis data telah dihapus!"

test: ## Run tests
	docker compose exec api go test -v ./...

shell: ## Masuk ke API container shell
	docker compose exec api sh

db-shell: ## Masuk ke PostgreSQL shell
	docker compose exec postgres psql -U postgres -d pos_fnb

db-migrate: ## Jalankan migration (auto migrate akan jalan saat startup)
	docker compose exec api ./pos-fnb-api

dev: ## Development mode (rebuild dan restart)
	docker compose down
	docker compose build
	docker compose up

dev-up: ## Start development mode dengan hot reload
	docker compose -f docker-compose.dev.yml up --build

dev-down: ## Stop development mode
	docker compose -f docker-compose.dev.yml down

dev-logs: ## Lihat logs development mode
	docker compose -f docker-compose.dev.yml logs -f
