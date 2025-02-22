package models

type Item struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string `gorm:"type:varchar(255);not null" json:"name"`
	Description string `json:"description"`
}
