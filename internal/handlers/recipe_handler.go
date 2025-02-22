package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/GroceryTrak/GroceryTrakService/config"
	"github.com/GroceryTrak/GroceryTrakService/internal/models"
	"github.com/go-chi/chi/v5"
)

// Get a recipe by ID
func GetRecipeHandler(w http.ResponseWriter, r *http.Request) {
	uintID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}
	id := uint(uintID)

	var recipe models.Recipe
	if err := config.DB.Preload("Ingredients.Item").First(&recipe, "id = ?", id).Error; err != nil {
		http.Error(w, "Recipe not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipe)
}

func CreateRecipeHandler(w http.ResponseWriter, r *http.Request) {
	var newRecipe models.Recipe
	if err := json.NewDecoder(r.Body).Decode(&newRecipe); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Create Recipe along with Ingredients
	if err := config.DB.Create(&newRecipe).Error; err != nil {
		http.Error(w, "Failed to create recipe", http.StatusInternalServerError)
		return
	}

	// Manually fetch related `Item` data for each ingredient
	for i := range newRecipe.Ingredients {
		var item models.Item
		if err := config.DB.First(&item, "id = ?", newRecipe.Ingredients[i].ItemID).Error; err == nil {
			newRecipe.Ingredients[i].Item = item
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newRecipe)
}

// Update a recipe by ID (with conditional ingredient updates)
func UpdateRecipeHandler(w http.ResponseWriter, r *http.Request) {
	uintID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}
	id := uint(uintID)

	var existingRecipe models.Recipe
	if err := config.DB.Preload("Ingredients.Item").First(&existingRecipe, "id = ?", id).Error; err != nil {
		http.Error(w, "Recipe not found", http.StatusNotFound)
		return
	}

	var updatePayload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updatePayload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Update basic fields if present
	if name, ok := updatePayload["name"].(string); ok {
		existingRecipe.Name = name
	}
	if instruction, ok := updatePayload["instruction"].(string); ok {
		existingRecipe.Instruction = instruction
	}
	if difficultyStr, ok := updatePayload["difficulty"].(string); ok {
		if difficulty, err := models.ParseDifficulty(difficultyStr); err == nil {
			existingRecipe.Difficulty = difficulty
		} else {
			http.Error(w, "Invalid difficulty value", http.StatusBadRequest)
			return
		}
	}
	if duration, ok := updatePayload["duration"].(float64); ok {
		existingRecipe.Duration = uint(duration)
	}
	if kcal, ok := updatePayload["kcal"].(float64); ok {
		existingRecipe.KCal = uint(kcal)
	}

	// Handle ingredients only if the field exists
	if ingredientsRaw, exists := updatePayload["ingredients"]; exists {
		config.DB.Where("recipe_id = ?", id).Delete(&models.RecipeItem{}) // Clear old items

		if ingredients, ok := ingredientsRaw.([]interface{}); ok {
			var newIngredients []models.RecipeItem
			for _, rawIng := range ingredients {
				ingMap, _ := rawIng.(map[string]interface{})

				itemID := uint(ingMap["item_id"].(float64))
				quantity := uint(ingMap["quantity"].(float64))
				unit, _ := ingMap["unit"].(string)

				var item models.Item
				if err := config.DB.First(&item, "id = ?", itemID).Error; err != nil {
					item = models.Item{} // Default empty item
				}

				newIngredients = append(newIngredients, models.RecipeItem{
					RecipeID: id,
					ItemID:   itemID,
					Quantity: quantity,
					Unit:     unit,
					Item:     item, // ✅ Assigning the item object
				})
			}
			existingRecipe.Ingredients = newIngredients
		}
	}

	if err := config.DB.Save(&existingRecipe).Error; err != nil {
		http.Error(w, "Failed to update recipe", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingRecipe)
}

// Delete a recipe by ID
func DeleteRecipeHandler(w http.ResponseWriter, r *http.Request) {
	uintID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}
	id := uint(uintID)

	// Delete Recipe (RecipeItems will be automatically deleted due to CASCADE)
	if err := config.DB.Delete(&models.Recipe{}, "id = ?", id).Error; err != nil {
		http.Error(w, "Failed to delete recipe", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Search recipes by substring
func SearchRecipesHandler(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("q")
	ingredientParam := r.URL.Query().Get("ingredients") // Comma-separated ingredient IDs
	diet := r.URL.Query().Get("diet")

	var recipes []models.Recipe
	query := config.DB.Preload("Ingredients.Item")

	// Filter by name if a keyword is provided
	if keyword != "" {
		searchTerm := "%" + keyword + "%"
		query = query.Where("name LIKE ?", searchTerm)
	}

	// Filter by ingredients if the parameter exists
	if ingredientParam != "" {
		ingredientIDs := strings.Split(ingredientParam, ",") // Convert string to slice
		query = query.Joins("JOIN recipe_items ri ON ri.recipe_id = recipes.id").
			Where("ri.item_id IN (?)", ingredientIDs).
			Group("recipes.id").
			Having("COUNT(DISTINCT ri.item_id) = ?", len(ingredientIDs)) // Ensure all selected items exist
	}

	// Filter by diet if provided
	if diet != "" {
		dietSearchTerm := "%" + diet + "%"
		query = query.Where("diet LIKE ?", dietSearchTerm)
	}

	// Execute query
	result := query.Find(&recipes)
	if result.Error != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(recipes)
}
