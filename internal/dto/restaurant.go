package dto

type RestaurantRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	LogoURL string `json:"logo_url"`
}

type RestaurantResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	LogoURL string `json:"logo_url"`
}
