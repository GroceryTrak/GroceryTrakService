package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/GroceryTrak/GroceryTrakService/internal/repository"
	"github.com/GroceryTrak/GroceryTrakService/internal/templates"
	"github.com/go-chi/chi/v5"
)

// @Summary Get an item
// @Description Get an item by its ID
// @Tags item
// @Accept json
// @Produce json
// @Param id path int true "Item ID"
// @Success 200 {object} templates.ItemResponse
// @Failure default {object} templates.ErrorResponse "Standard Error Responses"
// @Router /item/{id} [get]
func GetItemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	item, err := repository.GetItem(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(templates.NotFoundResponse{Error: "Item not found"})
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
// @Param item body templates.ItemRequest true "New Item"
// @Success 201 {object} templates.ItemResponse
// @Failure default {object} templates.ErrorResponse "Standard Error Responses"
// @Router /item [post]
func CreateItemHandler(w http.ResponseWriter, r *http.Request) {
	var newItem templates.ItemRequest
	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(templates.BadRequestResponse{Error: "Invalid request data"})
		return
	}

	createdItem, err := repository.CreateItem(newItem)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(templates.InternalServerErrorResponse{Error: "Failed to create item"})
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
// @Param item body templates.ItemRequest true "Updated Item Data"
// @Success 200 {object} templates.ItemResponse
// @Failure default {object} templates.ErrorResponse "Standard Error Responses"
// @Router /item/{id} [put]
func UpdateItemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(templates.BadRequestResponse{Error: "Invalid item ID"})
		return
	}

	var updatedItem templates.ItemRequest
	if err := json.NewDecoder(r.Body).Decode(&updatedItem); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(templates.BadRequestResponse{Error: "Invalid request data"})
		return
	}

	item, err := repository.UpdateItem(uint(id), updatedItem)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(templates.NotFoundResponse{Error: "Item not found"})
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
// @Failure default {object} templates.ErrorResponse "Standard Error Responses"
// @Router /item/{id} [delete]
func DeleteItemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(templates.BadRequestResponse{Error: "Invalid item ID"})
		return
	}

	if err := repository.DeleteItem(uint(id)); err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(templates.NotFoundResponse{Error: "Item not found"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Search items
// @Description Searches for items that match the provided keyword in their name or description
// @Tags item
// @Accept json
// @Produce json
// @Param q query string true "Search keyword"
// @Success 200 {object} templates.ItemsResponse
// @Failure default {object} templates.ErrorResponse "Standard Error Responses"
// @Router /item/search [get]
func SearchItemsHandler(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("q")
	if keyword == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(templates.BadRequestResponse{Error: "Search keyword is required"})
		return
	}

	items, err := repository.SearchItems(keyword)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(templates.InternalServerErrorResponse{Error: "Database error"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}
