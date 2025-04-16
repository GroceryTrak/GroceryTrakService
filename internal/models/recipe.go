package models

type Recipe struct {
	ID            uint                `gorm:"type:serial;primaryKey" json:"id"`
	SpoonacularID uint                `json:"spoonacular_id"`
	Title         string              `gorm:"type:varchar(255);not null" json:"title"`
	Summary       string              `gorm:"type:text" json:"summary"`
	Instructions  []RecipeInstruction `gorm:"foreignKey:RecipeID;constraint:OnDelete:CASCADE" json:"instructions"`
	Servings      float32             `json:"servings"`
	ReadyTime     int16               `json:"ready_time"`
	CookingTime   int16               `json:"cooking_time"`
	PrepTime      int16               `json:"prep_time"`
	Image         string              `json:"image"`
	KCal          float32             `json:"kcal"`
	Vegan         bool                `json:"vegan"`
	Vegetarian    bool                `json:"vegetarian"`
	Ingredients   []RecipeItem        `gorm:"foreignKey:RecipeID;constraint:OnDelete:CASCADE;" json:"ingredients"`
	Nutrients     []RecipeNutrient    `gorm:"foreignKey:RecipeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"nutrients"`
}
