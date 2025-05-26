package dtos

type UserResponse struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	Gender            string  `json:"gender"`
	Height            float64 `json:"height"`
	HeightUnit        string  `json:"heightUnit"`
	CurrentWeight     float64 `json:"currentWeight"`
	CurrentWeightUnit string  `json:"currentWeightUnit"`
	TargetWeight      float64 `json:"targetWeight"`
	TargetWeightUnit  string  `json:"targetWeightUnit"`
	ActivityLevel     string  `json:"activityLevel"`
	DietType          string  `json:"dietType"`
	Goal              string  `json:"goal"`
}

type UpdateUserRequest struct {
	Name              string  `json:"name" validate:"required"`
	Gender            string  `json:"gender" validate:"required"`
	Height            float64 `json:"height" validate:"required"`
	HeightUnit        string  `json:"heightUnit" validate:"required"`
	CurrentWeight     float64 `json:"currentWeight" validate:"required"`
	CurrentWeightUnit string  `json:"currentWeightUnit" validate:"required"`
	TargetWeight      float64 `json:"targetWeight" validate:"required"`
	TargetWeightUnit  string  `json:"targetWeightUnit" validate:"required"`
	ActivityLevel     string  `json:"activityLevel" validate:"required"`
	DietType          string  `json:"dietType" validate:"required"`
	Goal              string  `json:"goal" validate:"required"`
}
