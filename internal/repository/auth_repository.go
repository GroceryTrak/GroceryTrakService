package repository

import (
	"errors"

	"github.com/GroceryTrak/GroceryTrakService/config"
	"github.com/GroceryTrak/GroceryTrakService/internal/dtos"
	"github.com/GroceryTrak/GroceryTrakService/internal/models"
	"github.com/GroceryTrak/GroceryTrakService/internal/utils"
)

func RegisterUser(req dtos.RegisterRequest, role models.Role) (dtos.RegisterResponse, error) {
	var existingUser models.User
	result := config.DB.Where("username = ?", req.Username).First(&existingUser).Error
	if result == nil {
		return dtos.RegisterResponse{}, errors.New("user already exists")
	}

	user := models.User{
		Username: req.Username,
		Password: utils.HashPassword(req.Password),
		Role:     "user",
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return dtos.RegisterResponse{}, err
	}

	return dtos.RegisterResponse{Message: "User registered successfully"}, nil
}

func LoginUser(req dtos.LoginRequest) (dtos.LoginResponse, error) {
	var user models.User
	result := config.DB.Where("username = ?", req.Username).First(&user)
	if result.Error != nil {
		return dtos.LoginResponse{}, errors.New("user not found")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return dtos.LoginResponse{}, errors.New("invalid credentials")
	}

	token, err := utils.GenerateJWT(user.ID, user.Username, user.Role)
	if err != nil {
		return dtos.LoginResponse{}, errors.New("could not generate token")
	}

	return dtos.LoginResponse{Token: token}, nil
}
