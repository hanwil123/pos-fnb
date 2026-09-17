package models

import "github.com/google/uuid"

type TableStatus string

const (
	TableAvailable TableStatus = "available"
	TableOccupied  TableStatus = "occupied"
)

// Table merepresentasikan meja fisik di restoran. QRToken bersifat statis
// (dicetak sekali di kertas/stiker) dan di-sign dengan HMAC agar tidak bisa dipalsukan.
type Table struct {
	BaseModel
	RestaurantID uuid.UUID   `gorm:"type:uuid;not null;index" json:"restaurant_id"`
	TableNumber  string      `gorm:"not null" json:"table_number"`
	QRToken      string      `gorm:"unique;not null;index" json:"qr_token"`
	Status       TableStatus `gorm:"type:varchar(20);default:'available'" json:"status"`
}

type SessionStatus string

const (
	SessionActive  SessionStatus = "active"
	SessionClosed  SessionStatus = "closed"
	SessionExpired SessionStatus = "expired"
)

// TableSession dibuat setiap kali QR meja discan, merepresentasikan satu
// "kunjungan" customer di meja tersebut sampai session ditutup/expired.
type TableSession struct {
	BaseModel
	TableID      uuid.UUID     `gorm:"type:uuid;not null;index" json:"table_id"`
	SessionToken string        `gorm:"unique;not null;index" json:"session_token"`
	Status       SessionStatus `gorm:"type:varchar(20);default:'active'" json:"status"`
	ExpiresAt    *string       `json:"expires_at,omitempty"`
}
