package dto

// ── Requests ──────────────────────────────────────────────────────────────────

type CreateOrderItemRequest struct {
	MenuItemID string `json:"menu_item_id"`
	Quantity   int    `json:"quantity"`
	Notes      string `json:"notes"`
}

type CreateOrderRequest struct {
	SessionID     string                   `json:"session_id"`
	TableID       string                   `json:"table_id"`
	CustomerName  string                   `json:"customer_name"`
	PaymentMethod string                   `json:"payment_method"` // "online" atau "cashier"
	Items         []CreateOrderItemRequest `json:"items"`
}

type ScanCashierQRRequest struct {
	Token string `json:"token"`
}

type ConfirmCashPaymentRequest struct {
	OrderID string `json:"order_id"`
}

// ── Responses ─────────────────────────────────────────────────────────────────

type OrderItemResponse struct {
	ID         string  `json:"id"`
	OrderID    string  `json:"order_id"`
	MenuItemID string  `json:"menu_item_id"`
	Quantity   int     `json:"quantity"`
	Notes      string  `json:"notes"`
	Subtotal   float64 `json:"subtotal"`
}

type OrderResponse struct {
	ID             string              `json:"id"`
	SessionID      string              `json:"session_id"`
	TableID        string              `json:"table_id"`
	CustomerName   string              `json:"customer_name"`
	Status         string              `json:"status"`
	PaymentMethod  string              `json:"payment_method"`
	PaymentStatus  string              `json:"payment_status"`
	TotalAmount    float64             `json:"total_amount"`
	OrderQRPayload string              `json:"order_qr_payload,omitempty"`
	Items          []OrderItemResponse `json:"items,omitempty"`
}
