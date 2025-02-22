package main

import (
	"log"
	"net/http"
	"os"

	"github.com/GroceryTrak/GroceryTrakService/config"
	"github.com/GroceryTrak/GroceryTrakService/internal/routes"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func main() {
	config.LoadConfig()
	config.InitRedis()
	config.InitPostgreSQL()

	r := chi.NewRouter()

	// Enable CORS
	env := os.Getenv("ENV")
	frontendDomain := os.Getenv("FRONTEND_DOMAIN")

	allowedOrigins := []string{"*"}
	if env == "production" && frontendDomain != "" {
		allowedOrigins = []string{frontendDomain}
	}

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	routes.SetupRoutes(r)

	log.Println("Server running on port 8080")
	http.ListenAndServe(":8080", r)
}
