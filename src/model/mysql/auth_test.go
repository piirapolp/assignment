package model_mysql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
)

func TestGetUserHashedPin_Success(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := &ModelMysqlRepository{}
	ctx := context.Background()
	userID := "user-123"

	query := "SELECT * FROM `user_pin` WHERE user_id = ? ORDER BY `user_pin`.`user_id` LIMIT ?"

	rows := sqlmock.NewRows([]string{"user_id", "hashed_pin"}).
		AddRow(userID, "hashed-pin-placeholder")

	mock.ExpectQuery(query).
		WithArgs(userID, 1).
		WillReturnRows(rows)

	got, err := repo.GetUserHashedPin(ctx, userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Basic sanity check; other fields (not provided) will be zero-values
	if got.UserId != userID {
		t.Fatalf("expected user_id %q, got %q", userID, got.UserId)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetUserHashedPin_NotFound(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := &ModelMysqlRepository{}
	ctx := context.Background()
	userID := "missing-user"

	query := "SELECT * FROM `user_pin` WHERE user_id = ? ORDER BY `user_pin`.`user_id` LIMIT ?"

	mock.ExpectQuery(query).
		WithArgs(userID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "hashed_pin"}))

	_, err := repo.GetUserHashedPin(ctx, userID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != "user not found" {
		t.Fatalf("expected 'user not found', got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetUserHashedPin_DBError(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := &ModelMysqlRepository{}
	ctx := context.Background()
	userID := "user-err"

	query := "SELECT * FROM `user_pin` WHERE user_id = ? ORDER BY `user_pin`.`user_id` LIMIT ?"

	mock.ExpectQuery(query).WithArgs(userID, 1).WillReturnError(errors.New("db failure"))

	_, err := repo.GetUserHashedPin(ctx, userID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != "db failure" {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestRevokeExistingTokenAndCreateNewToken_Success(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := &ModelMysqlRepository{}
	ctx := context.Background()
	userID := "user-1"

	mock.ExpectBegin()

	updateQuery := "UPDATE `tokens` SET `expired_at`=? WHERE user_id = ? AND expired_at > ?"
	mock.ExpectExec(updateQuery).WillReturnResult(sqlmock.NewResult(0, 1))

	insertQuery := "INSERT INTO `tokens` (`session_id`,`user_id`,`issued_at`,`expired_at`) VALUES (?,?,?,?)"
	mock.ExpectExec(insertQuery).
		WithArgs(sqlmock.AnyArg(), userID, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	selectGreetingQuery := "SELECT greeting FROM `user_greetings` WHERE user_id = ?"
	mockGreeting := "Welcome back!"
	mock.ExpectQuery(selectGreetingQuery).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"greeting"}).AddRow(mockGreeting))

	mock.ExpectCommit()

	token, greeting, err := repo.RevokeExistingTokenAndCreateNewToken(ctx, userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token == "" {
		t.Fatalf("expected non-empty token")
	}
	if greeting != mockGreeting {
		t.Fatalf("expected greeting %q, got %q", mockGreeting, greeting)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestRevokeExistingTokenAndCreateNewToken_UpdateError(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := &ModelMysqlRepository{}
	ctx := context.Background()
	userID := "user-2"

	mock.ExpectBegin()

	updateQuery := "UPDATE `tokens` SET `expired_at`=? WHERE user_id = ? AND expired_at > ?"
	mock.ExpectExec(updateQuery).WillReturnError(errors.New("update failed"))

	mock.ExpectRollback()

	token, greeting, err := repo.RevokeExistingTokenAndCreateNewToken(ctx, userID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if token != "" || greeting != "" {
		t.Fatalf("expected empty outputs on error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestRevokeExistingTokenAndCreateNewToken_InsertError(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := &ModelMysqlRepository{}
	ctx := context.Background()
	userID := "user-3"

	mock.ExpectBegin()

	updateQuery := "UPDATE `tokens` SET `expired_at`=? WHERE user_id = ? AND expired_at > ?"
	mock.ExpectExec(updateQuery).WillReturnResult(sqlmock.NewResult(0, 1))

	insertQuery := "INSERT INTO `tokens` (`session_id`,`user_id`,`issued_at`,`expired_at`) VALUES (?,?,?,?)"
	mock.ExpectExec(insertQuery).
		WithArgs(sqlmock.AnyArg(), userID, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("insert failed"))

	mock.ExpectRollback()

	token, greeting, err := repo.RevokeExistingTokenAndCreateNewToken(ctx, userID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if token != "" || greeting != "" {
		t.Fatalf("expected empty outputs on error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestRevokeExistingTokenAndCreateNewToken_GreetingQueryError(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := &ModelMysqlRepository{}
	ctx := context.Background()
	userID := "user-4"

	mock.ExpectBegin()

	updateQuery := "UPDATE `tokens` SET `expired_at`=? WHERE user_id = ? AND expired_at > ?"
	mock.ExpectExec(updateQuery).WillReturnResult(sqlmock.NewResult(0, 1))

	insertQuery := "INSERT INTO `tokens` (`session_id`,`user_id`,`issued_at`,`expired_at`) VALUES (?,?,?,?)"
	mock.ExpectExec(insertQuery).
		WithArgs(sqlmock.AnyArg(), userID, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	selectGreetingQuery := "SELECT greeting FROM `user_greetings` WHERE user_id = ?"
	mock.ExpectQuery(selectGreetingQuery).
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	mock.ExpectRollback()

	token, greeting, err := repo.RevokeExistingTokenAndCreateNewToken(ctx, userID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if token != "" || greeting != "" {
		t.Fatalf("expected empty outputs on error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
