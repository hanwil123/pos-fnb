package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/yourorg/pos-fnb-backend/internal/config"
	"github.com/yourorg/pos-fnb-backend/internal/handler"
	"github.com/yourorg/pos-fnb-backend/internal/middleware"
	"github.com/yourorg/pos-fnb-backend/internal/repository"
	"github.com/yourorg/pos-fnb-backend/internal/services"
	"gorm.io/gorm"
)

// Setup mendaftarkan semua route API, dikelompokkan berdasarkan
// siapa yang mengakses: public/customer vs staff (protected).
//
// Dependency injection flow: DB → Repository → Service → Handler
func Setup(db *gorm.DB, cfg *config.Config) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.CORS)

	r.Get("/health", handler.HealthCheck)

	// ── Inisialisasi Repositories ──────────────────────────────────────────────
	menuRepo := repository.NewMenuRepository(db)
	tableRepo := repository.NewTableRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	restaurantRepo := repository.NewRestaurantRepository(db)

	// ── Inisialisasi Services ──────────────────────────────────────────────────
	menuService := services.NewMenuService(menuRepo)
	tableService := services.NewTableService(tableRepo)
	orderService := services.NewOrderService(orderRepo, menuRepo, tableRepo)
	restaurantService := services.NewRestaurantService(restaurantRepo)

	// ── Inisialisasi Handlers ──────────────────────────────────────────────────
	menuHandler := handler.NewMenuHandler(menuService)
	tableHandler := handler.NewTableHandler(tableService)
	orderHandler := handler.NewOrderHandler(orderService, cfg)
	restaurantHandler := handler.NewRestaurantHandler(restaurantService)

	// ── Routes ─────────────────────────────────────────────────────────────────
	r.Route("/api/v1", func(v1 chi.Router) {
		// ── Public routes (diakses customer via scan QR) ───────────────────────
		v1.Post("/restaurants", restaurantHandler.CreateRestaurant)
		v1.Post("/tabless", tableHandler.CreateTable)
		v1.Get("/table/{qr_token}", tableHandler.ResolveTable)
		v1.Post("/menu-itemss", menuHandler.CreateMenuItem)
		v1.Post("/orders/scan-qrr", orderHandler.ScanCashierQR)
		v1.Post("/payments/confirmm", orderHandler.ConfirmCashPayment)
		v1.Get("/menu", menuHandler.GetMenu)
		v1.Get("/menu/recommendations", menuHandler.GetRecommendations)
		v1.Post("/orders", orderHandler.CreateOrder)
		v1.Get("/orders/{id}", orderHandler.GetOrder)

		// ── Staff routes (butuh JWT, role admin/cashier/kitchen) ───────────────
		v1.Route("/staff", func(staff chi.Router) {
			staff.Use(middleware.StaffAuth(cfg.JWTSecret))

			staff.Get("/tables", tableHandler.ListTables)
			staff.Post("/tables", tableHandler.CreateTable)
			staff.Post("/menu-items", menuHandler.CreateMenuItem)
			staff.Post("/orders/scan-qr", orderHandler.ScanCashierQR)
			staff.Post("/payments/confirm", orderHandler.ConfirmCashPayment)
		})
	})

	return r
}
