package models

type RecipeInstruction struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	RecipeID uint   `json:"recipe_id"`
	Number   int    `json:"number"`
	Step     string `gorm:"type:text" json:"step"`

	Recipe Recipe `gorm:"foreignKey:RecipeID;references:ID"`
}
