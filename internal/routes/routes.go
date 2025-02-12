package routes

import (
	"github.com/GroceryTrak/GroceryTrakService/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(r *chi.Mux) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", handlers.RegisterHandler)
		r.Post("/login", handlers.LoginHandler)
	})

	r.Route("/item", func(r chi.Router) {
		r.Get("/{id}", handlers.GetItemHandler)
		r.Post("/", handlers.CreateItemHandler)
		r.Put("/{id}", handlers.UpdateItemHandler)
		r.Delete("/{id}", handlers.DeleteItemHandler)
	})

	r.Route("/recipe", func(r chi.Router) {
		r.Get("/{id}", handlers.GetRecipeHandler)
		r.Post("/", handlers.CreateRecipeHandler)
		r.Put("/{id}", handlers.UpdateRecipeHandler)
		r.Delete("/{id}", handlers.DeleteRecipeHandler)
	})
}
