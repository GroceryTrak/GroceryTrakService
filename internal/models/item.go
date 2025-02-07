package models

type Item struct {
	ID          string `gorm:"primaryKey" json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
	Unit        string `json:"unit"`
}
