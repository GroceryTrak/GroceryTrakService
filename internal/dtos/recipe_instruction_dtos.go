package dtos

type RecipeInstructionRequest struct {
	Number uint   `json:"number" example:"1"`
	Step   string `json:"step" example:"Bring a large pot of salted water to boil"`
}

type RecipeInstructionResponse struct {
	Number uint   `json:"number" example:"1"`
	Step   string `json:"step" example:"Bring a large pot of salted water to boil"`
}
