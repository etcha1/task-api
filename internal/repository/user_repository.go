package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/etcha1/task-api/internal/model"
	"github.com/etcha1/task-api/internal/utils"
	"github.com/jackc/pgx/v5"
)

// UserRepository handles database operations for Users.
type UserRepository struct {
	db queryExecutor
}

// NewUserRepository acts as a constructor to inject the DB connection.
func NewUserRepository(db queryExecutor) *UserRepository {
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

func (ur *UserRepository) GetUser(ctx context.Context, userData model.User) (*model.User, error) {
	singleRow, err := ur.db.Query(ctx, "SELECT id, email, password_hash, created_at FROM users WHERE email = $1", userData.Email)
	if err != nil {
		return nil, fmt.Errorf("getUser: %w", err)
	}
	defer singleRow.Close()

	user, err := pgx.CollectOneRow(singleRow, pgx.RowToStructByName[model.User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("getUser: %w", err)
	}

	checkPasswordHash := utils.CheckPasswordHash(userData.Password, user.Password)
	if !checkPasswordHash {
		return nil, fmt.Errorf("getUser: password does not match")
	}

	return &user, nil
}
