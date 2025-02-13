package models

type RecipeItem struct {
	RecipeID uint   `gorm:"primaryKey;constraint:OnDelete:CASCADE;" json:"recipe_id"`
	ItemID   uint   `gorm:"primaryKey;constraint:OnDelete:CASCADE;" json:"item_id"`
	Item     Item   `gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE;" json:"item"`
	Quantity uint   `json:"quantity"`
	Unit     string `json:"unit"`
}
