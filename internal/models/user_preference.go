package models

type UserPreference struct {
	UserID uint   `gorm:"primaryKey" json:"user_id"`
	Diet   string `gorm:"type:varchar(20)" json:"diet"`
}
