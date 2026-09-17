package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type MenuCategory struct {
	BaseModel
	RestaurantID uuid.UUID `gorm:"type:uuid;not null;index" json:"restaurant_id"`
	Name         string    `gorm:"not null" json:"name"`
	SortOrder    int       `gorm:"default:0" json:"sort_order"`
}

// MenuItem merepresentasikan satu item menu.
// Kolom Embedding sengaja belum diaktifkan di scaffold ini (butuh ekstensi
// pgvector + tipe kolom khusus) — ditambahkan nanti di Phase 3 (fitur AI).
type MenuItem struct {
	BaseModel
	RestaurantID uuid.UUID      `gorm:"type:uuid;not null;index" json:"restaurant_id"`
	CategoryID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"category_id"`
	Name         string         `gorm:"not null" json:"name"`
	Description  string         `json:"description"`
	Price        float64        `gorm:"type:numeric(12,2);not null" json:"price"`
	ImageURL     string         `json:"image_url"`
	IsSignature  bool           `gorm:"default:false" json:"is_signature"`
	IsAvailable  bool           `gorm:"default:true" json:"is_available"`
	Tags         pq.StringArray `gorm:"type:text[]" json:"tags"`
}

