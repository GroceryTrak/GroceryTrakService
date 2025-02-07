package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/GroceryTrak/GroceryTrakService/internal/repository"
	"github.com/GroceryTrak/GroceryTrakService/internal/models"
	"github.com/GroceryTrak/GroceryTrakService/internal/utils"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err = repository.RegisterUser(user.Username, user.Password, user.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.Write([]byte("User registered successfully"))
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	username, role, err := repository.AuthenticateUser(user.Username, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(username, role)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(token))
}
