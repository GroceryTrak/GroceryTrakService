package dtos

type ItemRequest struct {
	ID            uint              `json:"id" example:"1"`
	Name          string            `json:"name" example:"Milk"`
	Image         string            `json:"image" example:"milk.jpg"`
	SpoonacularID uint              `json:"spoonacular_id" example:"1"`
	Nutrients     []NutrientRequest `json:"nutrients"`
}

type ItemResponse struct {
	ID            uint               `json:"id" example:"1"`
	Name          string             `json:"name" example:"Milk"`
	Image         string             `json:"image" example:"milk.jpg"`
	SpoonacularID uint               `json:"spoonacular_id" example:"1"`
	Nutrients     []NutrientResponse `json:"nutrients"`
}

type ItemsResponse struct {
	Items []ItemResponse `json:"items"`
}

type ItemQuery struct {
	Name string `json:"name" example:"pasta"`
}
