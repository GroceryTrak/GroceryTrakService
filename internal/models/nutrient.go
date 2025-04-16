package models

type Nutrient struct {
	ID                  uint    `gorm:"serial;primaryKey" json:"id"`
	ItemID              uint    `json:"item_id"`
	Name                string  `gorm:"type:varchar(100);not null" json:"name"`
	Amount              float64 `json:"amount"`
	Unit                string  `gorm:"type:varchar(20)" json:"unit"`
	PercentOfDailyNeeds float64 `json:"percentOfDailyNeeds"`
}
