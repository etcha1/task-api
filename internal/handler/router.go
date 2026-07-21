// @title Task API
// @version 1.0
// @description API for managing users and tasks.
// @host localhost:3000
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/etcha1/task-api/internal/auth"
	"github.com/etcha1/task-api/internal/middleware"
	"github.com/etcha1/task-api/internal/model"
	"github.com/etcha1/task-api/internal/repository"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/etcha1/task-api/docs"
)

func RegisterRoutes(r *chi.Mux, userRepo *repository.UserRepository, taskRepo *repository.TaskRepository) {
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.With(middleware.ValidateBody[model.User]()).Group(func(r chi.Router) {
		r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
			registerHandler(w, r, userRepo)
		})
		r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
			loginHandler(w, r, userRepo)
		})
	})
	r.With(middleware.ValidateBody[model.Task]()).Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(auth.TokenAuth))
		r.Use(jwtauth.Authenticator(auth.TokenAuth))
		r.Get("/tasks", func(w http.ResponseWriter, r *http.Request) {
			tasksHandler(w, r, taskRepo)
		})
		r.Route("/{taskID}", func(r chi.Router) {
			r.Use(singleTaskHandler(taskRepo)) // Middleware to load task by ID
			r.Get("/", getTaskHandler)         // GET /task/123
			r.Put("/", func(w http.ResponseWriter, r *http.Request) {
				updateTaskHandler(w, r, taskRepo)
			}) // PUT /task/123
			r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
				deleteTaskHandler(w, r, taskRepo)
			}) // DELETE /task/123
			r.Patch("/complete", func(w http.ResponseWriter, r *http.Request) {
				patchTaskHandler(w, r, taskRepo)
			}) // PATCH /task/123/complete
		})
	})
	r.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:3000/docs/doc.json"), // Points to the generated json
	))
}

// RegisterUser godoc
// @Summary Register a new user
// @Description Create a new user account and return the created user.
// @Tags auth
// @Accept json
// @Produce json
// @Param user body model.User true "User registration payload"
// @Success 201 {object} model.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /register [post]
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

// LoginUser godoc
// @Summary Log in an existing user
// @Description Authenticate a user and return a JWT token.
// @Tags auth
// @Accept json
// @Produce json
// @Param user body model.User true "User login payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /login [post]
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

// ListTasks godoc
// @Summary List tasks
// @Description Retrieve all tasks for the authenticated user.
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /tasks [get]
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

func singleTaskHandler(taskRepo *repository.TaskRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			taskID := chi.URLParam(r, "taskID")
			task, err := taskRepo.GetTaskByID(r.Context(), taskID)
			if err != nil {
				log.Printf("Error retrieving task: %v", err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve task"})
				return
			}

			if task == nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{"error": "Task not found"})
				return
			}

			// Store the task in the request context for the next handlers to use
			ctx := context.WithValue(r.Context(), "task", task)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetTask godoc
// @Summary Get a task by ID
// @Description Fetch a single task for the authenticated user.
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param taskID path string true "Task ID"
// @Success 200 {object} model.Task
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /{taskID} [get]
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	task := r.Context().Value("task").(*model.Task)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

// UpdateTask godoc
// @Summary Update a task
// @Description Update an existing task for the authenticated user.
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param taskID path string true "Task ID"
// @Param task body model.Task true "Task update payload"
// @Success 200 {object} model.Task
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /{taskID} [put]
func updateTaskHandler(w http.ResponseWriter, r *http.Request, taskRepo *repository.TaskRepository) {
	task := r.Context().Value("task").(*model.Task)

	var updatedTask model.Task
	err := json.NewDecoder(r.Body).Decode(&updatedTask)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON payload"})
		return
	}

	err = taskRepo.UpdateTask(r.Context(), &updatedTask)
	if err != nil {
		log.Printf("Error updating task: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update task"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

// DeleteTask godoc
// @Summary Delete a task
// @Description Delete an existing task for the authenticated user.
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param taskID path string true "Task ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /{taskID} [delete]
func deleteTaskHandler(w http.ResponseWriter, r *http.Request, taskRepo *repository.TaskRepository) {
	task := r.Context().Value("task").(*model.Task)

	err := taskRepo.DeleteTask(r.Context(), task.ID)
	if err != nil {
		log.Printf("Error deleting task: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete task"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Task deleted successfully", "taskID": strconv.Itoa(task.ID)})
}

// CompleteTask godoc
// @Summary Mark a task as complete
// @Description Mark an existing task as completed for the authenticated user.
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param taskID path string true "Task ID"
// @Success 200 {object} model.Task
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /{taskID}/complete [patch]
func patchTaskHandler(w http.ResponseWriter, r *http.Request, taskRepo *repository.TaskRepository) {
	task := r.Context().Value("task").(*model.Task)
	task.Completed = true // Mark the task as completed

	err := taskRepo.UpdateTask(r.Context(), task)
	if err != nil {
		log.Printf("Error updating task: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update task"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}
