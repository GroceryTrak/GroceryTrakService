package models

type Role string

const (
	UserRole  Role = "user"
	AdminRole Role = "admin"
)

type User struct {
	ID       uint   `gorm:"type:serial;primaryKey" json:"id"`
	Username string `gorm:"type:varchar(50);unique;not null" json:"username"`
	Password string `gorm:"type:varchar(255);not null" json:"password"`
	Role     Role   `gorm:"type:role;not null;default:'user'" json:"role"`
}
