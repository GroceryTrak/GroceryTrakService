package routes

import (
	"github.com/GroceryTrak/GroceryTrakService/internal/handlers"
	"github.com/GroceryTrak/GroceryTrakService/internal/middlewares"
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
		r.Get("/search", handlers.SearchItemsHandler)
		r.Post("/predict", handlers.PredictItemsHandler)
	})

	r.Route("/recipe", func(r chi.Router) {
		r.Get("/{id}", handlers.GetRecipeHandler)
		r.Post("/", handlers.CreateRecipeHandler)
		r.Put("/{id}", handlers.UpdateRecipeHandler)
		r.Delete("/{id}", handlers.DeleteRecipeHandler)
		r.Get("/search", handlers.SearchRecipesHandler)
	})

	r.Route("/user_item", func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware)

		r.Get("/", handlers.GetAllUserItemHandler)
		r.Get("/{item_id}", handlers.GetUserItemHandler)
		r.Post("/", handlers.CreateUserItemHandler)
		r.Put("/{item_id}", handlers.UpdateUserItemHandler)
		r.Delete("/{item_id}", handlers.DeleteUserItemHandler)
	})
}
