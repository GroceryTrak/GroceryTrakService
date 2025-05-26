package models

type User struct {
	ID                string  `gorm:"primaryKey;type:varchar(255);not null" json:"id"`
	Name              string  `gorm:"type:varchar(100)" json:"name"`
	Gender            string  `gorm:"type:varchar(20)" json:"gender"`
	Age               int     `gorm:"type:int" json:"age"`
	Height            float64 `gorm:"type:decimal(5,2)" json:"height"`
	HeightUnit        string  `gorm:"type:varchar(10)" json:"heightUnit"`
	CurrentWeight     float64 `gorm:"type:decimal(5,2)" json:"currentWeight"`
	CurrentWeightUnit string  `gorm:"type:varchar(10)" json:"currentWeightUnit"`
	TargetWeight      float64 `gorm:"type:decimal(5,2)" json:"targetWeight"`
	TargetWeightUnit  string  `gorm:"type:varchar(10)" json:"targetWeightUnit"`
	ActivityLevel     string  `gorm:"type:varchar(50)" json:"activityLevel"`
	DietType          string  `gorm:"type:varchar(50)" json:"dietType"`
	Goal              string  `gorm:"type:varchar(50)" json:"goal"`
}
