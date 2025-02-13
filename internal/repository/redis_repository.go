package repository

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/GroceryTrak/GroceryTrakService/config"
	"github.com/GroceryTrak/GroceryTrakService/internal/utils"
)

func RegisterUser(username, password, role string) error {
	userKey := fmt.Sprintf("user:%s", username)

	exists, err := config.RedisClient.Exists(config.Ctx, userKey).Result()
	if err != nil {
		return err
	}
	if exists > 0 {
		return errors.New("user already exists")
	}

	hashedPassword := utils.HashPassword(password)

	err = config.RedisClient.HSet(config.Ctx, userKey, map[string]interface{}{
		"password": hashedPassword,
		"role":     role,
	}).Err()

	if err != nil {
		return err
	}

	config.RedisClient.Expire(config.Ctx, userKey, 24*time.Hour)
	return nil
}

func AuthenticateUser(username, password string) (uint, string, string, error) {
	userKey := fmt.Sprintf("user:%s", username)

	exists, err := config.RedisClient.Exists(config.Ctx, userKey).Result()
	if err != nil || exists == 0 {
		return 0, "", "", errors.New("user not found")
	}

	storedPassword, err := config.RedisClient.HGet(config.Ctx, userKey, "password").Result()
	if err != nil {
		return 0, "", "", err
	}

	if !utils.CheckPassword(password, storedPassword) {
		return 0, "", "", errors.New("invalid credentials")
	}

	idStr, err := config.RedisClient.HGet(config.Ctx, userKey, "id").Result()
	if err != nil {
		return 0, "", "", err
	}

	role, _ := config.RedisClient.HGet(config.Ctx, userKey, "role").Result()

	// Convert ID from string to uint
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, "", "", errors.New("invalid user ID format")
	}

	return uint(id), username, role, nil
}
