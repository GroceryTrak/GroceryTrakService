package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/GroceryTrak/GroceryTrakService/config"
	"github.com/GroceryTrak/GroceryTrakService/internal/models"
	"github.com/go-chi/chi/v5"
)

// Get a recipe by ID
func GetRecipeHandler(w http.ResponseWriter, r *http.Request) {
	uintID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "Recipe not found", http.StatusNotFound)
		return
	}
	id := uint(uintID)

	var recipe models.Recipe
	if err := config.DB.Preload("Ingredients").First(&recipe, "id = ?", id).Error; err != nil {
		http.Error(w, "Recipe not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipe)
}

// Create a new recipe
func CreateRecipeHandler(w http.ResponseWriter, r *http.Request) {
	var newRecipe models.Recipe
	if err := json.NewDecoder(r.Body).Decode(&newRecipe); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := config.DB.Create(&newRecipe).Error; err != nil {
		http.Error(w, "Failed to create recipe", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newRecipe)
}

// Update a recipe by ID
func UpdateRecipeHandler(w http.ResponseWriter, r *http.Request) {
	uintID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "Recipe not found", http.StatusNotFound)
		return
	}
	id := uint(uintID)

	var existingRecipe models.Recipe
	if err := config.DB.First(&existingRecipe, "id = ?", id).Error; err != nil {
		http.Error(w, "Recipe not found", http.StatusNotFound)
		return
	}

	var updatedRecipe models.Recipe
	if err := json.NewDecoder(r.Body).Decode(&updatedRecipe); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedRecipe.ID = id
	if err := config.DB.Save(&updatedRecipe).Error; err != nil {
		http.Error(w, "Failed to update recipe", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedRecipe)
}

// Delete a recipe by ID
func DeleteRecipeHandler(w http.ResponseWriter, r *http.Request) {
	uintID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "Recipe not found", http.StatusNotFound)
		return
	}
	id := uint(uintID)

	if err := config.DB.Delete(&models.Recipe{}, "id = ?", id).Error; err != nil {
		http.Error(w, "Failed to delete recipe", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
