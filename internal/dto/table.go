package dto

// ── Requests ──────────────────────────────────────────────────────────────────

type CreateTableRequest struct {
	RestaurantID string `json:"restaurant_id"`
	TableNumber  string `json:"table_number"`
}

// ── Responses ─────────────────────────────────────────────────────────────────

type TableResponse struct {
	ID           string `json:"id"`
	RestaurantID string `json:"restaurant_id"`
	TableNumber  string `json:"table_number"`
	QRToken      string `json:"qr_token"`
	Status       string `json:"status"`
	QRCodeImage  string `json:"qr_code_image,omitempty"` // Base64 encoded QR code image
}

type TableSessionResponse struct {
	Table        TableResponse `json:"table"`
	SessionToken string        `json:"session_token"`
	SessionID    string        `json:"session_id"`
}
