package models

type UserItem struct {
	UserID uint    `gorm:"primaryKey;constraint:OnDelete:CASCADE;" json:"user_id"`
	ItemID uint    `gorm:"primaryKey;constraint:OnDelete:CASCADE;" json:"item_id"`
	Amount float32 `json:"amount"`
	Unit   string  `gorm:"type:varchar(20)" json:"unit"`
}
