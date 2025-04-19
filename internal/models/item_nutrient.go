package models

type ItemNutrient struct {
	ItemID              uint    `gorm:"primaryKey" json:"-"`
	Name                string  `gorm:"primaryKey" json:"name"`
	Amount              float64 `json:"amount"`
	Unit                string  `gorm:"type:varchar(20)" json:"unit"`
	PercentOfDailyNeeds float64 `json:"percentOfDailyNeeds"`
}
