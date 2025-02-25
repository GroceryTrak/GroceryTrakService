package routes

import (
	"os"

	_ "github.com/GroceryTrak/GroceryTrakService/docs"
	"github.com/GroceryTrak/GroceryTrakService/internal/handlers"
	"github.com/GroceryTrak/GroceryTrakService/internal/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupRoutes(r *chi.Mux) {
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

	r.Route("/swagger", func(r chi.Router) {
		r.Get("/*", httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"),
		))
	})

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
		r.Post("/predict", handlers.PredictUserItemsHandler)
	})
}
