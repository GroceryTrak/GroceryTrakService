package dtos

type RecipeNutrientRequest struct {
	Name                string  `json:"name" example:"Calories"`
	Amount              float64 `json:"amount" example:"200"`
	Unit                string  `json:"unit" example:"kcal"`
	PercentOfDailyNeeds float64 `json:"percentOfDailyNeeds" example:"100"`
}

type RecipeNutrientResponse struct {
	Name                string  `json:"name" example:"Calories"`
	Amount              float64 `json:"amount" example:"200"`
	Unit                string  `json:"unit" example:"kcal"`
	PercentOfDailyNeeds float64 `json:"percentOfDailyNeeds" example:"100"`
}

type RecipeNutrientsResponse struct {
	Nutrients []RecipeNutrientResponse `json:"nutrients"`
}
