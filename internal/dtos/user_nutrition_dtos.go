package dtos

type UserNutritionResponse struct {
	UserID         string  `json:"user_id"`
	Calories       int     `json:"calories"`
	CarbPercent    float32 `json:"carb_percent"`
	ProteinPercent float32 `json:"protein_percent"`
	FatPercent     float32 `json:"fat_percent"`
}
