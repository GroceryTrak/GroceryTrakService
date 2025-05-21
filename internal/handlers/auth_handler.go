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
		json.NewEncoder(w).Encode(dtos.UnauthorizedResponse{Error: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(resp)
}

// @Summary Confirms a user's email
// @Description Confirms a user's email using the OTP code sent to their email during sign-up
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dtos.ConfirmRequest true "OTP Confirmation Request"
// @Success 200 {object} dtos.ConfirmResponse
// @Failure default {object} dtos.ErrorResponse "Standard Error Responses"
// @Router /auth/confirm [post]
func (h *AuthHandler) ConfirmHandler(w http.ResponseWriter, r *http.Request) {
	var req dtos.ConfirmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dtos.BadRequestResponse{Error: "Invalid request format"})
		return
	}

	resp, err := h.Repo.ConfirmUser(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.ErrorResponse{Error: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(resp)
}

// @Summary Resends the confirmation code to the user's email
// @Description Resends the email confirmation code in case the user did not receive it
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dtos.ResendRequest true "Resend Code Request"
// @Success 200 {object} dtos.ResendResponse
// @Failure default {object} dtos.ErrorResponse "Standard Error Responses"
// @Router /auth/resend [post]
func (h *AuthHandler) ResendHandler(w http.ResponseWriter, r *http.Request) {
	var req dtos.ResendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dtos.BadRequestResponse{Error: "Invalid request format"})
		return
	}

	resp, err := h.Repo.ResendCode(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.ErrorResponse{Error: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(resp)
}
