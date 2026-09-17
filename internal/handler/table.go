package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/yourorg/pos-fnb-backend/internal/dto"
	"github.com/yourorg/pos-fnb-backend/internal/services"
	"github.com/yourorg/pos-fnb-backend/internal/utils"
)

type TableHandler struct {
	service services.TableService
}

func NewTableHandler(service services.TableService) *TableHandler {
	return &TableHandler{service: service}
}

// ResolveTable dipanggil saat customer scan QR meja.
// GET /api/v1/table/{qr_token}
//
// Alur:
//  1. Cari table berdasarkan qr_token.
//  2. Jika ada session aktif yang belum expired, reuse session itu.
//  3. Jika tidak ada, buat table_session baru.
func (h *TableHandler) ResolveTable(w http.ResponseWriter, r *http.Request) {
	qrToken := chi.URLParam(r, "qr_token")

	result, err := h.service.ResolveTable(qrToken)
	if err != nil {
		utils.Error(w, http.StatusNotFound, "meja tidak ditemukan, QR code mungkin tidak valid")
		return
	}

	utils.Success(w, http.StatusOK, result)
}

// ListTables — admin: daftar semua meja di satu restoran.
// GET /api/v1/staff/tables?restaurant_id=
func (h *TableHandler) ListTables(w http.ResponseWriter, r *http.Request) {
	restaurantID := r.URL.Query().Get("restaurant_id")

	result, err := h.service.ListTables(restaurantID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "gagal mengambil daftar meja: "+err.Error())
		return
	}

	utils.Success(w, http.StatusOK, result)
}

// CreateTable — admin: tambah meja baru + generate qr_token unik.
// POST /api/v1/staff/tables
func (h *TableHandler) CreateTable(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTableRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.Error(w, http.StatusBadRequest, "body request tidak valid")
		return
	}

	if req.RestaurantID == "" || req.TableNumber == "" {
		utils.Error(w, http.StatusBadRequest, "restaurant_id dan table_number wajib diisi")
		return
	}

	result, err := h.service.CreateTable(&req)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "gagal membuat meja: "+err.Error())
		return
	}

	utils.Success(w, http.StatusCreated, result)
}
