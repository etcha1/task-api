package repository

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/etcha1/task-api/internal/model"
	"github.com/jackc/pgx/v5"
	pgxconn "github.com/jackc/pgx/v5/pgconn"
)

type stubDB struct {
	execSQL  string
	execArgs []any
	execErr  error
}

func (s *stubDB) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return nil, nil
}

func (s *stubDB) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return nil
}

func (s *stubDB) Exec(ctx context.Context, sql string, args ...any) (pgxconn.CommandTag, error) {
	s.execSQL = sql
	s.execArgs = append([]any(nil), args...)
	if s.execErr != nil {
		return pgxconn.CommandTag{}, s.execErr
	}
	return pgxconn.NewCommandTag("INSERT 1 1"), nil
}

func TestCreateUserHashesPasswordAndExecutesInsert(t *testing.T) {
	db := &stubDB{}
	repo := NewUserRepository(db)

	user := model.User{Email: "test@example.com", Password: "plain-password"}

	inserted, err := repo.CreateUser(context.Background(), user)
	if err != nil {
		t.Fatalf("CreateUser returned error: %v", err)
	}

	if !inserted {
		t.Fatal("CreateUser should report a successful insert")
	}

	if db.execSQL == "" {
		t.Fatal("expected Exec to be called")
	}

	if got := db.execArgs[0]; got != "test@example.com" {
		t.Fatalf("expected email argument %q, got %v", "test@example.com", got)
	}

	hashedPassword, ok := db.execArgs[1].(string)
	if !ok {
		t.Fatalf("expected password argument to be a string, got %T", db.execArgs[1])
	}

	if hashedPassword == "plain-password" {
		t.Fatal("expected password to be hashed before persistence")
	}

	if hashedPassword == "" {
		t.Fatal("expected password hash to be non-empty")
	}
}

func TestCreateTaskExecutesInsert(t *testing.T) {
	db := &stubDB{}
	repo := NewTaskRepository(db)

	task := &model.Task{
		UserID:      1,
		Title:       "Write tests",
		Description: "Add unit tests",
		Completed:   false,
		DueDate:     time.Date(2026, 7, 21, 12, 0, 0, 0, time.UTC),
	}

	if err := repo.CreateTask(context.Background(), task); err != nil {
		t.Fatalf("CreateTask returned error: %v", err)
	}

	if db.execSQL == "" {
		t.Fatal("expected Exec to be called")
	}

	if !strings.Contains(db.execSQL, "INSERT INTO tasks") {
		t.Fatalf("expected INSERT SQL, got %q", db.execSQL)
	}
}

func TestDeleteTaskExecutesDelete(t *testing.T) {
	db := &stubDB{}
	repo := NewTaskRepository(db)

	if err := repo.DeleteTask(context.Background(), 7); err != nil {
		t.Fatalf("DeleteTask returned error: %v", err)
	}

	if db.execSQL == "" {
		t.Fatal("expected Exec to be called")
	}

	if !strings.Contains(db.execSQL, "DELETE FROM tasks") {
		t.Fatalf("expected DELETE SQL, got %q", db.execSQL)
	}
}
