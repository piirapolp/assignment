package model_mysql

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestGetUserBanners_Success(t *testing.T) {
	mock, teardown := setupMockDB(t)
	defer teardown()

	repo := &ModelMysqlRepository{}
	userID := "test-user-id"

	query := "SELECT * FROM `banners` WHERE user_id = ?"

	rows := sqlmock.NewRows([]string{"user_id"}).
		AddRow(userID)

	mock.ExpectQuery(query).WithArgs(userID).WillReturnRows(rows)

	ctx := context.Background()
	res, err := repo.GetUserBanners(ctx, userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 1 {
		t.Fatalf("expected 1 banner, got %d", len(res))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet mock expectations: %v", err)
	}
}

func TestGetUserBanners_NotFound(t *testing.T) {
	mock, teardown := setupMockDB(t)
	defer teardown()

	repo := &ModelMysqlRepository{}
	userID := "missing-user-id"

	query := "SELECT * FROM `banners` WHERE user_id = ?"

	mock.ExpectQuery(query).WithArgs(userID).WillReturnError(gorm.ErrRecordNotFound)

	ctx := context.Background()
	res, err := repo.GetUserBanners(ctx, userID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != "banner not found" {
		t.Fatalf("expected 'banner not found', got %q", err.Error())
	}
	if len(res) != 0 {
		t.Fatalf("expected empty result, got %d", len(res))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet mock expectations: %v", err)
	}
}

func TestGetUserBanners_DBError(t *testing.T) {
	mock, teardown := setupMockDB(t)
	defer teardown()

	repo := &ModelMysqlRepository{}
	userID := "any-user-id"

	query := "SELECT * FROM `banners` WHERE user_id = ?"

	dbErr := errors.New("db error")
	mock.ExpectQuery(query).WithArgs(userID).WillReturnError(dbErr)

	ctx := context.Background()
	res, err := repo.GetUserBanners(ctx, userID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, dbErr) {
		t.Fatalf("expected %v, got %v", dbErr, err)
	}
	if len(res) != 0 {
		t.Fatalf("expected empty result, got %d", len(res))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet mock expectations: %v", err)
	}
}

func TestGetUserBanners_Success_ExtraColumns(t *testing.T) {
	mock, teardown := setupMockDB(t)
	defer teardown()

	repo := &ModelMysqlRepository{}
	userID := "user-with-banners"

	query := "SELECT * FROM `banners` WHERE user_id = ?"

	now := time.Now()
	rows := sqlmock.NewRows([]string{"user_id", "created_at", "updated_at"}).AddRow(userID, now, now)

	mock.ExpectQuery(query).WithArgs(userID).WillReturnRows(rows)

	ctx := context.Background()
	res, err := repo.GetUserBanners(ctx, userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 1 {
		t.Fatalf("expected 1 banner, got %d", len(res))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet mock expectations: %v", err)
	}
}
