package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/GroceryTrak/GroceryTrakService/config"
	"github.com/GroceryTrak/GroceryTrakService/internal/clients"
	"github.com/GroceryTrak/GroceryTrakService/internal/handlers"
	"github.com/GroceryTrak/GroceryTrakService/internal/repository"
	"github.com/GroceryTrak/GroceryTrakService/internal/routes"
	"github.com/go-chi/chi/v5"
)

// @title GroceryTrak API
// @version 1.0
// @description This is the API documentation for GroceryTrak.
// @termsOfService http://swagger.io/terms/
// @contact.name GroceryTrak
// @contact.email grocerytrak@gmail.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @response 400 {object} dtos.BadRequestResponse "Bad Request"
// @response 401 {object} dtos.UnauthorizedResponse "Unauthorized"
// @response 403 {object} dtos.ForbiddenResponse "Forbidden"
// @response 404 {object} dtos.NotFoundResponse "Not Found"
// @response 409 {object} dtos.ConflictResponse "Conflict"
// @response 500 {object} dtos.InternalServerErrorResponse "Internal Server Error"
func main() {
	config.LoadConfig()
	config.InitRedis()
	config.InitPostgreSQL()
	config.InitSpoonacularClient()
	clients.InitAWSCognito()

	// Initialize queue repository
	routes.InitQueue(config.RedisClient)

	// Create repositories
	itemRepo := repository.NewItemRepository(config.DB)
	itemQueueRepo := repository.NewItemQueueRepository(config.RedisClient)

	// Create and start ItemQueueHandler
	queueHandler := handlers.NewItemQueueHandler(
		itemQueueRepo,
		config.SpoonacularClient,
		itemRepo,
	)

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start queue handler in a goroutine
	go func() {
		if err := queueHandler.Start(ctx); err != nil {
			log.Printf("Queue handler error: %v", err)
		}
	}()

	// Setup HTTP server
	r := chi.NewRouter()
	routes.SetupRoutes(r)

	// Create server
	server := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: r,
	}

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Println("Server running on port 8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-quit
	log.Println("Shutting down server...")

	// Shutdown server gracefully
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// Cancel context to stop queue handler
	cancel()
	log.Println("Server exited properly")
}
