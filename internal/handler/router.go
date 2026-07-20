package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/etcha1/task-api/internal/model"
	"github.com/etcha1/task-api/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func RegisterRoutes(r *chi.Mux, userRepo *repository.UserRepository) {
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Group(func(r chi.Router) {
		r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
			registerHandler(w, r, userRepo)
		})
	})
}

func registerHandler(w http.ResponseWriter, r *http.Request, userRepo *repository.UserRepository) {
	// Handle the registration logic here
	var userData model.User

	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON payload"})
		return
	}

	isUserCreated, err := userRepo.CreateUser(r.Context(), userData)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create user"})
		return
	}

	// Set header and write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(isUserCreated)
}
