package dtos

type UserItemRequest struct {
	ItemID uint    `json:"item_id" example:"456"`
	Amount float32 `json:"amount" example:"2.0"`
	Unit   string  `json:"unit" example:"kg"`
}

type UserItemResponse struct {
	Item   ItemResponse `json:"item"`
	Amount float32      `json:"amount" example:"2.0"`
	Unit   string       `json:"unit" example:"kg"`
}

type UserItemsResponse struct {
	UserItems []UserItemResponse `json:"user_items"`
}

type UserItemQuery struct {
	Name string `json:"name" example:"pasta"`
}
