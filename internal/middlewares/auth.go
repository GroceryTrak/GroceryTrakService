package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/GroceryTrak/GroceryTrakService/internal/utils"
)

type contextKey string

const (
	IDKey   contextKey = "id"
	UserKey contextKey = "username"
	RoleKey contextKey = "role"
)

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

		userID, username, role, err := utils.VerifyToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

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

func GetUsernameFromContext(r *http.Request) string {
	if username, ok := r.Context().Value(UserKey).(string); ok {
		return username
	}
	return ""
}

func GetRoleFromContext(r *http.Request) string {
	if role, ok := r.Context().Value(RoleKey).(string); ok {
		return role
	}
	return ""
}
