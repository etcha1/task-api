package repository

import (
	"context"
	"fmt"

	"github.com/etcha1/task-api/internal/model"
	"github.com/jackc/pgx/v5"
)

// TaskRepository handles database operations for Tasks.
type TaskRepository struct {
	db *pgx.Conn // The connection handle is stored here
}

// NewTaskRepository acts as a constructor to inject the DB connection.
func NewTaskRepository(db *pgx.Conn) *TaskRepository {
	return &TaskRepository{db: db}
}

func (tr *TaskRepository) GetTasks(ctx context.Context, userId int) ([]model.Task, error) {
	rows, err := tr.db.Query(ctx, "SELECT id, user_id, title, description, completed, due_date, created_at FROM tasks WHERE user_id = $1", userId)
	if err != nil {
		return nil, fmt.Errorf("getTasks: %w", err)
	}
	defer rows.Close()

	tasks, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Task])
	if err != nil {
		return nil, fmt.Errorf("getTasks: %w", err)
	}

	return tasks, nil
}

func (tr *TaskRepository) GetTaskByID(ctx context.Context, taskID string) (*model.Task, error) {
	row := tr.db.QueryRow(ctx, "SELECT id, user_id, title, description, completed, due_date, created_at FROM tasks WHERE id = $1", taskID)

	var task model.Task
	err := row.Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Completed, &task.DueDate, &task.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Task not found
		}
		return nil, fmt.Errorf("getTaskByID: %w", err)
	}

	return &task, nil
}

func (tr *TaskRepository) CreateTask(ctx context.Context, task *model.Task) error {
	_, err := tr.db.Exec(ctx, "INSERT INTO tasks (user_id, title, description, completed, due_date) VALUES ($1, $2, $3, $4, $5)",
		task.UserID, task.Title, task.Description, task.Completed, task.DueDate)
	if err != nil {
		return fmt.Errorf("createTask: %w", err)
	}
	return nil
}

func (tr *TaskRepository) UpdateTask(ctx context.Context, task *model.Task) error {
	_, err := tr.db.Exec(ctx, "UPDATE tasks SET title = $1, description = $2, completed = $3, due_date = $4 WHERE id = $5",
		task.Title, task.Description, task.Completed, task.DueDate, task.ID)
	if err != nil {
		return fmt.Errorf("updateTask: %w", err)
	}
	return nil
}

func (tr *TaskRepository) DeleteTask(ctx context.Context, taskID int) error {
	_, err := tr.db.Exec(ctx, "DELETE FROM tasks WHERE id = $1", taskID)
	if err != nil {
		return fmt.Errorf("deleteTask: %w", err)
	}
	return nil
}
