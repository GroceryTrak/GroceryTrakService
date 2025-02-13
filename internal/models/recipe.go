package models

type Recipe struct {
	ID          uint         `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string       `gorm:"not null" json:"name"`
	Instruction string       `json:"instruction"`
	Difficulty  string       `json:"difficulty"`
	Duration    uint         `json:"duration"`
	KCal        uint         `json:"kcal"`
	Ingredients []RecipeItem `gorm:"foreignKey:RecipeID;constraint:OnDelete:CASCADE;" json:"ingredients"`
}
