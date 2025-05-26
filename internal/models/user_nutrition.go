package models

type UserNutrition struct {
	UserID         string  `gorm:"primaryKey;type:varchar(255);not null" json:"user_id"`
	Calories       int     `gorm:"type:int;not null" json:"calories"`
	CarbPercent    float32 `gorm:"type:decimal(5,2);not null" json:"carb_percent"`
	ProteinPercent float32 `gorm:"type:decimal(5,2);not null" json:"protein_percent"`
	FatPercent     float32 `gorm:"type:decimal(5,2);not null" json:"fat_percent"`
}
