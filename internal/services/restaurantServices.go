package services

import (
	"github.com/yourorg/pos-fnb-backend/internal/dto"
	"github.com/yourorg/pos-fnb-backend/internal/models"
	"github.com/yourorg/pos-fnb-backend/internal/repository"
)

type RestaurantService interface {
	CreateRestaurant(req *dto.RestaurantRequest) (*dto.RestaurantResponse, error)
	// GetRestaurant(restaurantID string) (*dto.RestaurantResponse, error)
	// UpdateRestaurant(restaurantID string, req *dto.UpdateRestaurantRequest) (*dto.RestaurantResponse, error)
	// DeleteRestaurant(restaurantID string) error
}

type restaurantService struct {
	restaurantRepo repository.RestaurantRepository
}

func NewRestaurantService(
	restaurantRepo repository.RestaurantRepository,
) RestaurantService {
	return &restaurantService{
		restaurantRepo: restaurantRepo,
	}
}

func (s *restaurantService) CreateRestaurant(req *dto.RestaurantRequest) (*dto.RestaurantResponse, error) {
	restaurant := &models.Restaurant{
		Name:    req.Name,
		Address: req.Address,
		LogoURL: req.LogoURL,
	}

	createdRestaurant, err := s.restaurantRepo.Create(restaurant)
	if err != nil {
		return nil, err
	}

	return &dto.RestaurantResponse{
		ID:      createdRestaurant.ID.String(),
		Name:    createdRestaurant.Name,
		Address: createdRestaurant.Address,
		LogoURL: createdRestaurant.LogoURL,
	}, nil
}
