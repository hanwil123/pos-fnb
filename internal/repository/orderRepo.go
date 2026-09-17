package repository

import (
	"github.com/google/uuid"
	"github.com/yourorg/pos-fnb-backend/internal/models"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *models.Order) error
	CreateItems(items []models.OrderItem) error
	FindByID(id uuid.UUID) (*models.Order, error)
	FindByIDWithItems(id uuid.UUID) (*models.Order, error)
	Update(order *models.Order) error
	UpdateStatus(orderID uuid.UUID, status models.OrderStatus) error
	UpdatePaymentStatus(orderID uuid.UUID, paymentStatus models.PaymentStatus) error
	
	// Payment methods
	CreatePayment(payment *models.Payment) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) CreateItems(items []models.OrderItem) error {
	return r.db.Create(&items).Error
}

func (r *orderRepository) FindByID(id uuid.UUID) (*models.Order, error) {
	var order models.Order
	if err := r.db.First(&order, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) FindByIDWithItems(id uuid.UUID) (*models.Order, error) {
	var order models.Order
	if err := r.db.Preload("Items").First(&order, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) Update(order *models.Order) error {
	return r.db.Save(order).Error
}

func (r *orderRepository) UpdateStatus(orderID uuid.UUID, status models.OrderStatus) error {
	return r.db.Model(&models.Order{}).Where("id = ?", orderID).Update("status", status).Error
}

func (r *orderRepository) UpdatePaymentStatus(orderID uuid.UUID, paymentStatus models.PaymentStatus) error {
	return r.db.Model(&models.Order{}).Where("id = ?", orderID).Update("payment_status", paymentStatus).Error
}

func (r *orderRepository) CreatePayment(payment *models.Payment) error {
	return r.db.Create(payment).Error
}
