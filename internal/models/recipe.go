package models

type Recipe struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string `json:"name"`
	Instruction string `json:"instruction"`
	Difficulty  string `json:"difficulty"`
	Duration    uint   `json:"duration"`
	KCal        uint   `json:"kcal"`
	Ingredients []Item `gorm:"many2many:recipes_items;" json:"ingredients"`
}
