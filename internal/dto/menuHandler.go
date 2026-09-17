package dto

// ── Requests ──────────────────────────────────────────────────────────────────

type CreateMenuItemRequest struct {
	RestaurantID string   `json:"restaurant_id"`
	CategoryID   string   `json:"category_id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Price        float64  `json:"price"`
	ImageURL     string   `json:"image_url"`
	IsSignature  bool     `json:"is_signature"`
	Tags         []string `json:"tags"`
}

type UpdateMenuItemRequest struct {
	CategoryID  string   `json:"category_id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	ImageURL    string   `json:"image_url"`
	IsSignature bool     `json:"is_signature"`
	Tags        []string `json:"tags"`
}

type DeleteMenuItemRequest struct {
	ID string `json:"id"`
}

// ── Responses ─────────────────────────────────────────────────────────────────

type MenuItemResponse struct {
	ID           string   `json:"id"`
	RestaurantID string   `json:"restaurant_id"`
	CategoryID   string   `json:"category_id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Price        float64  `json:"price"`
	ImageURL     string   `json:"image_url"`
	IsSignature  bool     `json:"is_signature"`
	IsAvailable  bool     `json:"is_available"`
	Tags         []string `json:"tags"`
}

type MenuCategoryResponse struct {
	ID           string             `json:"id"`
	RestaurantID string             `json:"restaurant_id"`
	Name         string             `json:"name"`
	SortOrder    int                `json:"sort_order"`
	Items        []MenuItemResponse `json:"items"`
}

type RecommendationsResponse struct {
	Signature []MenuItemResponse `json:"signature"`
	Trending  []MenuItemResponse `json:"trending"`
}
