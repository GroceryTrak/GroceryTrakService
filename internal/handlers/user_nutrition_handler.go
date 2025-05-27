package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/GroceryTrak/GroceryTrakService/internal/dtos"
	"github.com/GroceryTrak/GroceryTrakService/internal/middlewares"
	"github.com/GroceryTrak/GroceryTrakService/internal/repository"
)

type UserNutritionHandler struct {
	Repo repository.UserNutritionRepository
}

func NewUserNutritionHandler(repo repository.UserNutritionRepository) *UserNutritionHandler {
	return &UserNutritionHandler{
		Repo: repo,
	}
}

// @Summary Get user nutrition goals
// @Description Get the authenticated user's nutrition goals
// @Tags user-nutrition
// @Produce json
// @Success 200 {object} dtos.UserNutritionResponse
// @Failure default {object} dtos.ErrorResponse "Standard Error Responses"
// @Router /user/nutrition [get]
func (h *UserNutritionHandler) GetUserNutritionHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.IDKey).(string)

	nutrition, err := h.Repo.GetUserNutrition(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.InternalServerErrorResponse{Error: "Failed to get user nutrition"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nutrition)
}

// @Summary Create user nutrition
// @Description Create nutrition profile for the authenticated user
// @Tags user-nutrition
// @Produce json
// @Success 201 {object} dtos.UserNutritionResponse
// @Failure default {object} dtos.ErrorResponse "Standard Error Responses"
// @Router /user/nutrition [post]
func (h *UserNutritionHandler) CreateUserNutritionHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.IDKey).(string)

	nutrition, err := h.Repo.CreateUserNutrition(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.InternalServerErrorResponse{Error: "Failed to create user nutrition"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(nutrition)
}

// @Summary Update user nutrition
// @Description Update nutrition profile for the authenticated user
// @Tags user-nutrition
// @Produce json
// @Success 200 {object} dtos.UserNutritionResponse
// @Failure default {object} dtos.ErrorResponse "Standard Error Responses"
// @Router /user/nutrition [put]
func (h *UserNutritionHandler) UpdateUserNutritionHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.IDKey).(string)

	nutrition, err := h.Repo.UpdateUserNutrition(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.InternalServerErrorResponse{Error: "Failed to update user nutrition"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nutrition)
}

// @Summary Delete user nutrition
// @Description Delete nutrition profile for the authenticated user
// @Tags user-nutrition
// @Produce json
// @Success 204 "No Content"
// @Failure default {object} dtos.ErrorResponse "Standard Error Responses"
// @Router /user/nutrition [delete]
func (h *UserNutritionHandler) DeleteUserNutritionHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.IDKey).(string)

	if err := h.Repo.DeleteUserNutrition(userID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.InternalServerErrorResponse{Error: "Failed to delete user nutrition"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
