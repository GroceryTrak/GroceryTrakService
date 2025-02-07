package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/GroceryTrak/GroceryTrakService/internal/handlers"
)

func SetupRoutes(r *chi.Mux) {
	r.Post("/register", handlers.RegisterHandler)
	r.Post("/login", handlers.LoginHandler)
}