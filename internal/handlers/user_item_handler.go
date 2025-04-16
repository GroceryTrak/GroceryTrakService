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

	"github.com/GroceryTrak/GroceryTrakService/internal/dtos"
	"github.com/GroceryTrak/GroceryTrakService/internal/middlewares"
	"github.com/GroceryTrak/GroceryTrakService/internal/repository"
	"github.com/go-chi/chi/v5"
)

type UserItemHandler struct {
	Repo repository.UserItemRepository
}

func NewUserItemHandler(repo repository.UserItemRepository) *UserItemHandler {
	return &UserItemHandler{Repo: repo}
}

// @Summary Get all user's items
// @Description Get all items for the authenticated user
// @Tags user_item
// @Produce json
// @Success 200 {object} dtos.UserItemsResponse
// @Failure default {object} dtos.ErrorResponse "Standard Error Responses"
// @Router /user_item [get]
func (h *UserItemHandler) GetAllUserItemsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.IDKey).(uint)

	userItems, err := h.Repo.GetAllUserItems(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.InternalServerErrorResponse{Error: "Failed to get all user items"})
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
func (h *UserItemHandler) GetUserItemHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.IDKey).(uint)

	itemID, err := strconv.ParseUint(chi.URLParam(r, "item_id"), 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dtos.BadRequestResponse{Error: "Invalid user item ID"})
		return
	}

	userItem, err := h.Repo.GetUserItem(uint(itemID), userID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(dtos.NotFoundResponse{Error: "Failed to get user item"})
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
func (h *UserItemHandler) CreateUserItemHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.IDKey).(uint)

	var req dtos.UserItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dtos.BadRequestResponse{Error: "Invalid request payload"})
		return
	}

	userItem, err := h.Repo.CreateUserItem(req, userID)
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
func (h *UserItemHandler) UpdateUserItemHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.IDKey).(uint)

	itemID, err := strconv.ParseUint(chi.URLParam(r, "item_id"), 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dtos.BadRequestResponse{Error: "Invalid item ID"})
		return
	}

	var req dtos.UserItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dtos.BadRequestResponse{Error: "Invalid request payload"})
		return
	}

	userItem, err := h.Repo.UpdateUserItem(req, uint(itemID), userID)
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
func (h *UserItemHandler) DeleteUserItemHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.IDKey).(uint)

	itemID, err := strconv.ParseUint(chi.URLParam(r, "item_id"), 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dtos.BadRequestResponse{Error: "Invalid item ID"})
		return
	}

	err = h.Repo.DeleteUserItem(uint(itemID), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.InternalServerErrorResponse{Error: "Failed to delete user item"})
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
func (h *UserItemHandler) SearchUserItemsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.IDKey).(uint)

	query := dtos.UserItemQuery{}
	query.Name = r.URL.Query().Get("name")
	userItems, err := h.Repo.SearchUserItems(query, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.InternalServerErrorResponse{Error: "Failed to search user items"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userItems)
}

// @Summary Detect items from an uploaded image
// @Description Detect items from an uploaded image for the authenticated user using OpenAI's vision model
// @Tags user_item
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Image file"
// @Success 200 {array} dtos.UserItemsResponse
// @Failure default {object} dtos.ErrorResponse "Standard Error Responses"
// @Router /user_item/detect [post]
func (h *UserItemHandler) DetectUserItemsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.IDKey).(uint)

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.InternalServerErrorResponse{Error: "OPENAI_API_KEY is not set"})
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10 MB max file size
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dtos.BadRequestResponse{Error: "Failed to parse form"})
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dtos.BadRequestResponse{Error: "Image file is required"})
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.InternalServerErrorResponse{Error: "Failed to read image"})
		return
	}

	userItems, err := h.Repo.DetectUserItems(fileBytes, userID, apiKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.InternalServerErrorResponse{Error: err.Error()})
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
func (h *UserItemHandler) PredictUserItemsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.IDKey).(uint)

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dtos.BadRequestResponse{Error: "Failed to parse form"})
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dtos.BadRequestResponse{Error: "Image file is required"})
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.InternalServerErrorResponse{Error: "Failed to read image"})
		return
	}

	detectDomain := os.Getenv("HUGGINGFACE_URL")
	if detectDomain == "" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.InternalServerErrorResponse{Error: "HUGGINGFACE_URL is not set"})
		return
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	filePart, err := writer.CreateFormFile("image", "upload.png")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.InternalServerErrorResponse{Error: "Failed to create form file"})
		return
	}

	if _, err := filePart.Write(fileBytes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.InternalServerErrorResponse{Error: "Failed to write file data"})
		return
	}

	writer.Close()

	predictURL := fmt.Sprintf("%s/predict", detectDomain)
	req, err := http.NewRequest("POST", predictURL, &buf)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.InternalServerErrorResponse{Error: "Failed to create request"})
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.InternalServerErrorResponse{Error: "Failed to get prediction"})
		return
	}
	defer resp.Body.Close()

	var response struct {
		Items          []string `json:"items"`
		AnnotatedImage string   `json:"annotated_image"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.InternalServerErrorResponse{Error: "Failed to parse prediction response"})
		return
	}

	result, err := h.Repo.PredictUserItems(response.Items, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.InternalServerErrorResponse{Error: "Failed to process prediction results"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
