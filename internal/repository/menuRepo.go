package repository

import (
	"github.com/google/uuid"
	"github.com/yourorg/pos-fnb-backend/internal/models"
	"gorm.io/gorm"
)

type MenuRepository interface {
	GetCategoriesByRestaurantID(restaurantID uuid.UUID) ([]models.MenuCategory, error)
	GetAvailableItemsByRestaurantID(restaurantID uuid.UUID) ([]models.MenuItem, error)
	GetSignatureItems(restaurantID uuid.UUID) ([]models.MenuItem, error)
	FindItemByID(id uuid.UUID) (*models.MenuItem, error)
	CreateItem(item *models.MenuItem) error
	UpdateItem(item *models.MenuItem) error
	DeleteItem(id uuid.UUID) error
	CreateCategory(category *models.MenuCategory) error
	GetCategory(id uuid.UUID) (*models.MenuCategory, error)
	// GetAllCategories() ([]models.MenuCategory, error)
	// UpdateCategory(category *models.MenuCategory) error
	// DeleteCategory(id uuid.UUID) error
}

type menuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) MenuRepository {
	return &menuRepository{db: db}
}

func (r *menuRepository) GetCategoriesByRestaurantID(restaurantID uuid.UUID) ([]models.MenuCategory, error) {
	var categories []models.MenuCategory
	if err := r.db.Where("restaurant_id = ?", restaurantID).
		Order("sort_order ASC").
		Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *menuRepository) GetAvailableItemsByRestaurantID(restaurantID uuid.UUID) ([]models.MenuItem, error) {
	var items []models.MenuItem
	if err := r.db.Where("restaurant_id = ? AND is_available = ?", restaurantID, true).
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *menuRepository) GetSignatureItems(restaurantID uuid.UUID) ([]models.MenuItem, error) {
	var items []models.MenuItem
	if err := r.db.Where("restaurant_id = ? AND is_signature = ? AND is_available = ?",
		restaurantID, true, true).
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *menuRepository) FindItemByID(id uuid.UUID) (*models.MenuItem, error) {
	var item models.MenuItem
	if err := r.db.Where("id = ? AND is_available = ?", id, true).
		First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *menuRepository) CreateItem(item *models.MenuItem) error {
	return r.db.Create(item).Error
}

func (r *menuRepository) UpdateItem(item *models.MenuItem) error {
	return r.db.Save(item).Error
}

func (r *menuRepository) DeleteItem(id uuid.UUID) error {
	return r.db.Delete(&models.MenuItem{}, "id = ?", id).Error
}

func (r *menuRepository) CreateCategory(category *models.MenuCategory) error {
	return r.db.Create(category).Error
}

func (r *menuRepository) GetCategory(id uuid.UUID) (*models.MenuCategory, error) {
	var category models.MenuCategory
	if err := r.db.Where("id = ?", id).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// func (r *menuRepository) GetAllCategories() ([]models.MenuCategory, error) {
// 	var categories []models.MenuCategory
// 	if err := r.db.Find(&categories).Error; err != nil {
// 		return nil, err
// 	}
// 	return categories, nil
// }
