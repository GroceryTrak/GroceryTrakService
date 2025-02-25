package main

import (
	"log"
	"net/http"

	"github.com/GroceryTrak/GroceryTrakService/config"
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

// @response 400 {object} templates.BadRequestResponse "Bad Request"
// @response 401 {object} templates.UnauthorizedResponse "Unauthorized"
// @response 403 {object} templates.ForbiddenResponse "Forbidden"
// @response 404 {object} templates.NotFoundResponse "Not Found"
// @response 409 {object} templates.ConflictResponse "Conflict"
// @response 500 {object} templates.InternalServerErrorResponse "Internal Server Error"
func main() {
	config.LoadConfig()
	config.InitRedis()
	config.InitPostgreSQL()

	r := chi.NewRouter()
	routes.SetupRoutes(r)

	log.Println("Server running on port 8080")
	http.ListenAndServe(":8080", r)
}
