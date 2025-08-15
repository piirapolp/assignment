package model_mysql

import (
	datastoremysql "assignment/datastore/mysql"
	"assignment/logger"
	fake_logger "assignment/mocks/logger"
	"github.com/DATA-DOG/go-sqlmock"
	gorm_mysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func setupMockDB(t *testing.T) (sqlmock.Sqlmock, func()) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	g := gorm_mysql.New(gorm_mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	})
	db, err := gorm.Open(g, &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm with sqlmock: %v", err)
	}

	datastoremysql.DB = db

	teardown := func() {
		_ = sqlDB.Close()
		datastoremysql.DB = nil
	}

	return mock, teardown
}

func TestNewModelRepository_Defaults(t *testing.T) {
	origLogger := logger.Logger
	logger.Logger = fake_logger.NewLogger()
	defer func() { logger.Logger = origLogger }()

	repo := NewModelRepository()
	if repo == nil {
		t.Fatal("expected repo to be non-nil")
	}
	if repo.UserId != "system" {
		t.Fatalf("expected default UserId to be 'system', got %q", repo.UserId)
	}
	if repo.Logger == nil {
		t.Fatal("expected default Logger to be set (non-nil)")
	}
}

func TestConfigureRequestId_SetsFields(t *testing.T) {
	origLogger := logger.Logger
	logger.Logger = fake_logger.NewLogger()
	defer func() { logger.Logger = origLogger }()

	repo := NewModelRepository()
	rid := "req-123"

	repo.ConfigureRequestId(&rid)

	if repo.RequestId != rid {
		t.Fatalf("expected RequestId %q, got %q", rid, repo.RequestId)
	}
	if repo.Logger == nil {
		t.Fatal("expected Logger to remain non-nil after ConfigureRequestId")
	}
}

func TestConfigureUserId_SetsField(t *testing.T) {
	origLogger := logger.Logger
	logger.Logger = fake_logger.NewLogger()
	defer func() { logger.Logger = origLogger }()

	repo := NewModelRepository()
	uid := "user-xyz"

	repo.ConfigureUserId(&uid)

	if repo.UserId != uid {
		t.Fatalf("expected UserId %q, got %q", uid, repo.UserId)
	}
}
