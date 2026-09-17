package repository

import (
	"github.com/yourorg/pos-fnb-backend/internal/models"
	"gorm.io/gorm"
)

type RestaurantRepository interface {
	Create(restaurant *models.Restaurant) (*models.Restaurant, error)
}

type restaurantRepository struct {
	db *gorm.DB
}

func NewRestaurantRepository(db *gorm.DB) RestaurantRepository {
	return &restaurantRepository{db: db}
}

func (r *restaurantRepository) Create(restaurant *models.Restaurant) (*models.Restaurant, error) {
	if err := r.db.Create(restaurant).Error; err != nil {
		return nil, err
	}
	return restaurant, nil
}
