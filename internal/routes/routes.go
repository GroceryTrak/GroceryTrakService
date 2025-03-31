package routes

import (
	"os"

	"github.com/GroceryTrak/GroceryTrakService/config"
	_ "github.com/GroceryTrak/GroceryTrakService/docs"
	"github.com/GroceryTrak/GroceryTrakService/internal/handlers"
	"github.com/GroceryTrak/GroceryTrakService/internal/middlewares"
	"github.com/GroceryTrak/GroceryTrakService/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupDependencies() (*handlers.ItemHandler, *handlers.AuthHandler, *handlers.RecipeHandler, *handlers.UserItemHandler) {
	itemRepo := repository.NewItemRepository(config.DB)
	authRepo := repository.NewAuthRepository(config.DB)
	recipeRepo := repository.NewRecipeRepository(config.DB)
	userItemRepo := repository.NewUserItemRepository(config.DB)

	return handlers.NewItemHandler(itemRepo),
		handlers.NewAuthHandler(authRepo),
		handlers.NewRecipeHandler(recipeRepo),
		handlers.NewUserItemHandler(userItemRepo)
}

func SetupRoutes(r *chi.Mux) {
	itemHandler, authHandler, recipeHandler, userItemHandler := SetupDependencies()

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
		r.Post("/register", authHandler.RegisterHandler)
		r.Post("/login", authHandler.LoginHandler)
	})

	r.Route("/item", func(r chi.Router) {
		r.Get("/{id}", itemHandler.GetItemHandler)
		r.Post("/", itemHandler.CreateItemHandler)
		r.Put("/{id}", itemHandler.UpdateItemHandler)
		r.Delete("/{id}", itemHandler.DeleteItemHandler)
		r.Get("/search", itemHandler.SearchItemsHandler)
	})

	r.Route("/recipe", func(r chi.Router) {
		r.Get("/{id}", recipeHandler.GetRecipeHandler)
		r.Post("/", recipeHandler.CreateRecipeHandler)
		r.Put("/{id}", recipeHandler.UpdateRecipeHandler)
		r.Delete("/{id}", recipeHandler.DeleteRecipeHandler)
		r.Get("/search", recipeHandler.SearchRecipesHandler)
	})

	r.Route("/user_item", func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware)

		r.Get("/", userItemHandler.GetAllUserItemsHandler)
		r.Get("/{item_id}", userItemHandler.GetUserItemHandler)
		r.Post("/", userItemHandler.CreateUserItemHandler)
		r.Put("/{item_id}", userItemHandler.UpdateUserItemHandler)
		r.Delete("/{item_id}", userItemHandler.DeleteUserItemHandler)
		r.Post("/predict", userItemHandler.PredictUserItemsHandler)
	})
}
