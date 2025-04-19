package models

type RecipeInstruction struct {
	RecipeID uint   `gorm:"primaryKey" json:"-"`
	Number   uint   `gorm:"primaryKey" json:"number"`
	Step     string `gorm:"type:text" json:"step"`
}
