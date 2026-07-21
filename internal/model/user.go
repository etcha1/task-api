package model

import (
	"time"
)

type User struct {
	ID        int       `db:"id" json:"id"`
	Email     string    `db:"email" json:"email" validate:"required,email"`
	Password  string    `db:"password_hash" json:"password" validate:"required"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
