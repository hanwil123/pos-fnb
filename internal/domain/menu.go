package domain

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type MenuCategory struct {
	BaseModel
	RestaurantID uuid.UUID
	Name         string
	SortOrder    int
}

// MenuItem merepresentasikan satu item menu.
// Kolom Embedding sengaja belum diaktifkan di scaffold ini (butuh ekstensi
// pgvector + tipe kolom khusus) — ditambahkan nanti di Phase 3 (fitur AI).
type MenuItem struct {
	BaseModel
	RestaurantID uuid.UUID
	CategoryID   uuid.UUID
	Name         string
	Description  string
	Price        float64
	ImageURL     string
	IsSignature  bool
	IsAvailable  bool
	Tags         pq.StringArray
}
