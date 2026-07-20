package model

import "time"

type Task struct {
	ID          int       `db:"id" json:"id"`
	UserID      int       `db:"user_id" json:"user_id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Completed   bool      `db:"completed" json:"completed"`
	DueDate     time.Time `db:"due_date" json:"due_date"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}
