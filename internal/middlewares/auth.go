package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/GroceryTrak/GroceryTrakService/internal/utils"
)

// Context keys for storing user information
type contextKey string

const (
	IDKey   contextKey = "id"
	UserKey contextKey = "username"
	RoleKey contextKey = "role"
)

// AuthMiddleware ensures the request is authenticated
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization format", http.StatusUnauthorized)
			return
		}

		tokenString := tokenParts[1]

		// Extract user info
		userID, username, role, err := utils.VerifyToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Store in request context
		ctx := context.WithValue(r.Context(), IDKey, userID)
		ctx = context.WithValue(ctx, UserKey, username)
		ctx = context.WithValue(ctx, RoleKey, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserIDFromContext(r *http.Request) uint {
	if userID, ok := r.Context().Value(IDKey).(uint); ok {
		return userID
	}
	return 0
}

// GetUsernameFromContext extracts the username from request context
func GetUsernameFromContext(r *http.Request) string {
	if username, ok := r.Context().Value(UserKey).(string); ok {
		return username
	}
	return ""
}

// GetRoleFromContext extracts the role from request context
func GetRoleFromContext(r *http.Request) string {
	if role, ok := r.Context().Value(RoleKey).(string); ok {
		return role
	}
	return ""
}
