package repository

import (
	"github.com/GroceryTrak/GroceryTrakService/internal/dtos"
	"github.com/GroceryTrak/GroceryTrakService/internal/models"
	"gorm.io/gorm"
)

type UserNutritionRepository interface {
	CalculateAndSaveNutrition(userID string) (dtos.UserNutritionResponse, error)
	GetUserNutrition(userID string) (dtos.UserNutritionResponse, error)
}

type UserNutritionRepositoryImpl struct {
	db *gorm.DB
}

func NewUserNutritionRepository(db *gorm.DB) UserNutritionRepository {
	return &UserNutritionRepositoryImpl{
		db: db,
	}
}

func (r *UserNutritionRepositoryImpl) CalculateAndSaveNutrition(userID string) (dtos.UserNutritionResponse, error) {
	// Get user data
	var user models.User
	if err := r.db.First(&user, "cognito_uid = ?", userID).Error; err != nil {
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
	case "extra_active":
		activityMultiplier = 1.9
	}

	tdee := bmr * activityMultiplier

	// Adjust for goal
	switch user.Goal {
	case "loss":
		tdee -= 500
	case "gain":
		tdee += 500
	}

	// Set macronutrient ratios based on diet type
	var carbPercent, proteinPercent, fatPercent float32
	switch user.DietType {
	case "balanced":
		carbPercent = 50
		proteinPercent = 20
		fatPercent = 30
	case "low_carb":
		carbPercent = 20
		proteinPercent = 40
		fatPercent = 40
	case "high_protein":
		carbPercent = 30
		proteinPercent = 40
		fatPercent = 30
	case "keto":
		carbPercent = 5
		proteinPercent = 20
		fatPercent = 75
	}

	// Create and save user nutrition record
	userNutrition := models.UserNutrition{
		UserID:         userID,
		Calories:       int(tdee),
		CarbPercent:    carbPercent,
		ProteinPercent: proteinPercent,
		FatPercent:     fatPercent,
	}

	if err := r.db.Save(&userNutrition).Error; err != nil {
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
