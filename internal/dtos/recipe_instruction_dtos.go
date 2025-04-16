package dtos

type RecipeInstructionRequest struct {
	Number int    `json:"number" example:"1"`
	Step   string `json:"step" example:"Bring a large pot of salted water to boil"`
}

type RecipeInstructionResponse struct {
	Number int    `json:"number" example:"1"`
	Step   string `json:"step" example:"Bring a large pot of salted water to boil"`
}
