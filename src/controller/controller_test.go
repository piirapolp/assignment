package controller

import (
	"assignment/logger"
	fake_logger "assignment/mocks/logger"
	mock_model "assignment/mocks/model"
	"assignment/model"
	"github.com/golang/mock/gomock"
	"testing"
)

func newTestController(repo model.ModelRepository) Controller {
	return Controller{
		Logger:          fake_logger.NewLogger(),
		ModelRepository: repo,
		UserId:          "test-user-id",
	}
}

func TestNew_InitializesControllerAndConfiguresRepo(t *testing.T) {
	origLogger := logger.Logger
	logger.Logger = fake_logger.NewLogger()
	defer func() { logger.Logger = origLogger }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_model.NewMockModelRepository(ctrl)

	reqID := "req-123"
	userID := "user-456"

	gomock.InOrder(
		mockRepo.EXPECT().ConfigureRequestId(&reqID).Times(1),
		mockRepo.EXPECT().ConfigureUserId(&userID).Times(1),
	)

	c := New(&reqID, &userID, mockRepo)

	if c.RequestId != reqID {
		t.Fatalf("RequestId mismatch: got %q, want %q", c.RequestId, reqID)
	}
	if c.UserId != userID {
		t.Fatalf("UserId mismatch: got %q, want %q", c.UserId, userID)
	}
	if c.ModelRepository != mockRepo {
		t.Fatalf("ModelRepository not set correctly")
	}
	if c.Logger == nil {
		t.Fatalf("Logger should be initialized")
	}
}
