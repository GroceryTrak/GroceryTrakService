package models

type Item struct {
	ID            uint           `gorm:"type:serial;primaryKey" json:"id"`
	Name          string         `gorm:"type:varchar(255);not null" json:"name"`
	Image         string         `json:"image"`
	SpoonacularID uint           `json:"spoonacular_id"`
	Nutrients     []ItemNutrient `gorm:"foreignKey:ItemID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"nutrients"`
}
