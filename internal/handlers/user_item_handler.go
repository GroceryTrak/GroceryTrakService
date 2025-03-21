package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/GroceryTrak/GroceryTrakService/config"
	"github.com/GroceryTrak/GroceryTrakService/internal/dtos"
	"github.com/GroceryTrak/GroceryTrakService/internal/middlewares"
	"github.com/GroceryTrak/GroceryTrakService/internal/models"
	"github.com/GroceryTrak/GroceryTrakService/internal/repository"
	"github.com/go-chi/chi/v5"
)

// @Summary Get all user's items
// @Description Get all items for the authenticated user
// @Tags user_item
// @Produce json
// @Success 200 {object} dtos.UserItemsResponse
// @Failure default {object} dtos.ErrorResponse "Standard Error Responses"
// @Router /user_item [get]
func GetAllUserItemsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.IDKey).(uint)

	userItems, err := repository.GetAllUserItems(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.ErrorResponse{Error: "Failed to fetch user items"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userItems)
}

// @Summary Get a user's item by ItemID
// @Description Get a specific item for the authenticated user by item ID
// @Tags user_item
// @Produce json
// @Param item_id path int true "Item ID"
// @Success 200 {object} dtos.UserItemResponse
// @Failure default {object} dtos.ErrorResponse "Standard Error Responses"
// @Router /user_item/{item_id} [get]
func GetUserItemHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.IDKey).(uint)

	itemID, err := strconv.ParseUint(chi.URLParam(r, "item_id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		json.NewEncoder(w).Encode(dtos.ErrorResponse{Error: "Invalid user item ID"})
		return
	}

	userItem, err := repository.GetUserItem(uint(itemID), userID)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(dtos.NotFoundResponse{Error: "User item not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userItem)
}

// @Summary Create a new user_item for the authenticated user
// @Description Create a new item for the authenticated user
// @Tags user_item
// @Accept json
// @Produce json
// @Param userItem body dtos.UserItemRequest true "Create User Item"
// @Success 201 {object} dtos.UserItemResponse
// @Failure default {object} dtos.ErrorResponse "Standard Error Responses"
// @Router /user_item [post]
func CreateUserItemHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.IDKey).(uint)

	var req dtos.UserItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		json.NewEncoder(w).Encode(dtos.ErrorResponse{Error: "Invalid request payload"})
		return
	}

	userItem, err := repository.CreateUserItem(req, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.InternalServerErrorResponse{Error: "Failed to create user item"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userItem)
}

// @Summary Update a user_item for the authenticated user
// @Description Update a user_item for the authenticated user
// @Tags user_item
// @Accept json
// @Produce json
// @Param item_id path int true "Item ID"
// @Param userItem body dtos.UserItemRequest true "Update User Item"
// @Success 200 {object} dtos.UserItemResponse
// @Failure default {object} dtos.ErrorResponse "Standard Error Responses"
// @Router /user_item/{item_id} [put]
func UpdateUserItemHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.IDKey).(uint)

	itemID, err := strconv.ParseUint(chi.URLParam(r, "item_id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		json.NewEncoder(w).Encode(dtos.ErrorResponse{Error: "Invalid item ID"})
		return
	}

	var req dtos.UserItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dtos.BadRequestResponse{Error: "Invalid request data"})
		return
	}

	userItem, err := repository.UpdateUserItem(req, uint(itemID), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.InternalServerErrorResponse{Error: "Failed to update user item"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userItem)
}

// @Summary Delete a user_item for the authenticated user
// @Description Delete a user_item for the authenticated user
// @Tags user_item
// @Param item_id path int true "Item ID"
// @Success 204
// @Failure default {object} dtos.ErrorResponse "Standard Error Responses"
// @Router /user_item/{item_id} [delete]
func DeleteUserItemHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.IDKey).(uint)

	itemID, err := strconv.ParseUint(chi.URLParam(r, "item_id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		json.NewEncoder(w).Encode(dtos.ErrorResponse{Error: "Invalid item ID"})
		return
	}

	err = repository.DeleteUserItem(uint(itemID), userID)
	if err != nil {
		http.Error(w, "Failed to delete user_item", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.ErrorResponse{Error: "Failed to delete user_item"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Search user items
// @Description Searches for user items by name
// @Tags user_item
// @Accept json
// @Produce json
// @Param name query string false "Name of user item"
// @Success 200 {object} dtos.UserItemsResponse
// @Failure default {object} dtos.ErrorResponse "Standard Error Responses"
// @Router /user_item/search [get]
func SearchUserItemsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.IDKey).(uint)

	query := dtos.UserItemQuery{}
	query.Name = r.URL.Query().Get("name")
	userItems, err := repository.SearchUserItems(query, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.InternalServerErrorResponse{Error: "Database error"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userItems)
}

// @Summary Predict items from an uploaded image
// @Description Predict items from an uploaded image for the authenticated user
// @Tags user_item
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Image file"
// @Success 200 {object} dtos.UserItemsResponse
// @Failure default {object} dtos.ErrorResponse "Standard Error Responses"
// @Router /user_item/predict [post]
func PredictUserItemsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.IDKey).(uint)

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		json.NewEncoder(w).Encode(dtos.ErrorResponse{Error: "Failed to parse form"})
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Image file is required", http.StatusBadRequest)
		json.NewEncoder(w).Encode(dtos.ErrorResponse{Error: "Image file is required"})
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read image", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.ErrorResponse{Error: "Failed to read image"})
		return
	}

	detectDomain := os.Getenv("HUGGINGFACE_URL")
	if detectDomain == "" {
		http.Error(w, "HUGGINGFACE_URL is not set", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.ErrorResponse{Error: "HUGGINGFACE_URL is not set"})
		return
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	filePart, err := writer.CreateFormFile("image", "upload.png")
	if err != nil {
		http.Error(w, "Failed to create form file", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.ErrorResponse{Error: "Failed to create form file"})
		return
	}

	if _, err := filePart.Write(fileBytes); err != nil {
		http.Error(w, "Failed to write file data", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.ErrorResponse{Error: "Failed to write file data"})
		return
	}

	writer.Close()

	predictURL := fmt.Sprintf("%s/predict", detectDomain)
	req, err := http.NewRequest("POST", predictURL, &buf)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.ErrorResponse{Error: "Failed to create request"})
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to get prediction", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.ErrorResponse{Error: "Failed to get prediction"})
		return
	}
	defer resp.Body.Close()

	var response struct {
		Items          []string `json:"items"`
		AnnotatedImage string   `json:"annotated_image"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		http.Error(w, "Failed to parse prediction response", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.ErrorResponse{Error: "Failed to parse prediction response"})
		return
	}

	var result []models.UserItem

	for _, class := range response.Items {
		if class == "" {
			continue
		}

		var existingItem models.Item
		if err := config.DB.Where("name LIKE ?", "%"+class+"%").First(&existingItem).Error; err == nil {
		} else {
			newItem := models.Item{
				Name:  class,
				Image: "",
			}
			config.DB.Create(&newItem)
			existingItem = newItem
		}

		userItem := models.UserItem{
			UserID: userID,
			ItemID: existingItem.ID,
		}
		result = append(result, userItem)

		var existingUserItem models.UserItem
		if err := config.DB.Where("user_id = ? AND item_id = ?", userItem.UserID, userItem.ItemID).First(&existingUserItem).Error; err != nil {
			config.DB.Create(&userItem)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
