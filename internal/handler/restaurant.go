package handler

import (
	"net/http"

	"github.com/yourorg/pos-fnb-backend/internal/dto"
	"github.com/yourorg/pos-fnb-backend/internal/services"
	"github.com/yourorg/pos-fnb-backend/internal/utils"
)

type RestaurantHandler struct {
	restaurantService services.RestaurantService
}

func NewRestaurantHandler(restaurantService services.RestaurantService) *RestaurantHandler {
	return &RestaurantHandler{
		restaurantService: restaurantService,
	}
}

func (h *RestaurantHandler) CreateRestaurant(w http.ResponseWriter, r *http.Request) {
	var req dto.RestaurantRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.Error(w, http.StatusBadRequest, "body request tidak valid")
		return
	}
	createdRestaurant, err := h.restaurantService.CreateRestaurant(&req)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(w, http.StatusCreated, createdRestaurant)
}
