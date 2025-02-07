package main

import (
	"log"
	"net/http"

	"github.com/GroceryTrak/GroceryTrakService/config"
	"github.com/GroceryTrak/GroceryTrakService/internal/routes"
	"github.com/go-chi/chi/v5"
)

func main() {
	config.LoadConfig()
	config.InitRedis()
	config.InitPostgreSQL()

	r := chi.NewRouter()
	routes.SetupRoutes(r)

	log.Println("Server running on port 8080")
	http.ListenAndServe(":8080", r)
}
