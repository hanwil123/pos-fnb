package models

import "github.com/google/uuid"

type Restaurant struct {
	BaseModel
	Name    string `gorm:"not null" json:"name"`
	Address string `json:"address"`
	LogoURL string `json:"logo_url"`
}

// StaffRole merepresentasikan role staff: admin, cashier, kitchen.
type StaffRole string

const (
	RoleAdmin   StaffRole = "admin"
	RoleCashier StaffRole = "cashier"
	RoleKitchen StaffRole = "kitchen"
)

type StaffUser struct {
	BaseModel
	RestaurantID uuid.UUID `gorm:"type:uuid;not null;index" json:"restaurant_id"`
	Name         string    `gorm:"not null" json:"name"`
	Email        string    `gorm:"unique;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	Role         StaffRole `gorm:"type:varchar(20);not null" json:"role"`
}
