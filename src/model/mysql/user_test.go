package model_mysql

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"
	"testing"
)

func TestGetUser_Success(t *testing.T) {
	mock, teardown := setupMockDB(t)
	defer teardown()

	repo := &ModelMysqlRepository{}
	userID := "test-user-id"

	query := "SELECT * FROM `users` WHERE user_id = ? ORDER BY `users`.`user_id` LIMIT ?"

	rows := sqlmock.NewRows([]string{"name", "dummy_col_1"}).AddRow("Alice", "X123")

	mock.ExpectQuery(query).WithArgs(userID, 1).WillReturnRows(rows)

	ctx := context.Background()
	got, err := repo.GetUser(ctx, userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != "Alice" || got.DummyCol1 != "X123" {
		t.Fatalf("unexpected user: %+v", got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet mock expectations: %v", err)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	mock, teardown := setupMockDB(t)
	defer teardown()

	repo := &ModelMysqlRepository{}
	userID := "missing-user-id"

	query := "SELECT * FROM `users` WHERE user_id = ? ORDER BY `users`.`user_id` LIMIT ?"

	mock.ExpectQuery(query).WithArgs(userID, 1).WillReturnError(gorm.ErrRecordNotFound)

	ctx := context.Background()
	got, err := repo.GetUser(ctx, userID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != "user not found" {
		t.Fatalf("expected 'user not found', got %q", err.Error())
	}
	if (got != User{}) {
		t.Fatalf("expected zero-value user, got %+v", got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet mock expectations: %v", err)
	}
}

func TestGetUser_DBError(t *testing.T) {
	mock, teardown := setupMockDB(t)
	defer teardown()

	repo := &ModelMysqlRepository{}
	userID := "any-user-id"

	query := "SELECT * FROM `users` WHERE user_id = ? ORDER BY `users`.`user_id` LIMIT ?"

	dbErr := errors.New("db error")
	mock.ExpectQuery(query).WithArgs(userID, 1).WillReturnError(dbErr)

	ctx := context.Background()
	got, err := repo.GetUser(ctx, userID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, dbErr) {
		t.Fatalf("expected %v, got %v", dbErr, err)
	}
	if (got != User{}) {
		t.Fatalf("expected zero-value user, got %+v", got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet mock expectations: %v", err)
	}
}
