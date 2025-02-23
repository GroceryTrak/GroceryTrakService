package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/GroceryTrak/GroceryTrakService/config"
	"github.com/GroceryTrak/GroceryTrakService/internal/models"
	"github.com/go-chi/chi/v5"
)

// Get an item by ID
func GetItemHandler(w http.ResponseWriter, r *http.Request) {
	uintID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "Recipe not found", http.StatusNotFound)
		return
	}
	id := uint(uintID)

	var item models.Item
	if err := config.DB.First(&item, "id = ?", id).Error; err != nil {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// Create a new item
func CreateItemHandler(w http.ResponseWriter, r *http.Request) {
	var newItem models.Item
	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := config.DB.Create(&newItem).Error; err != nil {
		http.Error(w, "Failed to create item", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newItem)
}

// Update an item by ID
func UpdateItemHandler(w http.ResponseWriter, r *http.Request) {
	uintID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "Recipe not found", http.StatusNotFound)
		return
	}
	id := uint(uintID)

	var existingItem models.Item
	if err := config.DB.First(&existingItem, "id = ?", id).Error; err != nil {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	var updatedItem models.Item
	if err := json.NewDecoder(r.Body).Decode(&updatedItem); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedItem.ID = id
	if err := config.DB.Save(&updatedItem).Error; err != nil {
		http.Error(w, "Failed to update item", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedItem)
}

// Delete an item by ID
func DeleteItemHandler(w http.ResponseWriter, r *http.Request) {
	uintID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "Recipe not found", http.StatusNotFound)
		return
	}
	id := uint(uintID)

	if err := config.DB.Delete(&models.Item{}, "id = ?", id).Error; err != nil {
		http.Error(w, "Failed to delete item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Search items by substring
func SearchItemsHandler(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("q")
	if keyword == "" {
		http.Error(w, "Search keyword is required", http.StatusBadRequest)
		return
	}

	var items []models.Item
	searchTerm := "%" + keyword + "%"

	result := config.DB.Where("name LIKE ? OR description LIKE ?", searchTerm, searchTerm).Find(&items)
	if result.Error != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(items)
}
