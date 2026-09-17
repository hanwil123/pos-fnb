package repository

import (
	"github.com/google/uuid"
	"github.com/yourorg/pos-fnb-backend/internal/models"
	"gorm.io/gorm"
)

type TableRepository interface {
	FindByQRToken(qrToken string) (*models.Table, error)
	FindByID(id uuid.UUID) (*models.Table, error)
	FindByRestaurantID(restaurantID uuid.UUID) ([]models.Table, error)
	FindAll() ([]models.Table, error)
	Create(table *models.Table) error
	Update(table *models.Table) error
	
	// Table Session methods
	FindActiveSession(tableID uuid.UUID) (*models.TableSession, error)
	CreateSession(session *models.TableSession) error
	UpdateSession(session *models.TableSession) error
}

type tableRepository struct {
	db *gorm.DB
}

func NewTableRepository(db *gorm.DB) TableRepository {
	return &tableRepository{db: db}
}

func (r *tableRepository) FindByQRToken(qrToken string) (*models.Table, error) {
	var table models.Table
	if err := r.db.Where("qr_token = ?", qrToken).First(&table).Error; err != nil {
		return nil, err
	}
	return &table, nil
}

func (r *tableRepository) FindByID(id uuid.UUID) (*models.Table, error) {
	var table models.Table
	if err := r.db.First(&table, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &table, nil
}

func (r *tableRepository) FindByRestaurantID(restaurantID uuid.UUID) ([]models.Table, error) {
	var tables []models.Table
	if err := r.db.Where("restaurant_id = ?", restaurantID).Find(&tables).Error; err != nil {
		return nil, err
	}
	return tables, nil
}

func (r *tableRepository) FindAll() ([]models.Table, error) {
	var tables []models.Table
	if err := r.db.Find(&tables).Error; err != nil {
		return nil, err
	}
	return tables, nil
}

func (r *tableRepository) Create(table *models.Table) error {
	return r.db.Create(table).Error
}

func (r *tableRepository) Update(table *models.Table) error {
	return r.db.Save(table).Error
}

func (r *tableRepository) FindActiveSession(tableID uuid.UUID) (*models.TableSession, error) {
	var session models.TableSession
	err := r.db.Where("table_id = ? AND status = ?", tableID, models.SessionActive).
		Order("created_at DESC").
		First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *tableRepository) CreateSession(session *models.TableSession) error {
	return r.db.Create(session).Error
}

func (r *tableRepository) UpdateSession(session *models.TableSession) error {
	return r.db.Save(session).Error
}
