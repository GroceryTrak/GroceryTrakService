package models

type RecipeNutrient struct {
	ID                  uint    `gorm:"serial;primaryKey" json:"id"`
	RecipeID            uint    `json:"recipe_id"`
	Name                string  `gorm:"type:varchar(100);not null" json:"name"`
	Amount              float64 `json:"amount"`
	Unit                string  `gorm:"type:varchar(20)" json:"unit"`
	PercentOfDailyNeeds float64 `json:"percentOfDailyNeeds"`

	Recipe Recipe `gorm:"foreignKey:RecipeID;references:ID"`
}
