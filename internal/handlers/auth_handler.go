package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/GroceryTrak/GroceryTrakService/internal/dtos"
	"github.com/GroceryTrak/GroceryTrakService/internal/repository"
)

type AuthHandler struct {
	Repo repository.AuthRepository
}

func NewAuthHandler(repo repository.AuthRepository) *AuthHandler {
	return &AuthHandler{Repo: repo}
}

// @Summary Register a new user
// @Description Creates a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dtos.RegisterRequest true "User Registration Request"
// @Success 201 {object} dtos.RegisterResponse
// @Failure default {object} dtos.ErrorResponse "Standard Error Responses"
// @Router /auth/register [post]
func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req dtos.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dtos.BadRequestResponse{Error: "Invalid request format"})
		return
	}

	resp, err := h.Repo.RegisterUser(req, "user") // Default role to "user"
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(dtos.ConflictResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// @Summary Logs in a user
// @Description Logs in a user and returns a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dtos.LoginRequest true "User Login Request"
// @Success 200 {object} dtos.LoginResponse
// @Failure default {object} dtos.ErrorResponse "Standard Error Responses"
// @Router /auth/login [post]
func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req dtos.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dtos.BadRequestResponse{Error: "Invalid request format"})
		return
	}

	resp, err := h.Repo.LoginUser(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(dtos.UnauthorizedResponse{Error: "Invalid credentials"})
		return
	}

	json.NewEncoder(w).Encode(resp)
}
