package dtos

type RecipeRequest struct {
	Title         string                     `json:"title" example:"Spaghetti Carbonara"`
	Summary       string                     `json:"summary" example:"A classic Italian pasta dish with eggs, cheese, pancetta, and black pepper."`
	SpoonacularID uint                       `json:"spoonacular_id" example:"12345"`
	Instructions  []RecipeInstructionRequest `json:"instructions"`
	Servings      float32                    `json:"servings" example:"4"`
	ReadyTime     int16                      `json:"ready_time" example:"30"`
	CookingTime   int16                      `json:"cooking_time" example:"20"`
	PrepTime      int16                      `json:"prep_time" example:"10"`
	Image         string                     `json:"image" example:"https://example.com/spaghetti.jpg"`
	KCal          float32                    `json:"kcal" example:"450.5"`
	Vegan         bool                       `json:"vegan" example:"false"`
	Vegetarian    bool                       `json:"vegetarian" example:"false"`
	Ingredients   []RecipeItemRequest        `json:"ingredients"`
	Nutrients     []RecipeNutrientRequest    `json:"nutrients"`
}

type RecipeResponse struct {
	ID            uint                        `json:"id" example:"1"`
	Title         string                      `json:"title" example:"Spaghetti Carbonara"`
	Summary       string                      `json:"summary" example:"A classic Italian pasta dish with eggs, cheese, pancetta, and black pepper."`
	SpoonacularID uint                        `json:"spoonacular_id" example:"12345"`
	Instructions  []RecipeInstructionResponse `json:"instructions"`
	Servings      float32                     `json:"servings" example:"4"`
	ReadyTime     int16                       `json:"ready_time" example:"30"`
	CookingTime   int16                       `json:"cooking_time" example:"20"`
	PrepTime      int16                       `json:"prep_time" example:"10"`
	Image         string                      `json:"image" example:"https://example.com/spaghetti.jpg"`
	KCal          float32                     `json:"kcal" example:"450.5"`
	Vegan         bool                        `json:"vegan" example:"false"`
	Vegetarian    bool                        `json:"vegetarian" example:"false"`
	Ingredients   []RecipeItemResponse        `json:"ingredients"`
	Nutrients     []RecipeNutrientResponse    `json:"nutrients"`
}

type DietCount struct {
	Vegan      bool  `json:"vegan"`
	Vegetarian bool  `json:"vegetarian"`
	Count      int64 `json:"count"`
}

type RecipesResponse struct {
	Recipes    []RecipeResponse `json:"recipes"`
	Count      int              `json:"count"`
	DietCounts []DietCount      `json:"diet_counts"`
}

type RecipeQuery struct {
	Title       string   `json:"title" example:"pasta"`
	Ingredients []string `json:"ingredients" example:"1,2,3"`
	Diet        string   `json:"diet" example:"vegan"`
}
