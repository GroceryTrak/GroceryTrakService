package dtos

type NutrientRequest struct {
	Name                string  `json:"name" example:"Calories"`
	Amount              float64 `json:"amount" example:"200"`
	Unit                string  `json:"unit" example:"kcal"`
	PercentOfDailyNeeds float64 `json:"percentOfDailyNeeds" example:"100"`
}

type NutrientResponse struct {
	Name                string  `json:"name" example:"Calories"`
	Amount              float64 `json:"amount" example:"200"`
	Unit                string  `json:"unit" example:"kcal"`
	PercentOfDailyNeeds float64 `json:"percentOfDailyNeeds" example:"100"`
}

type NutrientsResponse struct {
	Nutrients []NutrientResponse `json:"nutrients"`
}
