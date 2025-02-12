package models

type Item struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
	Unit        string `json:"unit"`
}
