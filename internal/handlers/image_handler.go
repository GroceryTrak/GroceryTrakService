package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/GroceryTrak/GroceryTrakService/internal/dtos"
)

// @Summary Proxy an image
// @Description Proxy an image from a given URL
// @Tags image
// @Produce image/*
// @Param url query string true "Image URL"
// @Success 200 {file} binary
// @Failure default {object} dtos.ErrorResponse "Standard Error Responses"
// @Router /image [get]
func ImageProxyHandler(w http.ResponseWriter, r *http.Request) {
	imageUrl := r.URL.Query().Get("url")
	if imageUrl == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dtos.BadRequestResponse{Error: "Missing 'url' query parameter"})
		return
	}

	parsedUrl, err := url.ParseRequestURI(imageUrl)
	if err != nil || (parsedUrl.Scheme != "http" && parsedUrl.Scheme != "https") {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dtos.BadRequestResponse{Error: "Invalid image URL"})
		return
	}

	resp, err := http.Get(imageUrl)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dtos.BadRequestResponse{Error: "Failed to fetch image"})
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.Header().Set("Access-Control-Allow-Origin", "*")

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dtos.InternalServerErrorResponse{Error: "Failed to stream image"})
		return
	}
}
