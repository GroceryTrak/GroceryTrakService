package models

type RecipeNutrient struct {
	RecipeID            uint    `gorm:"primaryKey" json:"-"`
	Name                string  `gorm:"primaryKey" json:"name"`
	Amount              float64 `json:"amount"`
	Unit                string  `json:"unit"`
	PercentOfDailyNeeds float64 `json:"percent_of_daily_needs"`
}
