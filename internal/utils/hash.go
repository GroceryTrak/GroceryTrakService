package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword encrypts a password using bcrypt
func HashPassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword)
}

// CheckPassword compares a hashed password with a plaintext one
func CheckPassword(plainPassword, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
