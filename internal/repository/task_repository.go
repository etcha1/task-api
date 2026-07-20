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
