package utils

import (
	"os"
	"time"

	"github.com/GroceryTrak/GroceryTrakService/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateJWT(id uint, username string, role models.Role) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       id,
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	return token.SignedString(jwtSecret)
}

// VerifyToken parses and validates a JWT, returning the id, username and role
func VerifyToken(tokenString string) (uint, string, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return 0, "", "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id := uint(claims["id"].(float64))
		username := claims["username"].(string)
		role := claims["role"].(string)
		return id, username, role, nil
	}

	return 0, "", "", err
}
