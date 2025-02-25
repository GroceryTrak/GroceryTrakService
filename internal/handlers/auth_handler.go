package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/GroceryTrak/GroceryTrakService/internal/repository"
	"github.com/GroceryTrak/GroceryTrakService/internal/templates"
	"github.com/GroceryTrak/GroceryTrakService/internal/utils"
)

// @Summary Register a new user
// @Description Creates a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body templates.RegisterRequest true "User Registration Request"
// @Failure default {object} templates.ErrorResponse "Standard Error Responses"
// @Router /auth/register [post]
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req templates.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(templates.BadRequestResponse{Error: "Invalid request format"})
		return
	}

	err := repository.RegisterUser(req.Username, req.Password, "user") // Default role to "user"
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(templates.ConflictResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(templates.RegisterResponse{Message: "User registered successfully"})
}

// @Summary Logs in a user
// @Description Logs in a user and returns a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body templates.LoginRequest true "User Login Request"
// @Success 200 {object} templates.AuthResponse
// @Failure default {object} templates.ErrorResponse "Standard Error Responses"
// @Router /auth/login [post]
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req templates.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(templates.BadRequestResponse{Error: "Invalid request format"})
		return
	}

	id, username, role, err := repository.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(templates.UnauthorizedResponse{Error: "Invalid credentials"})
		return
	}

	token, err := utils.GenerateJWT(id, username, role)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(templates.InternalServerErrorResponse{Error: "Could not generate token"})
		return
	}

	json.NewEncoder(w).Encode(templates.AuthResponse{Token: token})
}
