package repository

import (
	"github.com/GroceryTrak/GroceryTrakService/internal/dtos"
	"github.com/GroceryTrak/GroceryTrakService/internal/models"
	"gorm.io/gorm"
)

type UserNutritionRepository interface {
	CreateUserNutrition(userID string) (dtos.UserNutritionResponse, error)
	GetUserNutrition(userID string) (dtos.UserNutritionResponse, error)
	UpdateUserNutrition(userID string) (dtos.UserNutritionResponse, error)
	DeleteUserNutrition(userID string) error
}

type UserNutritionRepositoryImpl struct {
	db *gorm.DB
}

func NewUserNutritionRepository(db *gorm.DB) UserNutritionRepository {
	return &UserNutritionRepositoryImpl{
		db: db,
	}
}

func (r *UserNutritionRepositoryImpl) CreateUserNutrition(userID string) (dtos.UserNutritionResponse, error) {
	// Get user data
	var user models.User
	if err := r.db.First(&user, "id = ?", userID).Error; err != nil {
		return dtos.UserNutritionResponse{}, err
	}

	// Convert units to metric if needed
	weightKg := user.CurrentWeight
	if user.CurrentWeightUnit == "lb" {
		weightKg = user.CurrentWeight * 0.453592
	}

	heightCm := user.Height
	if user.HeightUnit == "inch" {
		heightCm = user.Height * 2.54
	}

	// Calculate BMR using Mifflin-St Jeor Equation
	var bmr float64
	if user.Gender == "male" {
		bmr = 10*weightKg + 6.25*heightCm - 5*float64(user.Age) + 5
	} else {
		bmr = 10*weightKg + 6.25*heightCm - 5*float64(user.Age) - 161
	}

	// Apply activity multiplier
	var activityMultiplier float64
	switch user.ActivityLevel {
	case "sedentary":
		activityMultiplier = 1.2
	case "lightly_active":
		activityMultiplier = 1.375
	case "moderately_active":
		activityMultiplier = 1.55
	case "very_active":
		activityMultiplier = 1.725
	case "super_active":
		activityMultiplier = 1.9
	default:
		activityMultiplier = 1.2 // Default to sedentary if unknown
	}

	tdee := bmr * activityMultiplier

	// Adjust for goal
	switch user.Goal {
	case "lose_weight_slow":
		tdee -= 250
	case "lose_weight_standard":
		tdee -= 500
	case "lose_weight_fast":
		tdee -= 750
	case "gain_weight_slow":
		tdee += 250
	case "gain_weight_standard":
		tdee += 500
	case "gain_weight_fast":
		tdee += 750
	case "gain_muscle_reduce_fat":
		tdee += 250 // Slight surplus for muscle gain
	case "maintain_weight":
		// No adjustment needed
	default:
		// Default to maintain weight
	}

	// Set macronutrient ratios based on diet type
	var carbPercent, proteinPercent, fatPercent float32
	switch user.DietType {
	case "balanced":
		carbPercent = 50
		proteinPercent = 20
		fatPercent = 30
	case "high_protein":
		carbPercent = 30
		proteinPercent = 40
		fatPercent = 30
	case "low_carb":
		carbPercent = 20
		proteinPercent = 40
		fatPercent = 40
	case "vegetarian":
		carbPercent = 55
		proteinPercent = 15
		fatPercent = 30
	case "vegan":
		carbPercent = 60
		proteinPercent = 15
		fatPercent = 25
	case "keto":
		carbPercent = 5
		proteinPercent = 20
		fatPercent = 75
	case "mediterranean":
		carbPercent = 45
		proteinPercent = 25
		fatPercent = 30
	default:
		carbPercent = 50
		proteinPercent = 20
		fatPercent = 30
	}

	// Create and save user nutrition record
	userNutrition := models.UserNutrition{
		UserID:         userID,
		Calories:       int(tdee),
		CarbPercent:    carbPercent,
		ProteinPercent: proteinPercent,
		FatPercent:     fatPercent,
	}

	if err := r.db.Create(&userNutrition).Error; err != nil {
		return dtos.UserNutritionResponse{}, err
	}

	return dtos.UserNutritionResponse{
		UserID:         userNutrition.UserID,
		Calories:       userNutrition.Calories,
		CarbPercent:    userNutrition.CarbPercent,
		ProteinPercent: userNutrition.ProteinPercent,
		FatPercent:     userNutrition.FatPercent,
	}, nil
}

func (r *UserNutritionRepositoryImpl) GetUserNutrition(userID string) (dtos.UserNutritionResponse, error) {
	var userNutrition models.UserNutrition
	if err := r.db.First(&userNutrition, "user_id = ?", userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// If not found, create new nutrition record
			return r.CreateUserNutrition(userID)
		}
		return dtos.UserNutritionResponse{}, err
	}

	return dtos.UserNutritionResponse{
		UserID:         userNutrition.UserID,
		Calories:       userNutrition.Calories,
		CarbPercent:    userNutrition.CarbPercent,
		ProteinPercent: userNutrition.ProteinPercent,
		FatPercent:     userNutrition.FatPercent,
	}, nil
}

func (r *UserNutritionRepositoryImpl) UpdateUserNutrition(userID string) (dtos.UserNutritionResponse, error) {
	// Delete existing nutrition record
	if err := r.DeleteUserNutrition(userID); err != nil {
		return dtos.UserNutritionResponse{}, err
	}

	// Create new nutrition record
	return r.CreateUserNutrition(userID)
}

func (r *UserNutritionRepositoryImpl) DeleteUserNutrition(userID string) error {
	return r.db.Delete(&models.UserNutrition{}, "user_id = ?", userID).Error
}
