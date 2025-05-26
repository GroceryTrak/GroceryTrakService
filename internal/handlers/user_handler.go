package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/GroceryTrak/GroceryTrakService/internal/dtos"
	"github.com/GroceryTrak/GroceryTrakService/internal/middlewares"
	"github.com/GroceryTrak/GroceryTrakService/internal/repository"
)

type UserHandler struct {
	Repo repository.UserRepository
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{Repo: repo}
}

// @Summary Get user profile
// @Description Get the authenticated user's profile
// @Tags user
// @Produce json
// @Success 200 {object} dtos.UserResponse
// @Failure default {object} dtos.ErrorResponse "Standard Error Responses"
// @Router /user [get]
func (h *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.IDKey).(string)

	user, err := h.Repo.GetUser(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.InternalServerErrorResponse{Error: "Failed to get user profile"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// @Summary Update user profile
// @Description Update the authenticated user's profile
// @Tags user
// @Accept json
// @Produce json
// @Param user body dtos.UpdateUserRequest true "Update User Profile"
// @Success 200 {object} dtos.UserResponse
// @Failure default {object} dtos.ErrorResponse "Standard Error Responses"
// @Router /user [put]
func (h *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.IDKey).(string)

	var req dtos.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dtos.BadRequestResponse{Error: "Invalid request payload"})
		return
	}

	user, err := h.Repo.UpdateUser(req, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.InternalServerErrorResponse{Error: "Failed to update user profile"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
