package model

import "time"

type Task struct {
	ID          int       `db:"id" json:"id"`
	UserID      int       `db:"user_id" json:"user_id"`
	Title       string    `db:"title" json:"title" validate:"required"`
	Description string    `db:"description" json:"description" validate:"required"`
	Completed   bool      `db:"completed" json:"completed" validate:"boolean"`
	DueDate     time.Time `db:"due_date" json:"due_date" validate:"datetime"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}
