package templates

type ItemRequest struct {
	Name        string `json:"name" example:"Milk"`
	Description string `json:"description" example:"Organic whole milk"`
}

type ItemResponse struct {
	ID          uint   `json:"id" example:"1"`
	Name        string `json:"name" example:"Milk"`
	Description string `json:"description" example:"Organic whole milk"`
}

type ItemsResponse struct {
	Items []ItemResponse `json:"items"`
}
