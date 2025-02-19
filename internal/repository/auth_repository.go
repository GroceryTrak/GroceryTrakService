package repository

import (
	"errors"

	"github.com/GroceryTrak/GroceryTrakService/config"
	"github.com/GroceryTrak/GroceryTrakService/internal/models"
	"github.com/GroceryTrak/GroceryTrakService/internal/utils"
)

func RegisterUser(username, password, role string) error {
	var existingUser models.User
	result := config.DB.Where("username = ?", username).First(&existingUser)
	if result.Error == nil {
		return errors.New("user already exists")
	}

	// Hash the password
	hashedPassword := utils.HashPassword(password)

	user := models.User{
		Username: username,
		Password: hashedPassword,
		Role:     role,
	}

	// Save user to PostgreSQL
	if err := config.DB.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func AuthenticateUser(username, password string) (uint, string, string, error) {
	var user models.User
	result := config.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return 0, "", "", errors.New("user not found")
	}

	// Verify password
	if !utils.CheckPassword(password, user.Password) {
		return 0, "", "", errors.New("invalid credentials")
	}

	return user.ID, user.Username, user.Role, nil
}
