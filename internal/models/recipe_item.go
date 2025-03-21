package models

type RecipeItem struct {
	RecipeID uint    `gorm:"primaryKey;constraint:OnDelete:CASCADE;" json:"recipe_id"`
	ItemID   uint    `gorm:"primaryKey;constraint:OnDelete:CASCADE;" json:"item_id"`
	Amount   float32 `json:"amount"`
	Unit     string  `gorm:"type:varchar(20)" json:"unit"`

	Recipe Recipe `gorm:"foreignKey:RecipeID;references:ID"`
	Item   Item   `gorm:"foreignKey:ItemID;references:ID"`
}
