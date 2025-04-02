package middlewares

import (
	"net/http"
	"os"
	"time"

	"github.com/GroceryTrak/GroceryTrakService/config"
	"github.com/redis/go-redis/v9"
)

const (
	maxRequestSize = 10 * 1024 * 1024 // 10MB
)

func RequestSizeLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
			if r.ContentLength > maxRequestSize {
				http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func ProductionURLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		env := os.Getenv("ENV")
		if env == "production" {
			origin := r.Header.Get("Origin")
			referer := r.Header.Get("Referer")
			flutterURL := os.Getenv("FLUTTER_URL")

			if origin == flutterURL || referer == flutterURL {
				next.ServeHTTP(w, r)
				return
			}

			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		env := os.Getenv("ENV")
		if env != "production" {
			next.ServeHTTP(w, r)
			return
		}

		clientIP := r.RemoteAddr
		if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
			clientIP = forwardedFor
		}

		key := "rate_limit:" + clientIP

		val, err := config.RedisClient.Get(r.Context(), key).Int()
		if err == redis.Nil {
			config.RedisClient.Set(r.Context(), key, 1, 500*time.Millisecond)
			next.ServeHTTP(w, r)
			return
		}

		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		if val >= 1 {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		config.RedisClient.Incr(r.Context(), key)
		next.ServeHTTP(w, r)
	})
}

func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		if os.Getenv("ENV") == "production" {
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline';")
		next.ServeHTTP(w, r)
	})
}
