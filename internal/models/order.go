package models

import "github.com/google/uuid"

type OrderStatus string

const (
	OrderPending   OrderStatus = "pending"
	OrderPaid      OrderStatus = "paid"
	OrderPreparing OrderStatus = "preparing"
	OrderServed    OrderStatus = "served"
	OrderCompleted OrderStatus = "completed"
	OrderCancelled OrderStatus = "cancelled"
)

type PaymentMethod string

const (
	PaymentOnline  PaymentMethod = "online"
	PaymentCashier PaymentMethod = "cashier"
)

type PaymentStatus string

const (
	PaymentUnpaid PaymentStatus = "unpaid"
	PaymentPaid   PaymentStatus = "paid"
	PaymentFailed PaymentStatus = "failed"
)

// Order merepresentasikan satu pesanan yang dibuat dari satu table session.
// OrderQRPayload diisi (signed JSON/JWT) ketika payment_method = cashier,
// untuk discan kasir saat konfirmasi pembayaran cash.
type Order struct {
	BaseModel
	SessionID       uuid.UUID     `gorm:"type:uuid;not null;index" json:"session_id"`
	TableID         uuid.UUID     `gorm:"type:uuid;not null;index" json:"table_id"`
	CustomerName    string        `json:"customer_name"`
	Status          OrderStatus   `gorm:"type:varchar(20);default:'pending'" json:"status"`
	PaymentMethod   PaymentMethod `gorm:"type:varchar(20)" json:"payment_method"`
	PaymentStatus   PaymentStatus `gorm:"type:varchar(20);default:'unpaid'" json:"payment_status"`
	TotalAmount     float64       `gorm:"type:numeric(12,2);not null" json:"total_amount"`
	OrderQRPayload  string        `json:"order_qr_payload,omitempty"`
	Items           []OrderItem   `gorm:"foreignKey:OrderID" json:"items,omitempty"`
}

type OrderItem struct {
	BaseModel
	OrderID    uuid.UUID `gorm:"type:uuid;not null;index" json:"order_id"`
	MenuItemID uuid.UUID `gorm:"type:uuid;not null;index" json:"menu_item_id"`
	Quantity   int       `gorm:"not null" json:"quantity"`
	Notes      string    `json:"notes"`
	Subtotal   float64   `gorm:"type:numeric(12,2);not null" json:"subtotal"`
}
