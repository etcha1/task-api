package main

import (
	"context"
	"log"
	"net/http"

	"github.com/etcha1/task-api/internal/database"
	"github.com/etcha1/task-api/internal/handler"
	"github.com/etcha1/task-api/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get database connection
	db := database.GetConnection()
	defer db.Close(context.Background())

	userRepo := repository.NewUserRepository(db)

	// Initialize the router
	r := chi.NewRouter()
	handler.RegisterRoutes(r, userRepo)
	log.Println("Server starting on :3000...")
	http.ListenAndServe(":3000", r)
}
