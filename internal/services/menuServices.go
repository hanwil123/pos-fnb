package services

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/yourorg/pos-fnb-backend/internal/dto"
	"github.com/yourorg/pos-fnb-backend/internal/models"
	"github.com/yourorg/pos-fnb-backend/internal/repository"
)

type MenuService interface {
	GetMenu(restaurantID string) ([]dto.MenuCategoryResponse, error)
	GetRecommendations(restaurantID string) (*dto.RecommendationsResponse, error)
	CreateMenuItem(req *dto.CreateMenuItemRequest) (*dto.MenuItemResponse, error)
}

type menuService struct {
	repo repository.MenuRepository
}

func NewMenuService(repo repository.MenuRepository) MenuService {
	return &menuService{repo: repo}
}

func (s *menuService) GetMenu(restaurantID string) ([]dto.MenuCategoryResponse, error) {
	restaurantUUID, err := uuid.Parse(restaurantID)
	if err != nil {
		return nil, err
	}

	categories, err := s.repo.GetCategoriesByRestaurantID(restaurantUUID)
	if err != nil {
		return nil, err
	}

	items, err := s.repo.GetAvailableItemsByRestaurantID(restaurantUUID)
	if err != nil {
		return nil, err
	}

	// Kelompokkan item per kategori
	grouped := make(map[string][]dto.MenuItemResponse)
	for _, item := range items {
		key := item.CategoryID.String()
		grouped[key] = append(grouped[key], dto.MenuItemResponse{
			ID:           item.ID.String(),
			RestaurantID: item.RestaurantID.String(),
			CategoryID:   item.CategoryID.String(),
			Name:         item.Name,
			Description:  item.Description,
			Price:        item.Price,
			ImageURL:     item.ImageURL,
			IsSignature:  item.IsSignature,
			IsAvailable:  item.IsAvailable,
			Tags:         item.Tags,
		})
	}

	result := make([]dto.MenuCategoryResponse, 0, len(categories))
	for _, cat := range categories {
		result = append(result, dto.MenuCategoryResponse{
			ID:           cat.ID.String(),
			RestaurantID: cat.RestaurantID.String(),
			Name:         cat.Name,
			SortOrder:    cat.SortOrder,
			Items:        grouped[cat.ID.String()],
		})
	}

	return result, nil
}

func (s *menuService) GetRecommendations(restaurantID string) (*dto.RecommendationsResponse, error) {
	restaurantUUID, err := uuid.Parse(restaurantID)
	if err != nil {
		return nil, err
	}

	signatureItems, err := s.repo.GetSignatureItems(restaurantUUID)
	if err != nil {
		return nil, err
	}

	signature := make([]dto.MenuItemResponse, 0, len(signatureItems))
	for _, item := range signatureItems {
		signature = append(signature, dto.MenuItemResponse{
			ID:           item.ID.String(),
			RestaurantID: item.RestaurantID.String(),
			CategoryID:   item.CategoryID.String(),
			Name:         item.Name,
			Description:  item.Description,
			Price:        item.Price,
			ImageURL:     item.ImageURL,
			IsSignature:  item.IsSignature,
			IsAvailable:  item.IsAvailable,
			Tags:         item.Tags,
		})
	}

	return &dto.RecommendationsResponse{
		Signature: signature,
		Trending:  []dto.MenuItemResponse{}, // TODO: isi dari menu_trending materialized view
	}, nil
}

func (s *menuService) CreateMenuItem(req *dto.CreateMenuItemRequest) (*dto.MenuItemResponse, error) {
	restaurantID, err := uuid.Parse(req.RestaurantID)
	if err != nil {
		return nil, err
	}

	categoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return nil, err
	}

	item := &models.MenuItem{
		RestaurantID: restaurantID,
		CategoryID:   categoryID,
		Name:         req.Name,
		Description:  req.Description,
		Price:        req.Price,
		ImageURL:     req.ImageURL,
		IsSignature:  req.IsSignature,
		IsAvailable:  true,
		Tags:         pq.StringArray(req.Tags),
	}

	if err := s.repo.CreateItem(item); err != nil {
		return nil, err
	}

	return &dto.MenuItemResponse{
		ID:           item.ID.String(),
		RestaurantID: item.RestaurantID.String(),
		CategoryID:   item.CategoryID.String(),
		Name:         item.Name,
		Description:  item.Description,
		Price:        item.Price,
		ImageURL:     item.ImageURL,
		IsSignature:  item.IsSignature,
		IsAvailable:  item.IsAvailable,
		Tags:         item.Tags,
	}, nil
}

