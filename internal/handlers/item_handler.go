package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/GroceryTrak/GroceryTrakService/internal/dtos"
	"github.com/GroceryTrak/GroceryTrakService/internal/repository"
	"github.com/go-chi/chi/v5"
)

type ItemHandler struct {
	Repo repository.ItemRepository
}

func NewItemHandler(repo repository.ItemRepository) *ItemHandler {
	return &ItemHandler{Repo: repo}
}

// @Summary Get an item
// @Description Get an item by its ID
// @Tags item
// @Accept json
// @Produce json
// @Param id path int true "Item ID"
// @Success 200 {object} dtos.ItemResponse
// @Failure default {object} dtos.ErrorResponse "Standard Error Responses"
// @Router /item/{id} [get]
func (h *ItemHandler) GetItemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	item, err := h.Repo.GetItem(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(dtos.NotFoundResponse{Error: "Item not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// @Summary Create an item
// @Description Add a new item to the database
// @Tags item
// @Accept json
// @Produce json
// @Param item body dtos.ItemRequest true "New Item"
// @Success 201 {object} dtos.ItemResponse
// @Failure default {object} dtos.ErrorResponse "Standard Error Responses"
// @Router /item [post]
func (h *ItemHandler) CreateItemHandler(w http.ResponseWriter, r *http.Request) {
	var newItem dtos.ItemRequest
	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dtos.BadRequestResponse{Error: "Invalid request data"})
		return
	}

	createdItem, err := h.Repo.CreateItem(newItem)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.InternalServerErrorResponse{Error: "Failed to create item"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdItem)
}

// @Summary Update an item
// @Description Update an existing item by ID
// @Tags item
// @Accept json
// @Produce json
// @Param id path int true "Item ID"
// @Param item body dtos.ItemRequest true "Updated Item Data"
// @Success 200 {object} dtos.ItemResponse
// @Failure default {object} dtos.ErrorResponse "Standard Error Responses"
// @Router /item/{id} [put]
func (h *ItemHandler) UpdateItemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dtos.BadRequestResponse{Error: "Invalid item ID"})
		return
	}

	var updatedItem dtos.ItemRequest
	if err := json.NewDecoder(r.Body).Decode(&updatedItem); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dtos.BadRequestResponse{Error: "Invalid request data"})
		return
	}

	item, err := h.Repo.UpdateItem(uint(id), updatedItem)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.InternalServerErrorResponse{Error: "Failed to update item"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// @Summary Delete an item
// @Description Remove an item by its ID
// @Tags item
// @Accept json
// @Produce json
// @Param id path int true "Item ID"
// @Success 204 "No Content"
// @Failure default {object} dtos.ErrorResponse "Standard Error Responses"
// @Router /item/{id} [delete]
func (h *ItemHandler) DeleteItemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dtos.BadRequestResponse{Error: "Invalid item ID"})
		return
	}

	if err := h.Repo.DeleteItem(uint(id)); err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(dtos.NotFoundResponse{Error: "Item not found"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Search items
// @Description Searches for items that match the provided keyword in their name or description
// @Tags item
// @Accept json
// @Produce json
// @Param name query string true "Search keyword"
// @Success 200 {object} dtos.ItemsResponse
// @Failure default {object} dtos.ErrorResponse "Standard Error Responses"
// @Router /item/search [get]
func (h *ItemHandler) SearchItemsHandler(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("name")
	if keyword == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dtos.BadRequestResponse{Error: "Search keyword is required"})
		return
	}

	items, err := h.Repo.SearchItems(keyword)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.InternalServerErrorResponse{Error: "Database error"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}
