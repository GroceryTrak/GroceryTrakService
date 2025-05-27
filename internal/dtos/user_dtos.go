package dtos

type UserResponse struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	Gender            string  `json:"gender"`
	Height            float64 `json:"height"`
	HeightUnit        string  `json:"height_unit"`
	CurrentWeight     float64 `json:"current_weight"`
	CurrentWeightUnit string  `json:"current_weight_unit"`
	TargetWeight      float64 `json:"target_weight"`
	TargetWeightUnit  string  `json:"target_weight_unit"`
	ActivityLevel     string  `json:"activity_level"`
	DietType          string  `json:"diet_type"`
	Goal              string  `json:"goal"`
}

type UserRequest struct {
	Name              string  `json:"name" validate:"required"`
	Gender            string  `json:"gender" validate:"required"`
	Age               int     `json:"age" validate:"required"`
	Height            float64 `json:"height" validate:"required"`
	HeightUnit        string  `json:"height_unit" validate:"required"`
	CurrentWeight     float64 `json:"current_weight" validate:"required"`
	CurrentWeightUnit string  `json:"current_weight_unit" validate:"required"`
	TargetWeight      float64 `json:"target_weight" validate:"required"`
	TargetWeightUnit  string  `json:"target_weight_unit" validate:"required"`
	ActivityLevel     string  `json:"activity_level" validate:"required"`
	DietType          string  `json:"diet_type" validate:"required"`
	Goal              string  `json:"goal" validate:"required"`
}
