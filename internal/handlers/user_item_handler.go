package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"unicode"

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

// PredictItemsHandler handles image upload and item prediction
func PredictUserItemsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.IDKey).(uint) // Get authenticated user ID

	// Ensure method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10MB limit
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get file from request
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Image file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read file into memory
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read image", http.StatusInternalServerError)
		return
	}

	// Send image to DETECT_DOMAIN/predict
	detectDomain := os.Getenv("DETECT_DOMAIN")
	if detectDomain == "" {
		http.Error(w, "DETECT_DOMAIN is not set", http.StatusInternalServerError)
		return
	}

	// Create a buffer to store the multipart data
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Create the file field (key: "image", filename: "upload.png")
	filePart, err := writer.CreateFormFile("image", "upload.png")
	if err != nil {
		http.Error(w, "Failed to create form file", http.StatusInternalServerError)
		return
	}

	// Copy file data to the form field
	if _, err := filePart.Write(fileBytes); err != nil {
		http.Error(w, "Failed to write file data", http.StatusInternalServerError)
		return
	}

	// Close the writer to finalize the multipart form
	writer.Close()

	// Create request
	predictURL := fmt.Sprintf("%s/predict", detectDomain)
	req, err := http.NewRequest("POST", predictURL, &buf)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// Set Content-Type header with the boundary from writer
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to get prediction", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	log.Println(resp)

	// Decode response
	var response struct {
		Items          []string `json:"items"`
		AnnotatedImage string   `json:"annotated_image"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		http.Error(w, "Failed to parse prediction response", http.StatusInternalServerError)
		return
	}

	log.Println("Received classes:", response.Items)

	var result []models.UserItem

	for _, class := range response.Items {
		log.Println("Processing class:", class)
		if class == "" {
			continue
		}

		// Capitalize first letter
		class = capitalizeFirstLetter(class)

		// Search for existing item
		var existingItem models.Item
		if err := config.DB.Where("name LIKE ?", "%"+class+"%").First(&existingItem).Error; err == nil {
		} else {
			// Create new item if not found
			newItem := models.Item{
				Name:        class,
				Description: "",
				ImageLink:   "",
			}
			config.DB.Create(&newItem)
			existingItem = newItem
		}

		// Create UserItem entry
		userItem := models.UserItem{
			UserID: userID, // Convert to uint if necessary
			ItemID: existingItem.ID,
			Item:   existingItem,
		}
		result = append(result, userItem)

		// Check if UserItem already exists to prevent duplicates
		var existingUserItem models.UserItem
		if err := config.DB.Where("user_id = ? AND item_id = ?", userItem.UserID, userItem.ItemID).First(&existingUserItem).Error; err != nil {
			config.DB.Create(&userItem)
		}
	}

	// Return JSON result
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// capitalizeFirstLetter ensures the first letter is uppercase
func capitalizeFirstLetter(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}
