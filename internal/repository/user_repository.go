package repository

import (
	"context"
	"fmt"

	"github.com/etcha1/task-api/internal/model"
	"github.com/etcha1/task-api/internal/utils"
	"github.com/jackc/pgx/v5"
)

// UserRepository handles database operations for Users.
type UserRepository struct {
	db *pgx.Conn // The connection handle is stored here
}

// NewUserRepository acts as a constructor to inject the DB connection.
func NewUserRepository(db *pgx.Conn) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) CreateUser(ctx context.Context, user model.User) (bool, error) {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return false, fmt.Errorf("createUser: %v", err)
	}
	user.Password = hashedPassword

	result, err := ur.db.Exec(ctx, "INSERT INTO users (email, password_hash) VALUES ($1, $2)", user.Email, user.Password)
	if err != nil {
		return false, fmt.Errorf("createUser: %v", err)
	}

	return result.Insert(), nil
}
