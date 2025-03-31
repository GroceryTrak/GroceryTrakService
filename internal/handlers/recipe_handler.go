package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/GroceryTrak/GroceryTrakService/internal/dtos"
	"github.com/GroceryTrak/GroceryTrakService/internal/repository"
	"github.com/go-chi/chi/v5"
)

type RecipeHandler struct {
	Repo repository.RecipeRepository
}

func NewRecipeHandler(repo repository.RecipeRepository) *RecipeHandler {
	return &RecipeHandler{Repo: repo}
}

// @Summary Get a recipe
// @Description Retrieves a recipe by its ID
// @Tags recipe
// @Accept json
// @Produce json
// @Param id path int true "Recipe ID"
// @Success 200 {object} dtos.RecipeResponse
// @Failure default {object} dtos.ErrorResponse "Standard Error Responses"
// @Router /recipe/{id} [get]
func (h *RecipeHandler) GetRecipeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dtos.BadRequestResponse{Error: "Invalid recipe ID"})
		return
	}

	recipe, err := h.Repo.GetRecipe(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(dtos.NotFoundResponse{Error: "Recipe not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipe)
}

// @Summary Create a recipe
// @Description Creates a new recipe
// @Tags recipe
// @Accept json
// @Produce json
// @Param recipe body dtos.RecipeRequest true "Recipe Data"
// @Success 201 {object} dtos.RecipeResponse
// @Failure default {object} dtos.ErrorResponse "Standard Error Responses"
// @Router /recipe [post]
func (h *RecipeHandler) CreateRecipeHandler(w http.ResponseWriter, r *http.Request) {
	var req dtos.RecipeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dtos.BadRequestResponse{Error: "Invalid request data"})
		return
	}

	recipe, err := h.Repo.CreateRecipe(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.InternalServerErrorResponse{Error: "Failed to create recipe"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(recipe)
}

// @Summary Update a recipe
// @Description Updates an existing recipe by ID
// @Tags recipe
// @Accept json
// @Produce json
// @Param id path int true "Recipe ID"
// @Param recipe body dtos.RecipeRequest true "Updated Recipe Data"
// @Success 200 {object} dtos.RecipeResponse
// @Failure default {object} dtos.ErrorResponse "Standard Error Responses"
// @Router /recipe/{id} [put]
func (h *RecipeHandler) UpdateRecipeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dtos.BadRequestResponse{Error: "Invalid recipe ID"})
		return
	}

	var req dtos.RecipeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dtos.BadRequestResponse{Error: "Invalid request data"})
		return
	}

	recipe, err := h.Repo.UpdateRecipe(uint(id), req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.InternalServerErrorResponse{Error: "Failed to update recipe"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipe)
}

// @Summary Delete a recipe
// @Description Deletes a recipe by ID
// @Tags recipe
// @Accept json
// @Produce json
// @Param id path int true "Recipe ID"
// @Success 204
// @Failure default {object} dtos.ErrorResponse "Standard Error Responses"
// @Router /recipe/{id} [delete]
func (h *RecipeHandler) DeleteRecipeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dtos.BadRequestResponse{Error: "Invalid recipe ID"})
		return
	}

	if err := h.Repo.DeleteRecipe(uint(id)); err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(dtos.NotFoundResponse{Error: "Recipe not found"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Search recipes
// @Description Searches for recipes by title, ingredients, or diet type
// @Tags recipe
// @Accept json
// @Produce json
// @Param title query string false "Title of recipe"
// @Param ingredients query string false "Comma-separated ingredient IDs"
// @Param diet query string false "Diet type"
// @Success 200 {object} dtos.RecipesResponse
// @Failure default {object} dtos.ErrorResponse "Standard Error Responses"
// @Router /recipe/search [get]
func (h *RecipeHandler) SearchRecipesHandler(w http.ResponseWriter, r *http.Request) {
	query := dtos.RecipeQuery{}
	query.Title = r.URL.Query().Get("title")
	query.Diet = r.URL.Query().Get("diet")
	query.Ingredients = strings.Split(r.URL.Query().Get("ingredients"), ",")

	recipes, err := h.Repo.SearchRecipes(query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.InternalServerErrorResponse{Error: "Database error"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipes)
}
