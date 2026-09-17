# Clean Architecture — Dokumentasi

## Struktur Folder

Setelah refactoring, struktur project mengikuti **Clean Architecture** dengan pemisahan layer yang jelas:

```
internal/
├── dto/                   # Data Transfer Objects (Request & Response)
│   ├── menu.go           # Menu DTO (request & response)
│   ├── table.go          # Table DTO (request & response)
│   └── order.go          # Order DTO (request & response)
│
├── repository/            # Data Access Layer (komunikasi dengan database)
│   ├── menuRepo.go       # Menu repository interface & implementasi
│   ├── tableRepo.go      # Table repository interface & implementasi
│   └── orderRepo.go      # Order repository interface & implementasi
│
├── services/              # Business Logic Layer
│   ├── menuServices.go   # Menu service interface & implementasi
│   ├── tableService.go   # Table service interface & implementasi
│   └── orderService.go   # Order service interface & implementasi
│
├── handler/               # HTTP Handler Layer (thin layer, no business logic)
│   ├── menu.go           # Menu HTTP handlers
│   ├── table.go          # Table HTTP handlers
│   ├── order.go          # Order HTTP handlers
│   └── health.go         # Health check handler
│
├── models/                # Database Models (GORM entities)
├── middleware/            # HTTP Middleware (CORS, Auth, etc)
├── router/                # Route definitions & dependency injection
├── config/                # Environment config loader
├── database/              # Database connection & migration
└── utils/                 # Helper functions (response, QR token, etc)
```

---

## Flow Dependency Injection

**Router.go** adalah tempat wiring semua dependency:

```
DB (gorm.DB)
  ↓
Repository (interface + implementasi)
  ↓
Service (interface + implementasi)
  ↓
Handler (HTTP layer)
```

### Contoh wiring di `router/router.go`:

```go
// 1. Inisialisasi Repositories
menuRepo := repository.NewMenuRepository(db)
tableRepo := repository.NewTableRepository(db)
orderRepo := repository.NewOrderRepository(db)

// 2. Inisialisasi Services (inject repositories)
menuService := services.NewMenuService(menuRepo)
tableService := services.NewTableService(tableRepo)
orderService := services.NewOrderService(orderRepo, menuRepo, tableRepo)

// 3. Inisialisasi Handlers (inject services)
menuHandler := handler.NewMenuHandler(menuService)
tableHandler := handler.NewTableHandler(tableService)
orderHandler := handler.NewOrderHandler(orderService, cfg)
```

---

## Tanggung Jawab Setiap Layer

### 1. **DTO (Data Transfer Object)**
- Mendefinisikan struktur request dan response API
- Memisahkan struktur HTTP dari database models
- Validasi input di layer ini (bisa diperluas dengan validator library)

**Contoh:**
```go
type CreateMenuItemRequest struct {
    RestaurantID string   `json:"restaurant_id"`
    CategoryID   string   `json:"category_id"`
    Name         string   `json:"name"`
    Price        float64  `json:"price"`
    ...
}
```

### 2. **Repository**
- Interface untuk data access (agar mudah di-mock untuk testing)
- Semua query database ada di sini
- Tidak boleh ada business logic
- Return models dari database

**Contoh:**
```go
type MenuRepository interface {
    GetCategoriesByRestaurantID(restaurantID uuid.UUID) ([]models.MenuCategory, error)
    GetAvailableItemsByRestaurantID(restaurantID uuid.UUID) ([]models.MenuItem, error)
    CreateItem(item *models.MenuItem) error
    ...
}
```

### 3. **Service**
- Interface untuk business logic
- Koordinasi antar repository
- Validasi business rules
- Transformasi dari models ke DTO
- Transaction management (jika kompleks)

**Contoh:**
```go
type MenuService interface {
    GetMenu(restaurantID string) ([]dto.MenuCategoryResponse, error)
    CreateMenuItem(req *dto.CreateMenuItemRequest) (*dto.MenuItemResponse, error)
}
```

### 4. **Handler**
- Thin layer, hanya HTTP concerns (parsing request, response)
- Tidak ada business logic atau database query
- Delegate semua logic ke Service
- Return JSON response dengan utils helper

**Contoh:**
```go
func (h *MenuHandler) GetMenu(w http.ResponseWriter, r *http.Request) {
    restaurantID := r.URL.Query().Get("restaurant_id")
    
    result, err := h.service.GetMenu(restaurantID)
    if err != nil {
        utils.Error(w, http.StatusInternalServerError, err.Error())
        return
    }
    
    utils.Success(w, http.StatusOK, result)
}
```

---

## Keuntungan Arsitektur Ini

✅ **Separation of Concerns**: Setiap layer punya tanggung jawab yang jelas  
✅ **Testability**: Mudah membuat unit test dengan mock interface  
✅ **Maintainability**: Perubahan di satu layer tidak affect layer lain  
✅ **Scalability**: Mudah menambah fitur baru tanpa ubah existing code  
✅ **Reusability**: Service bisa dipanggil dari berbagai handler atau background job  
✅ **Clean Dependencies**: Dependency mengalir satu arah (Handler → Service → Repository → DB)

---

## Contoh Use Case: Create Order

**Request Flow:**
```
1. Handler.CreateOrder()
   ↓ Parse request body ke DTO
   ↓ Validasi basic (empty fields, dll)
   ↓
2. Service.CreateOrder(dto)
   ↓ Validasi business logic (payment method, items tidak kosong)
   ↓ Loop items → ambil harga dari MenuRepository (security: prevent price manipulation)
   ↓ Hitung total amount
   ↓ Buat order via OrderRepository
   ↓ Buat order items via OrderRepository
   ↓ Jika cashier payment → generate QR token via Utils
   ↓ Convert models ke DTO response
   ↓
3. Handler.CreateOrder()
   ↓ Return JSON response
```

**Code:**
```go
// Handler (thin)
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
    var req dto.CreateOrderRequest
    utils.DecodeJSON(r, &req)
    
    result, err := h.service.CreateOrder(&req, h.cfg.QRHMACSecret)
    if err != nil {
        utils.Error(w, http.StatusBadRequest, err.Error())
        return
    }
    
    utils.Success(w, http.StatusCreated, result)
}

// Service (business logic)
func (s *orderService) CreateOrder(req *dto.CreateOrderRequest, qrSecret string) (*dto.OrderResponse, error) {
    // Validasi
    // Loop items → get real price from menuRepo
    // Calculate total
    // Create order via orderRepo
    // Generate QR if needed
    // Transform to DTO
    return &dto.OrderResponse{...}, nil
}

// Repository (data access)
func (r *orderRepository) Create(order *models.Order) error {
    return r.db.Create(order).Error
}
```

---

## Next Steps (Rekomendasi)

1. **Unit Testing**: Buat unit test untuk service layer dengan mock repository
2. **Validation**: Tambahkan library validator (e.g., `go-playground/validator`) di DTO
3. **Error Handling**: Buat custom error types untuk handling error lebih spesifik
4. **Logging**: Tambahkan structured logging (e.g., `zap`, `logrus`) di service layer
5. **Transaction**: Wrap complex operations di service dengan database transaction
6. **Pagination**: Tambahkan pagination di repository untuk list endpoints

---

## Changelog

**[2026-09-16] Refactoring ke Clean Architecture**
- ✅ Pisahkan DTO dari models
- ✅ Buat repository layer dengan interface untuk menu, table, order
- ✅ Buat service layer dengan business logic
- ✅ Refactor handler menjadi thin layer (no business logic)
- ✅ Update router dengan proper dependency injection
- ✅ Build success tanpa error
