// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "GroceryTrak",
            "email": "grocerytrak@gmail.com"
        },
        "license": {
            "name": "MIT",
            "url": "https://opensource.org/licenses/MIT"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/login": {
            "post": {
                "description": "Logs in a user and returns a JWT token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Logs in a user",
                "parameters": [
                    {
                        "description": "User Login Request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dtos.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.LoginResponse"
                        }
                    },
                    "default": {
                        "description": "Standard Error Responses",
                        "schema": {
                            "$ref": "#/definitions/dtos.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/auth/register": {
            "post": {
                "description": "Creates a new user account",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Register a new user",
                "parameters": [
                    {
                        "description": "User Registration Request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dtos.RegisterRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/dtos.RegisterResponse"
                        }
                    },
                    "default": {
                        "description": "Standard Error Responses",
                        "schema": {
                            "$ref": "#/definitions/dtos.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/item": {
            "post": {
                "description": "Add a new item to the database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "item"
                ],
                "summary": "Create an item",
                "parameters": [
                    {
                        "description": "New Item",
                        "name": "item",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dtos.ItemRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/dtos.ItemResponse"
                        }
                    },
                    "default": {
                        "description": "Standard Error Responses",
                        "schema": {
                            "$ref": "#/definitions/dtos.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/item/search": {
            "get": {
                "description": "Searches for items that match the provided keyword in their name or description",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "item"
                ],
                "summary": "Search items",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Search keyword",
                        "name": "name",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.ItemsResponse"
                        }
                    },
                    "default": {
                        "description": "Standard Error Responses",
                        "schema": {
                            "$ref": "#/definitions/dtos.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/item/{id}": {
            "get": {
                "description": "Get an item by its ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "item"
                ],
                "summary": "Get an item",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Item ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.ItemResponse"
                        }
                    },
                    "default": {
                        "description": "Standard Error Responses",
                        "schema": {
                            "$ref": "#/definitions/dtos.ErrorResponse"
                        }
                    }
                }
            },
            "put": {
                "description": "Update an existing item by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "item"
                ],
                "summary": "Update an item",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Item ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Updated Item Data",
                        "name": "item",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dtos.ItemRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.ItemResponse"
                        }
                    },
                    "default": {
                        "description": "Standard Error Responses",
                        "schema": {
                            "$ref": "#/definitions/dtos.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "Remove an item by its ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "item"
                ],
                "summary": "Delete an item",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Item ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "default": {
                        "description": "Standard Error Responses",
                        "schema": {
                            "$ref": "#/definitions/dtos.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/recipe": {
            "post": {
                "description": "Creates a new recipe",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "recipe"
                ],
                "summary": "Create a recipe",
                "parameters": [
                    {
                        "description": "Recipe Data",
                        "name": "recipe",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dtos.RecipeRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/dtos.RecipeResponse"
                        }
                    },
                    "default": {
                        "description": "Standard Error Responses",
                        "schema": {
                            "$ref": "#/definitions/dtos.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/recipe/search": {
            "get": {
                "description": "Searches for recipes by title, ingredients, or diet type",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "recipe"
                ],
                "summary": "Search recipes",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Title of recipe",
                        "name": "title",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Comma-separated ingredient IDs",
                        "name": "ingredients",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Diet type",
                        "name": "diet",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.RecipesResponse"
                        }
                    },
                    "default": {
                        "description": "Standard Error Responses",
                        "schema": {
                            "$ref": "#/definitions/dtos.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/recipe/{id}": {
            "get": {
                "description": "Retrieves a recipe by its ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "recipe"
                ],
                "summary": "Get a recipe",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Recipe ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.RecipeResponse"
                        }
                    },
                    "default": {
                        "description": "Standard Error Responses",
                        "schema": {
                            "$ref": "#/definitions/dtos.ErrorResponse"
                        }
                    }
                }
            },
            "put": {
                "description": "Updates an existing recipe by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "recipe"
                ],
                "summary": "Update a recipe",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Recipe ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Updated Recipe Data",
                        "name": "recipe",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dtos.RecipeRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.RecipeResponse"
                        }
                    },
                    "default": {
                        "description": "Standard Error Responses",
                        "schema": {
                            "$ref": "#/definitions/dtos.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "Deletes a recipe by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "recipe"
                ],
                "summary": "Delete a recipe",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Recipe ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "default": {
                        "description": "Standard Error Responses",
                        "schema": {
                            "$ref": "#/definitions/dtos.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/user_item": {
            "get": {
                "description": "Get all items for the authenticated user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user_item"
                ],
                "summary": "Get all user's items",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.UserItemsResponse"
                        }
                    },
                    "default": {
                        "description": "Standard Error Responses",
                        "schema": {
                            "$ref": "#/definitions/dtos.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Create a new item for the authenticated user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user_item"
                ],
                "summary": "Create a new user_item for the authenticated user",
                "parameters": [
                    {
                        "description": "Create User Item",
                        "name": "userItem",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dtos.UserItemRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/dtos.UserItemResponse"
                        }
                    },
                    "default": {
                        "description": "Standard Error Responses",
                        "schema": {
                            "$ref": "#/definitions/dtos.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/user_item/predict": {
            "post": {
                "description": "Predict items from an uploaded image for the authenticated user",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user_item"
                ],
                "summary": "Predict items from an uploaded image",
                "parameters": [
                    {
                        "type": "file",
                        "description": "Image file",
                        "name": "image",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.UserItemsResponse"
                        }
                    },
                    "default": {
                        "description": "Standard Error Responses",
                        "schema": {
                            "$ref": "#/definitions/dtos.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/user_item/search": {
            "get": {
                "description": "Searches for user items by name",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user_item"
                ],
                "summary": "Search user items",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Name of user item",
                        "name": "name",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.UserItemsResponse"
                        }
                    },
                    "default": {
                        "description": "Standard Error Responses",
                        "schema": {
                            "$ref": "#/definitions/dtos.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/user_item/{item_id}": {
            "get": {
                "description": "Get a specific item for the authenticated user by item ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user_item"
                ],
                "summary": "Get a user's item by ItemID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Item ID",
                        "name": "item_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.UserItemResponse"
                        }
                    },
                    "default": {
                        "description": "Standard Error Responses",
                        "schema": {
                            "$ref": "#/definitions/dtos.ErrorResponse"
                        }
                    }
                }
            },
            "put": {
                "description": "Update a user_item for the authenticated user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user_item"
                ],
                "summary": "Update a user_item for the authenticated user",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Item ID",
                        "name": "item_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Update User Item",
                        "name": "userItem",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dtos.UserItemRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.UserItemResponse"
                        }
                    },
                    "default": {
                        "description": "Standard Error Responses",
                        "schema": {
                            "$ref": "#/definitions/dtos.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete a user_item for the authenticated user",
                "tags": [
                    "user_item"
                ],
                "summary": "Delete a user_item for the authenticated user",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Item ID",
                        "name": "item_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "default": {
                        "description": "Standard Error Responses",
                        "schema": {
                            "$ref": "#/definitions/dtos.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dtos.BadRequestResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "Invalid request data"
                }
            }
        },
        "dtos.ConflictResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "User already exists"
                }
            }
        },
        "dtos.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "An error occurred"
                }
            }
        },
        "dtos.ForbiddenResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "Access denied"
                }
            }
        },
        "dtos.InternalServerErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "Internal server error"
                }
            }
        },
        "dtos.ItemRequest": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "image": {
                    "type": "string",
                    "example": "milk.jpg"
                },
                "name": {
                    "type": "string",
                    "example": "Milk"
                },
                "spoonacular_id": {
                    "type": "integer",
                    "example": 1
                }
            }
        },
        "dtos.ItemResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "image": {
                    "type": "string",
                    "example": "milk.jpg"
                },
                "name": {
                    "type": "string",
                    "example": "Milk"
                },
                "spoonacular_id": {
                    "type": "integer",
                    "example": 1
                }
            }
        },
        "dtos.ItemsResponse": {
            "type": "object",
            "properties": {
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dtos.ItemResponse"
                    }
                }
            }
        },
        "dtos.LoginRequest": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string",
                    "example": "securepassword123"
                },
                "username": {
                    "type": "string",
                    "example": "john_doe"
                }
            }
        },
        "dtos.LoginResponse": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string",
                    "example": "eyJhbGciOiJIUzI1NiIsInR..."
                }
            }
        },
        "dtos.NotFoundResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "Resource not found"
                }
            }
        },
        "dtos.RecipeItemRequest": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number",
                    "example": 2.5
                },
                "item_id": {
                    "type": "integer",
                    "example": 456
                },
                "unit": {
                    "type": "string",
                    "example": "cups"
                }
            }
        },
        "dtos.RecipeItemResponse": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number",
                    "example": 2.5
                },
                "item": {
                    "$ref": "#/definitions/dtos.ItemResponse"
                },
                "unit": {
                    "type": "string",
                    "example": "cups"
                }
            }
        },
        "dtos.RecipeRequest": {
            "type": "object",
            "properties": {
                "cooking_time": {
                    "type": "integer",
                    "example": 20
                },
                "image": {
                    "type": "string",
                    "example": "https://example.com/spaghetti.jpg"
                },
                "ingredients": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dtos.RecipeItemRequest"
                    }
                },
                "kcal": {
                    "type": "number",
                    "example": 450.5
                },
                "prep_time": {
                    "type": "integer",
                    "example": 10
                },
                "ready_time": {
                    "type": "integer",
                    "example": 30
                },
                "title": {
                    "type": "string",
                    "example": "Spaghetti Carbonara"
                },
                "vegan": {
                    "type": "boolean",
                    "example": false
                },
                "vegetarian": {
                    "type": "boolean",
                    "example": false
                }
            }
        },
        "dtos.RecipeResponse": {
            "type": "object",
            "properties": {
                "cooking_time": {
                    "type": "integer",
                    "example": 20
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "image": {
                    "type": "string",
                    "example": "https://example.com/spaghetti.jpg"
                },
                "ingredients": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dtos.RecipeItemResponse"
                    }
                },
                "kcal": {
                    "type": "number",
                    "example": 450.5
                },
                "prep_time": {
                    "type": "integer",
                    "example": 10
                },
                "ready_time": {
                    "type": "integer",
                    "example": 30
                },
                "title": {
                    "type": "string",
                    "example": "Spaghetti Carbonara"
                },
                "vegan": {
                    "type": "boolean",
                    "example": false
                },
                "vegetarian": {
                    "type": "boolean",
                    "example": false
                }
            }
        },
        "dtos.RecipesResponse": {
            "type": "object",
            "properties": {
                "recipes": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dtos.RecipeResponse"
                    }
                }
            }
        },
        "dtos.RegisterRequest": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string",
                    "example": "securepassword123"
                },
                "username": {
                    "type": "string",
                    "example": "john_doe"
                }
            }
        },
        "dtos.RegisterResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "User registered successfully"
                }
            }
        },
        "dtos.UnauthorizedResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "Invalid credentials"
                }
            }
        },
        "dtos.UserItemRequest": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number",
                    "example": 2
                },
                "item_id": {
                    "type": "integer",
                    "example": 456
                },
                "unit": {
                    "type": "string",
                    "example": "kg"
                }
            }
        },
        "dtos.UserItemResponse": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number",
                    "example": 2
                },
                "item": {
                    "$ref": "#/definitions/dtos.ItemResponse"
                },
                "unit": {
                    "type": "string",
                    "example": "kg"
                }
            }
        },
        "dtos.UserItemsResponse": {
            "type": "object",
            "properties": {
                "user_items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dtos.UserItemResponse"
                    }
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "GroceryTrak API",
	Description:      "This is the API documentation for GroceryTrak.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
