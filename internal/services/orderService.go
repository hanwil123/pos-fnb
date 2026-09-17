package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/yourorg/pos-fnb-backend/internal/dto"
	"github.com/yourorg/pos-fnb-backend/internal/models"
	"github.com/yourorg/pos-fnb-backend/internal/repository"
	"github.com/yourorg/pos-fnb-backend/internal/utils"
)

type OrderService interface {
	CreateOrder(req *dto.CreateOrderRequest, qrHMACSecret string) (*dto.OrderResponse, error)
	GetOrder(orderID string) (*dto.OrderResponse, error)
	ScanCashierQR(token string, qrHMACSecret string) (*dto.OrderResponse, error)
	ConfirmCashPayment(orderID string) error
}

type orderService struct {
	orderRepo repository.OrderRepository
	menuRepo  repository.MenuRepository
	tableRepo repository.TableRepository
}

func NewOrderService(
	orderRepo repository.OrderRepository,
	menuRepo repository.MenuRepository,
	tableRepo repository.TableRepository,
) OrderService {
	return &orderService{
		orderRepo: orderRepo,
		menuRepo:  menuRepo,
		tableRepo: tableRepo,
	}
}

func (s *orderService) CreateOrder(req *dto.CreateOrderRequest, qrHMACSecret string) (*dto.OrderResponse, error) {
	sessionID, err := uuid.Parse(req.SessionID)
	if err != nil {
		return nil, errors.New("session_id tidak valid")
	}

	tableID, err := uuid.Parse(req.TableID)
	if err != nil {
		return nil, errors.New("table_id tidak valid")
	}

	if len(req.Items) == 0 {
		return nil, errors.New("items tidak boleh kosong")
	}

	if req.PaymentMethod != string(models.PaymentOnline) && req.PaymentMethod != string(models.PaymentCashier) {
		return nil, errors.New("payment_method harus 'online' atau 'cashier'")
	}

	// Hitung total dan buat order items
	var totalAmount float64
	orderItems := make([]models.OrderItem, 0, len(req.Items))

	for _, reqItem := range req.Items {
		menuItemID, err := uuid.Parse(reqItem.MenuItemID)
		if err != nil {
			return nil, err
		}

		quantity := reqItem.Quantity
		if quantity < 1 {
			quantity = 1
		}

		// PENTING: harga selalu diambil dari database, TIDAK dipercaya dari client
		menuItem, err := s.menuRepo.FindItemByID(menuItemID)
		if err != nil {
			return nil, errors.New("menu item tidak ditemukan atau tidak tersedia")
		}

		subtotal := menuItem.Price * float64(quantity)
		totalAmount += subtotal

		orderItems = append(orderItems, models.OrderItem{
			MenuItemID: menuItemID,
			Quantity:   quantity,
			Notes:      reqItem.Notes,
			Subtotal:   subtotal,
		})
	}

	// Buat order
	order := &models.Order{
		SessionID:     sessionID,
		TableID:       tableID,
		CustomerName:  req.CustomerName,
		Status:        models.OrderPending,
		PaymentMethod: models.PaymentMethod(req.PaymentMethod),
		PaymentStatus: models.PaymentUnpaid,
		TotalAmount:   totalAmount,
	}

	if err := s.orderRepo.Create(order); err != nil {
		return nil, err
	}

	// Buat order items
	for i := range orderItems {
		orderItems[i].OrderID = order.ID
	}
	if err := s.orderRepo.CreateItems(orderItems); err != nil {
		return nil, err
	}

	// Jika bayar di kasir, generate signed QR payload
	if order.PaymentMethod == models.PaymentCashier {
		table, err := s.tableRepo.FindByID(tableID)
		if err != nil {
			return nil, err
		}

		payload := utils.NewCashierQRPayload(order.ID.String(), table.TableNumber, int64(totalAmount))
		token, err := utils.GenerateSignedToken(payload, qrHMACSecret)
		if err != nil {
			return nil, err
		}

		order.OrderQRPayload = token
		if err := s.orderRepo.Update(order); err != nil {
			return nil, err
		}
	}

	// Convert to response
	items := make([]dto.OrderItemResponse, 0, len(orderItems))
	for _, item := range orderItems {
		items = append(items, dto.OrderItemResponse{
			ID:         item.ID.String(),
			OrderID:    item.OrderID.String(),
			MenuItemID: item.MenuItemID.String(),
			Quantity:   item.Quantity,
			Notes:      item.Notes,
			Subtotal:   item.Subtotal,
		})
	}

	return &dto.OrderResponse{
		ID:             order.ID.String(),
		SessionID:      order.SessionID.String(),
		TableID:        order.TableID.String(),
		CustomerName:   order.CustomerName,
		Status:         string(order.Status),
		PaymentMethod:  string(order.PaymentMethod),
		PaymentStatus:  string(order.PaymentStatus),
		TotalAmount:    order.TotalAmount,
		OrderQRPayload: order.OrderQRPayload,
		Items:          items,
	}, nil
}

func (s *orderService) GetOrder(orderID string) (*dto.OrderResponse, error) {
	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		return nil, errors.New("order_id tidak valid")
	}

	order, err := s.orderRepo.FindByIDWithItems(orderUUID)
	if err != nil {
		return nil, errors.New("order tidak ditemukan")
	}

	items := make([]dto.OrderItemResponse, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, dto.OrderItemResponse{
			ID:         item.ID.String(),
			OrderID:    item.OrderID.String(),
			MenuItemID: item.MenuItemID.String(),
			Quantity:   item.Quantity,
			Notes:      item.Notes,
			Subtotal:   item.Subtotal,
		})
	}

	return &dto.OrderResponse{
		ID:             order.ID.String(),
		SessionID:      order.SessionID.String(),
		TableID:        order.TableID.String(),
		CustomerName:   order.CustomerName,
		Status:         string(order.Status),
		PaymentMethod:  string(order.PaymentMethod),
		PaymentStatus:  string(order.PaymentStatus),
		TotalAmount:    order.TotalAmount,
		OrderQRPayload: order.OrderQRPayload,
		Items:          items,
	}, nil
}

func (s *orderService) ScanCashierQR(token string, qrHMACSecret string) (*dto.OrderResponse, error) {
	var payload utils.QRPayload
	if err := utils.VerifySignedToken(token, qrHMACSecret, &payload); err != nil {
		return nil, errors.New("QR tidak valid")
	}

	orderUUID, err := uuid.Parse(payload.OrderID)
	if err != nil {
		return nil, errors.New("order_id dalam QR tidak valid")
	}

	order, err := s.orderRepo.FindByIDWithItems(orderUUID)
	if err != nil {
		return nil, errors.New("order tidak ditemukan")
	}

	if order.PaymentStatus == models.PaymentPaid {
		return nil, errors.New("order ini sudah dibayar sebelumnya")
	}

	items := make([]dto.OrderItemResponse, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, dto.OrderItemResponse{
			ID:         item.ID.String(),
			OrderID:    item.OrderID.String(),
			MenuItemID: item.MenuItemID.String(),
			Quantity:   item.Quantity,
			Notes:      item.Notes,
			Subtotal:   item.Subtotal,
		})
	}

	return &dto.OrderResponse{
		ID:             order.ID.String(),
		SessionID:      order.SessionID.String(),
		TableID:        order.TableID.String(),
		CustomerName:   order.CustomerName,
		Status:         string(order.Status),
		PaymentMethod:  string(order.PaymentMethod),
		PaymentStatus:  string(order.PaymentStatus),
		TotalAmount:    order.TotalAmount,
		OrderQRPayload: order.OrderQRPayload,
		Items:          items,
	}, nil
}

func (s *orderService) ConfirmCashPayment(orderID string) error {
	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		return errors.New("order_id tidak valid")
	}

	order, err := s.orderRepo.FindByID(orderUUID)
	if err != nil {
		return errors.New("order tidak ditemukan")
	}

	// Update payment status dan order status
	if err := s.orderRepo.UpdatePaymentStatus(order.ID, models.PaymentPaid); err != nil {
		return err
	}
	if err := s.orderRepo.UpdateStatus(order.ID, models.OrderPaid); err != nil {
		return err
	}

	// Catat payment record
	payment := &models.Payment{
		OrderID:  order.ID,
		Provider: "cash",
		Amount:   order.TotalAmount,
		Status:   "success",
	}
	if err := s.orderRepo.CreatePayment(payment); err != nil {
		return err
	}

	return nil
}
