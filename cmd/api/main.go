package main

import (
	"context"
	"log"

	"github.com/etcha1/task-api/internal/database"
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
}
