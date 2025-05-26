package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/GroceryTrak/GroceryTrakService/internal/clients"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

type contextKey string

const (
	IDKey contextKey = "id"
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

		input := &cognitoidentityprovider.GetUserInput{
			AccessToken: &tokenString,
		}

		_, err := clients.CognitoClient.GetUser(context.TODO(), input)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		claims, err := clients.ParseToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), IDKey, claims.Sub)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserIDFromContext(r *http.Request) string {
	if userID, ok := r.Context().Value(IDKey).(string); ok {
		return userID
	}
	return ""
}
