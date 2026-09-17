# Docker Helper Script untuk POS FnB Backend
param(
    [Parameter(Position=0)]
    [string]$Command = "help"
)

function Show-Help {
    Write-Host "🐳 Docker Helper Script - POS FnB Backend" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Perintah yang tersedia:" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Production:" -ForegroundColor Green
    Write-Host "  .\docker.ps1 build     - Build semua Docker images"
    Write-Host "  .\docker.ps1 up        - Jalankan semua services"
    Write-Host "  .\docker.ps1 down      - Stop semua services"
    Write-Host "  .\docker.ps1 restart   - Restart semua services"
    Write-Host "  .\docker.ps1 logs      - Lihat logs semua services"
    Write-Host ""
    Write-Host "Development (hot reload):" -ForegroundColor Cyan
    Write-Host "  .\docker.ps1 dev-up    - Start development mode dengan hot reload"
    Write-Host "  .\docker.ps1 dev-down  - Stop development mode"
    Write-Host "  .\docker.ps1 dev-logs  - Lihat logs development mode"
    Write-Host ""
    Write-Host "Utilities:" -ForegroundColor Yellow
    Write-Host "  .\docker.ps1 clean     - Stop services dan hapus volumes"
    Write-Host "  .\docker.ps1 test      - Jalankan test di dalam container"
    Write-Host "  .\docker.ps1 shell     - Masuk ke shell API container"
    Write-Host "  .\docker.ps1 db-shell  - Masuk ke PostgreSQL shell"
}

function Build-Images {
    Write-Host "🔨 Building Docker images..." -ForegroundColor Yellow
    docker compose build
}

function Start-Services {
    Write-Host "🚀 Starting services..." -ForegroundColor Green
    docker compose up -d
    Write-Host ""
    Write-Host "✅ Services started successfully!" -ForegroundColor Green
    Write-Host "   API: http://localhost:8080" -ForegroundColor Cyan
    Write-Host "   PostgreSQL: localhost:5432" -ForegroundColor Cyan
    Write-Host "   Redis: localhost:6379" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Gunakan '.\docker.ps1 logs' untuk melihat logs" -ForegroundColor Yellow
}

function Stop-Services {
    Write-Host "🛑 Stopping services..." -ForegroundColor Yellow
    docker compose down
}

function Restart-Services {
    Write-Host "🔄 Restarting services..." -ForegroundColor Yellow
    docker compose restart
}

function Show-Logs {
    Write-Host "📋 Showing logs (Ctrl+C untuk keluar)..." -ForegroundColor Yellow
    docker compose logs -f
}

function Clean-All {
    Write-Host "🧹 Cleaning up..." -ForegroundColor Red
    $confirm = Read-Host "Ini akan menghapus semua data di database dan redis. Lanjutkan? (y/N)"
    if ($confirm -eq 'y' -or $confirm -eq 'Y') {
        docker compose down -v
        Write-Host "✅ Cleanup completed!" -ForegroundColor Green
    } else {
        Write-Host "❌ Cleanup cancelled" -ForegroundColor Yellow
    }
}

function Run-Tests {
    Write-Host "🧪 Running tests..." -ForegroundColor Yellow
    docker compose exec api go test -v ./...
}

function Enter-Shell {
    Write-Host "🐚 Entering API container shell..." -ForegroundColor Yellow
    docker compose exec api sh
}

function Enter-DBShell {
    Write-Host "🗄️  Entering PostgreSQL shell..." -ForegroundColor Yellow
    docker compose exec postgres psql -U postgres -d pos_fnb
}

function Dev-Mode {
    Write-Host "👨‍💻 Starting development mode..." -ForegroundColor Cyan
    docker compose down
    docker compose build
    docker compose up
}

function Dev-Up {
    Write-Host "🚀 Starting development mode dengan hot reload..." -ForegroundColor Cyan
    docker compose -f docker-compose.dev.yml up --build
}

function Dev-Down {
    Write-Host "🛑 Stopping development mode..." -ForegroundColor Yellow
    docker compose -f docker-compose.dev.yml down
}

function Dev-Logs {
    Write-Host "📋 Development logs (Ctrl+C untuk keluar)..." -ForegroundColor Yellow
    docker compose -f docker-compose.dev.yml logs -f
}

# Main switch
switch ($Command.ToLower()) {
    "help"      { Show-Help }
    "build"     { Build-Images }
    "up"        { Start-Services }
    "down"      { Stop-Services }
    "restart"   { Restart-Services }
    "logs"      { Show-Logs }
    "clean"     { Clean-All }
    "test"      { Run-Tests }
    "shell"     { Enter-Shell }
    "db-shell"  { Enter-DBShell }
    "dev"       { Dev-Mode }
    "dev-up"    { Dev-Up }
    "dev-down"  { Dev-Down }
    "dev-logs"  { Dev-Logs }
    default     { 
        Write-Host "❌ Perintah tidak dikenali: $Command" -ForegroundColor Red
        Write-Host ""
        Show-Help 
    }
}
