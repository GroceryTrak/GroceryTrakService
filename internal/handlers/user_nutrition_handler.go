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

// @Summary Calculate and save user nutrition goals
// @Description Calculate and save nutrition goals based on user's profile
// @Tags user_nutrition
// @Produce json
// @Success 200 {object} dtos.UserNutritionResponse
// @Failure default {object} dtos.ErrorResponse "Standard Error Responses"
// @Router /user/nutrition [post]
func (h *UserNutritionHandler) CalculateNutritionHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.IDKey).(string)

	nutrition, err := h.Repo.CalculateAndSaveNutrition(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.InternalServerErrorResponse{Error: "Failed to calculate nutrition goals"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nutrition)
}

// @Summary Get user nutrition goals
// @Description Get the authenticated user's nutrition goals
// @Tags user_nutrition
// @Produce json
// @Success 200 {object} dtos.UserNutritionResponse
// @Failure default {object} dtos.ErrorResponse "Standard Error Responses"
// @Router /user/nutrition [get]
func (h *UserNutritionHandler) GetNutritionHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.IDKey).(string)

	nutrition, err := h.Repo.GetUserNutrition(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.InternalServerErrorResponse{Error: "Failed to get nutrition goals"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nutrition)
}
