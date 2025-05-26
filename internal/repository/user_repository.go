package repository

import (
	"github.com/GroceryTrak/GroceryTrakService/internal/dtos"
	"github.com/GroceryTrak/GroceryTrakService/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUser(userID string) (dtos.UserResponse, error)
	UpdateUser(req dtos.UpdateUserRequest, userID string) (dtos.UserResponse, error)
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{
		db: db,
	}
}

func (r *UserRepositoryImpl) GetUser(userID string) (dtos.UserResponse, error) {
	var user models.User
	if err := r.db.First(&user, "id = ?", userID).Error; err != nil {
		return dtos.UserResponse{}, err
	}

	return dtos.UserResponse{
		ID:                user.ID,
		Name:              user.Name,
		Gender:            user.Gender,
		Height:            user.Height,
		HeightUnit:        user.HeightUnit,
		CurrentWeight:     user.CurrentWeight,
		CurrentWeightUnit: user.CurrentWeightUnit,
		TargetWeight:      user.TargetWeight,
		TargetWeightUnit:  user.TargetWeightUnit,
		ActivityLevel:     user.ActivityLevel,
		DietType:          user.DietType,
		Goal:              user.Goal,
	}, nil
}

func (r *UserRepositoryImpl) UpdateUser(req dtos.UpdateUserRequest, userID string) (dtos.UserResponse, error) {
	var user models.User
	if err := r.db.First(&user, "id = ?", userID).Error; err != nil {
		return dtos.UserResponse{}, err
	}

	user.Name = req.Name
	user.Gender = req.Gender
	user.Height = req.Height
	user.HeightUnit = req.HeightUnit
	user.CurrentWeight = req.CurrentWeight
	user.CurrentWeightUnit = req.CurrentWeightUnit
	user.TargetWeight = req.TargetWeight
	user.TargetWeightUnit = req.TargetWeightUnit
	user.ActivityLevel = req.ActivityLevel
	user.DietType = req.DietType
	user.Goal = req.Goal

	if err := r.db.Save(&user).Error; err != nil {
		return dtos.UserResponse{}, err
	}

	return dtos.UserResponse{
		ID:                user.ID,
		Name:              user.Name,
		Gender:            user.Gender,
		Height:            user.Height,
		HeightUnit:        user.HeightUnit,
		CurrentWeight:     user.CurrentWeight,
		CurrentWeightUnit: user.CurrentWeightUnit,
		TargetWeight:      user.TargetWeight,
		TargetWeightUnit:  user.TargetWeightUnit,
		ActivityLevel:     user.ActivityLevel,
		DietType:          user.DietType,
		Goal:              user.Goal,
	}, nil
}
