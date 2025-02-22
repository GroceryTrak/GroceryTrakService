package models

type UserItem struct {
	UserID   uint   `gorm:"primaryKey;constraint:OnDelete:CASCADE;" json:"user_id"`
	ItemID   uint   `gorm:"primaryKey;constraint:OnDelete:CASCADE;" json:"item_id"`
	Item     Item   `gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE;" json:"item"`
	Quantity uint   `json:"quantity"`
	Unit     string `gorm:"type:varchar(20)" json:"unit"`
}
