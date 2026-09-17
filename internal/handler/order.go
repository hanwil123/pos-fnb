package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/yourorg/pos-fnb-backend/internal/config"
	"github.com/yourorg/pos-fnb-backend/internal/dto"
	"github.com/yourorg/pos-fnb-backend/internal/services"
	"github.com/yourorg/pos-fnb-backend/internal/utils"
)

type OrderHandler struct {
	service services.OrderService
	cfg     *config.Config
}

func NewOrderHandler(service services.OrderService, cfg *config.Config) *OrderHandler {
	return &OrderHandler{
		service: service,
		cfg:     cfg,
	}
}

// CreateOrder — customer: submit order dari cart.
// POST /api/v1/orders
//
// PENTING: harga selalu diambil ulang dari database (menu_items), TIDAK
// dipercaya dari client, untuk mencegah manipulasi harga oleh customer.
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateOrderRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.Error(w, http.StatusBadRequest, "body request tidak valid")
		return
	}

	result, err := h.service.CreateOrder(&req, h.cfg.QRHMACSecret)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(w, http.StatusCreated, result)
}

// GetOrder — ambil detail order + item-nya.
// GET /api/v1/orders/{id}
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "id")

	result, err := h.service.GetOrder(orderID)
	if err != nil {
		utils.Error(w, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(w, http.StatusOK, result)
}

// ScanCashierQR — staff kasir: verifikasi QR yang ditunjukkan customer,
// lalu tampilkan detail order untuk dikonfirmasi kasir.
// POST /api/v1/staff/orders/scan-qr  body: { "token": "..." }
func (h *OrderHandler) ScanCashierQR(w http.ResponseWriter, r *http.Request) {
	var req dto.ScanCashierQRRequest
	if err := utils.DecodeJSON(r, &req); err != nil || req.Token == "" {
		utils.Error(w, http.StatusBadRequest, "token wajib diisi")
		return
	}

	result, err := h.service.ScanCashierQR(req.Token, h.cfg.QRHMACSecret)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	utils.Success(w, http.StatusOK, result)
}

// ConfirmCashPayment — staff kasir: konfirmasi pembayaran cash setelah scan QR.
// POST /api/v1/staff/payments/confirm  body: { "order_id": "..." }
func (h *OrderHandler) ConfirmCashPayment(w http.ResponseWriter, r *http.Request) {
	var req dto.ConfirmCashPaymentRequest
	if err := utils.DecodeJSON(r, &req); err != nil || req.OrderID == "" {
		utils.Error(w, http.StatusBadRequest, "order_id wajib diisi")
		return
	}

	if err := h.service.ConfirmCashPayment(req.OrderID); err != nil {
		utils.Error(w, http.StatusInternalServerError, "gagal konfirmasi pembayaran: "+err.Error())
		return
	}

	utils.Success(w, http.StatusOK, map[string]string{"message": "pembayaran cash berhasil dikonfirmasi"})
}
