package dtos

type ItemNutrientRequest struct {
	Name                string  `json:"name" example:"Calories"`
	Amount              float64 `json:"amount" example:"200"`
	Unit                string  `json:"unit" example:"kcal"`
	PercentOfDailyNeeds float64 `json:"percentOfDailyNeeds" example:"100"`
}

type ItemNutrientResponse struct {
	Name                string  `json:"name" example:"Calories"`
	Amount              float64 `json:"amount" example:"200"`
	Unit                string  `json:"unit" example:"kcal"`
	PercentOfDailyNeeds float64 `json:"percentOfDailyNeeds" example:"100"`
}

type ItemNutrientsResponse struct {
	Nutrients []ItemNutrientResponse `json:"nutrients"`
}
