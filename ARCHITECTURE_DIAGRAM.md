# Clean Architecture Diagram

## Layer Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         HTTP REQUEST                             │
│                         (Customer / Staff)                       │
└─────────────────────────────────┬───────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────┐
│                      HANDLER LAYER (Thin)                        │
│  ┌───────────────┐  ┌───────────────┐  ┌───────────────┐       │
│  │  MenuHandler  │  │  TableHandler │  │  OrderHandler │       │
│  │               │  │               │  │               │       │
│  │ - GetMenu()   │  │ - Resolve()   │  │ - Create()    │       │
│  │ - CreateItem()│  │ - List()      │  │ - ScanQR()    │       │
│  └───────┬───────┘  └───────┬───────┘  └───────┬───────┘       │
│          │                  │                  │                 │
│          │  Inject Service  │                  │                 │
└──────────┼──────────────────┼──────────────────┼─────────────────┘
           │                  │                  │
           ▼                  ▼                  ▼
┌─────────────────────────────────────────────────────────────────┐
│                     SERVICE LAYER (Business Logic)               │
│  ┌───────────────┐  ┌───────────────┐  ┌───────────────┐       │
│  │  MenuService  │  │  TableService │  │  OrderService │       │
│  │               │  │               │  │               │       │
│  │ - Validasi    │  │ - QR Token    │  │ - Price Check │       │
│  │ - Transform   │  │ - Session Mgmt│  │ - Calculate   │       │
│  │ - Group Data  │  │ - Transform   │  │ - QR Generate │       │
│  └───────┬───────┘  └───────┬───────┘  └───────┬───────┘       │
│          │                  │                  │                 │
│          │  Inject Repo     │                  │                 │
└──────────┼──────────────────┼──────────────────┼─────────────────┘
           │                  │                  │
           ▼                  ▼                  ▼
┌─────────────────────────────────────────────────────────────────┐
│                   REPOSITORY LAYER (Data Access)                 │
│  ┌───────────────┐  ┌───────────────┐  ┌───────────────┐       │
│  │  MenuRepo     │  │  TableRepo    │  │  OrderRepo    │       │
│  │               │  │               │  │               │       │
│  │ - GetItems()  │  │ - FindByQR()  │  │ - Create()    │       │
│  │ - Create()    │  │ - Create()    │  │ - FindByID()  │       │
│  │ - Update()    │  │ - Session()   │  │ - Update()    │       │
│  └───────┬───────┘  └───────┬───────┘  └───────┬───────┘       │
│          │                  │                  │                 │
│          │  Use GORM        │                  │                 │
└──────────┼──────────────────┼──────────────────┼─────────────────┘
           │                  │                  │
           └──────────────────┴──────────────────┘
                              │
                              ▼
           ┌──────────────────────────────────┐
           │       PostgreSQL Database         │
           │  ┌────────┐  ┌────────┐          │
           │  │ Tables │  │ Orders │          │
           │  └────────┘  └────────┘          │
           │  ┌────────┐  ┌────────┐          │
           │  │  Menu  │  │Sessions│          │
           │  └────────┘  └────────┘          │
           └──────────────────────────────────┘
```

---

## Data Flow Example: Create Order

```
┌──────────┐
│ Customer │
└────┬─────┘
     │ POST /api/v1/orders
     │ { items: [...], payment_method: "cashier" }
     ▼
┌─────────────────────────────────────────────────────┐
│ OrderHandler.CreateOrder()                          │
│ - Parse JSON request ke dto.CreateOrderRequest      │
│ - Delegate ke service                               │
└────────────────────┬────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────┐
│ OrderService.CreateOrder()                          │
│ 1. Validasi payment_method, items tidak kosong      │
│ 2. Loop items:                                      │
│    - Parse UUID menu_item_id                        │
│    - Call menuRepo.FindItemByID() → get real price  │
│    - Calculate subtotal                             │
│ 3. Calculate total_amount                           │
│ 4. orderRepo.Create(order)                          │
│ 5. orderRepo.CreateItems(orderItems)                │
│ 6. If payment_method == "cashier":                  │
│    - tableRepo.FindByID() → get table_number        │
│    - utils.GenerateSignedToken() → QR payload       │
│    - orderRepo.Update(order with QR)                │
│ 7. Transform models → dto.OrderResponse             │
└────────────────────┬────────────────────────────────┘
                     │
                     ▼
         ┌───────────────────────┐
         │  Multiple Repositories │
         │  - MenuRepo            │
         │  - OrderRepo           │
         │  - TableRepo           │
         └──────────┬─────────────┘
                    │
                    ▼
         ┌──────────────────┐
         │   GORM Queries   │
         │   to PostgreSQL  │
         └──────────────────┘
```

---

## Dependency Injection Flow

```
main.go
  │
  ├─ config.Load()
  ├─ database.Connect()
  │
  └─ router.Setup(db, cfg)
       │
       ├─ repository.NewMenuRepository(db)        ──┐
       ├─ repository.NewTableRepository(db)       ──┼─ Inject ke Service
       ├─ repository.NewOrderRepository(db)       ──┘
       │
       ├─ services.NewMenuService(menuRepo)       ──┐
       ├─ services.NewTableService(tableRepo)     ──┼─ Inject ke Handler
       ├─ services.NewOrderService(orderRepo, ...) ─┘
       │
       ├─ handler.NewMenuHandler(menuService)     ──┐
       ├─ handler.NewTableHandler(tableService)   ──┼─ Register Routes
       └─ handler.NewOrderHandler(orderService)   ──┘
```

---

## Interface Benefits

### Testability Example

```go
// Production: inject real repository
menuService := services.NewMenuService(repository.NewMenuRepository(db))

// Unit Test: inject mock repository
mockRepo := &MockMenuRepository{
    GetCategoriesFunc: func(id uuid.UUID) ([]models.MenuCategory, error) {
        return []models.MenuCategory{{Name: "Appetizers"}}, nil
    },
}
menuService := services.NewMenuService(mockRepo)

// Test service logic tanpa database
result, _ := menuService.GetMenu("restaurant-id")
assert.Equal(t, "Appetizers", result[0].Name)
```

---

## Folder Structure Tree

```
pos-fnb-backend/
├── cmd/
│   └── api/
│       └── main.go                 # Entry point
├── internal/
│   ├── dto/                        # 🆕 Request/Response structures
│   │   ├── menu.go
│   │   ├── table.go
│   │   └── order.go
│   │
│   ├── repository/                 # 🆕 Data access layer
│   │   ├── menuRepo.go
│   │   ├── tableRepo.go
│   │   └── orderRepo.go
│   │
│   ├── services/                   # 🆕 Business logic layer
│   │   ├── menuServices.go
│   │   ├── tableService.go
│   │   └── orderService.go
│   │
│   ├── handler/                    # ♻️ Refactored (thin layer)
│   │   ├── menu.go
│   │   ├── table.go
│   │   ├── order.go
│   │   └── health.go
│   │
│   ├── router/                     # ♻️ Updated (DI wiring)
│   │   └── router.go
│   │
│   ├── models/                     # Database entities
│   ├── middleware/                 # HTTP middleware
│   ├── config/                     # Env config
│   ├── database/                   # DB connection
│   └── utils/                      # Helpers
│
├── ARCHITECTURE.md                 # 🆕 Architecture docs
├── REFACTORING_SUMMARY.md          # 🆕 Change summary
├── ARCHITECTURE_DIAGRAM.md         # 🆕 This file
├── go.mod
└── go.sum
```

**Legend:**
- 🆕 = New files
- ♻️ = Refactored files
