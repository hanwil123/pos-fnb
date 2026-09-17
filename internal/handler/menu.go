package handler

import (
	"net/http"

	"github.com/yourorg/pos-fnb-backend/internal/dto"
	"github.com/yourorg/pos-fnb-backend/internal/services"
	"github.com/yourorg/pos-fnb-backend/internal/utils"
)

type MenuHandler struct {
	service services.MenuService
}

func NewMenuHandler(service services.MenuService) *MenuHandler {
	return &MenuHandler{service: service}
}

// GetMenu — customer: ambil semua kategori + item menu untuk satu restoran.
// GET /api/v1/menu?restaurant_id=
func (h *MenuHandler) GetMenu(w http.ResponseWriter, r *http.Request) {
	restaurantID := r.URL.Query().Get("restaurant_id")
	if restaurantID == "" {
		utils.Error(w, http.StatusBadRequest, "restaurant_id wajib diisi")
		return
	}

	result, err := h.service.GetMenu(restaurantID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "gagal mengambil menu: "+err.Error())
		return
	}

	utils.Success(w, http.StatusOK, result)
}

// GetRecommendations — placeholder untuk fitur "menu signature & trending".
// Untuk MVP: return item dengan is_signature = true.
// Nanti di Phase 3 diganti/ditambah query ke materialized view menu_trending
// + pemanggilan AI service untuk personalisasi.
// GET /api/v1/menu/recommendations?restaurant_id=
func (h *MenuHandler) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	restaurantID := r.URL.Query().Get("restaurant_id")
	if restaurantID == "" {
		utils.Error(w, http.StatusBadRequest, "restaurant_id wajib diisi")
		return
	}

	result, err := h.service.GetRecommendations(restaurantID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "gagal mengambil rekomendasi: "+err.Error())
		return
	}

	utils.Success(w, http.StatusOK, result)
}

// CreateMenuItem — admin: tambah item menu baru.
// POST /api/v1/staff/menu-items
func (h *MenuHandler) CreateMenuItem(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateMenuItemRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.Error(w, http.StatusBadRequest, "body request tidak valid")
		return
	}

	if req.RestaurantID == "" || req.CategoryID == "" || req.Name == "" || req.Price <= 0 {
		utils.Error(w, http.StatusBadRequest, "restaurant_id, category_id, name, dan price wajib diisi")
		return
	}

	result, err := h.service.CreateMenuItem(&req)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "gagal membuat item menu: "+err.Error())
		return
	}

	utils.Success(w, http.StatusCreated, result)
}
