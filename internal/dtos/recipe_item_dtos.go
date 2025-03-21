package dtos

type RecipeItemRequest struct {
	ItemID uint    `json:"item_id" example:"456"`
	Amount float32 `json:"amount" example:"2.5"`
	Unit   string  `json:"unit" example:"cups"`
}

type RecipeItemResponse struct {
	Item   ItemResponse `json:"item"`
	Amount float32      `json:"amount" example:"2.5"`
	Unit   string       `json:"unit" example:"cups"`
}

type RecipeItemsResponse struct {
	RecipeItems []RecipeItemResponse `json:"recipe_items"`
}
