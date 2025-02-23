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

// PredictItemsHandler handles image upload and item prediction
func PredictItemsHandler(w http.ResponseWriter, r *http.Request) {
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

	var result []models.Item

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
			result = append(result, existingItem)
		} else {
			// Create new item if not found
			newItem := models.Item{
				Name:        class,
				Description: "",
				ImageLink:   "",
			}
			config.DB.Create(&newItem)
			result = append(result, newItem)
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
