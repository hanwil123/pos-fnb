package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseModel dipakai sebagai embed di semua model agar konsisten:
// UUID sebagai primary key (bukan auto-increment int) supaya aman
// dipakai di URL publik (QR code, dsb) tanpa membocorkan jumlah data.
type BaseModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate hook GORM: generate UUID otomatis sebelum insert jika belum diisi.
func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}
