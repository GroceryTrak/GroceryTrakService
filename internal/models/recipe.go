package models

import "errors"

type Difficulty string

const (
	Easy   Difficulty = "easy"
	Medium Difficulty = "medium"
	Hard   Difficulty = "hard"
)

func ParseDifficulty(value string) (Difficulty, error) {
	switch value {
	case string(Easy):
		return Easy, nil
	case string(Medium):
		return Medium, nil
	case string(Hard):
		return Hard, nil
	default:
		return "", errors.New("invalid difficulty level")
	}
}

type Recipe struct {
	ID          uint         `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string       `gorm:"type:varchar(255);not null" json:"name"`
	Instruction string       `json:"instruction"`
	Difficulty  Difficulty   `gorm:"type:difficulty" json:"difficulty"`
	Duration    uint         `json:"duration"`
	KCal        uint         `json:"kcal"`
	Diet        string       `json:"diet"`
	Ingredients []RecipeItem `gorm:"foreignKey:RecipeID;constraint:OnDelete:CASCADE;" json:"ingredients"`
}
