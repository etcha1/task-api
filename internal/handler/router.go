package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/etcha1/task-api/internal/auth"
	"github.com/etcha1/task-api/internal/model"
	"github.com/etcha1/task-api/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
)

func RegisterRoutes(r *chi.Mux, userRepo *repository.UserRepository, taskRepo *repository.TaskRepository) {
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Group(func(r chi.Router) {
		r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
			registerHandler(w, r, userRepo)
		})
		r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
			loginHandler(w, r, userRepo)
		})
	})
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(auth.TokenAuth))
		r.Use(jwtauth.Authenticator(auth.TokenAuth))
		r.Get("/tasks", func(w http.ResponseWriter, r *http.Request) {
			tasksHandler(w, r, taskRepo)
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

func loginHandler(w http.ResponseWriter, r *http.Request, userRepo *repository.UserRepository) {
	var userData model.User

	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON payload"})
		return
	}

	user, err := userRepo.GetUser(r.Context(), userData)
	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve user"})
		return
	}

	if user == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
		return
	}

	token, err := auth.NewToken(user.ID)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to issue token"})
		return
	}

	response := map[string]interface{}{
		"user":  user,
		"token": token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func tasksHandler(w http.ResponseWriter, r *http.Request, taskRepo *repository.TaskRepository) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	userID := claims["user_id"]

	tasks, err := taskRepo.GetTasks(r.Context(), int(userID.(float64)))
	if err != nil {
		log.Printf("Error retrieving tasks: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve tasks"})
		return
	}

	response := map[string]interface{}{
		"tasks": tasks,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
