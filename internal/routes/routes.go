package routes

import (
	"github.com/GroceryTrak/GroceryTrakService/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(r *chi.Mux) {
	r.Post("/register", handlers.RegisterHandler)
	r.Post("/login", handlers.LoginHandler)

	r.Route("/recipe", func(r chi.Router) {
		r.Get("/{id}", handlers.GetRecipeHandler)
		r.Post("/", handlers.CreateRecipeHandler)
		r.Put("/{id}", handlers.UpdateRecipeHandler)
		r.Delete("/{id}", handlers.DeleteRecipeHandler)
	})
}
