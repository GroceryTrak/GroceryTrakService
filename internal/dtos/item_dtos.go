package dtos

type ItemRequest struct {
	ID            uint   `json:"id" example:"1"`
	Name          string `json:"name" example:"Milk"`
	Image         string `json:"image" example:"milk.jpg"`
	SpoonacularID uint   `json:"spoonacular_id" example:"1"`
}

type ItemResponse struct {
	ID            uint   `json:"id" example:"1"`
	Name          string `json:"name" example:"Milk"`
	Image         string `json:"image" example:"milk.jpg"`
	SpoonacularID uint   `json:"spoonacular_id" example:"1"`
}

type ItemsResponse struct {
	Items []ItemResponse `json:"items"`
}

type ItemQuery struct {
	Name string `json:"name" example:"pasta"`
}
