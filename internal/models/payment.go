package models

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	BaseModel
	OrderID          uuid.UUID  `gorm:"type:uuid;not null;index" json:"order_id"`
	Provider         string     `gorm:"type:varchar(30)" json:"provider"` // midtrans, xendit, cash
	ProviderRef      string     `json:"provider_ref"`
	Amount           float64    `gorm:"type:numeric(12,2)" json:"amount"`
	Status           string     `gorm:"type:varchar(20)" json:"status"` // pending, success, failed
	PaidAt           *time.Time `json:"paid_at,omitempty"`
	HandledByStaffID *uuid.UUID `gorm:"type:uuid" json:"handled_by_staff_id,omitempty"`
}

type Review struct {
	BaseModel
	OrderID        uuid.UUID `gorm:"type:uuid;not null;index" json:"order_id"`
	Rating         int       `gorm:"check:rating >= 1 AND rating <= 5" json:"rating"`
	Comment        string    `json:"comment"`
	SentimentScore float64   `json:"sentiment_score"`
}
