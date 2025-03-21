package models

type Recipe struct {
	ID          uint         `gorm:"type:serial;primaryKey" json:"id"`
	Title       string       `gorm:"type:varchar(255);not null" json:"title"`
	ReadyTime   int16        `json:"ready_time"`
	CookingTime int16        `json:"cooking_time"`
	PrepTime    int16        `json:"prep_time"`
	Image       string       `json:"image"`
	KCal        float32      `json:"kcal"`
	Vegan       bool         `json:"vegan"`
	Vegetarian  bool         `json:"vegetarian"`
	Ingredients []RecipeItem `gorm:"foreignKey:RecipeID;constraint:OnDelete:CASCADE;" json:"ingredients"`
}
