package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/GroceryTrak/GroceryTrakService/config"
	"github.com/GroceryTrak/GroceryTrakService/internal/middlewares"
	"github.com/GroceryTrak/GroceryTrakService/internal/models"
	"github.com/go-chi/chi/v5"
)

// Get all user's items
func GetAllUserItemHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.IDKey).(uint) // Get authenticated user ID

	var userItems []models.UserItem
	if err := config.DB.Where("user_id = ?", userID).Preload("Item").Find(&userItems).Error; err != nil {
		http.Error(w, "Failed to fetch user items", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userItems)
}

// Get a user's item by ItemID
func GetUserItemHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.IDKey).(uint) // Get authenticated user ID

	itemID, err := strconv.ParseUint(chi.URLParam(r, "item_id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	var userItem models.UserItem
	if err := config.DB.Where("user_id = ? AND item_id = ?", userID, itemID).First(&userItem).Error; err != nil {
		http.Error(w, "UserItem not found", http.StatusNotFound)
		return
	}

	// Validate if the item exists
	var item models.Item
	if err := config.DB.Where("id = ?", itemID).First(&item).Error; err != nil {
		http.Error(w, "Item not found", http.StatusBadRequest)
		return
	}

	// Assign the fetched item details
	userItem.Item = item

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userItem)
}

// Create a new user_item for the authenticated user
func CreateUserItemHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.IDKey).(uint) // Get authenticated user ID

	var newUserItem models.UserItem
	if err := json.NewDecoder(r.Body).Decode(&newUserItem); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Ensure user can only create items for themselves
	newUserItem.UserID = userID

	// Validate if the item exists
	var item models.Item
	if err := config.DB.Where("id = ?", newUserItem.ItemID).First(&item).Error; err != nil {
		http.Error(w, "Item not found", http.StatusBadRequest)
		return
	}

	// Assign the fetched item details
	newUserItem.Item = item

	// Prevent duplicate entries
	existingUserItem := models.UserItem{}
	if err := config.DB.Where("user_id = ? AND item_id = ?", userID, newUserItem.ItemID).First(&existingUserItem).Error; err == nil {
		http.Error(w, "UserItem already exists", http.StatusConflict)
		return
	}

	if err := config.DB.Create(&newUserItem).Error; err != nil {
		http.Error(w, "Failed to create user_item", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUserItem)
}

// Update a user_item (only if owned by user)
func UpdateUserItemHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.IDKey).(uint)

	itemID, err := strconv.ParseUint(chi.URLParam(r, "item_id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	var existingUserItem models.UserItem
	if err := config.DB.Where("user_id = ? AND item_id = ?", userID, itemID).First(&existingUserItem).Error; err != nil {
		http.Error(w, "UserItem not found", http.StatusNotFound)
		return
	}

	var updatedUserItem models.UserItem
	if err := json.NewDecoder(r.Body).Decode(&updatedUserItem); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedUserItem.UserID = userID
	updatedUserItem.ItemID = uint(itemID)

	if err := config.DB.Save(&updatedUserItem).Error; err != nil {
		http.Error(w, "Failed to update user_item", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUserItem)
}

// Delete a user_item (only if owned by user)
func DeleteUserItemHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.IDKey).(uint)

	itemID, err := strconv.ParseUint(chi.URLParam(r, "item_id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	if err := config.DB.Delete(&models.UserItem{}, "user_id = ? AND item_id = ?", userID, itemID).Error; err != nil {
		http.Error(w, "Failed to delete user_item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
