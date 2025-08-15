package controller

import (
	"assignment/global"
	fake_logger "assignment/mocks/logger"
	mock_model "assignment/mocks/model"
	model_mysql "assignment/model/mysql"
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestController_GetUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_model.NewMockModelRepository(ctrl)
	c := Controller{
		Logger:          fake_logger.NewLogger(),
		ModelRepository: mockRepo,
	}

	input := GetUserInput{UserId: "user-123"}
	expected := model_mysql.User{Name: "Alice", DummyCol1: "x"}

	mockRepo.EXPECT().GetUser(gomock.Any(), input.UserId).Return(expected, nil).Times(1)

	out, err := c.GetUser(context.Background(), input)
	if err != nil {
		t.Fatalf("GetUser returned error: %v", err)
	}
	if out.UserInfo != expected {
		t.Fatalf("unexpected user info: got %+v, want %+v", out.UserInfo, expected)
	}
}

func TestController_GetUser_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_model.NewMockModelRepository(ctrl)
	c := Controller{
		Logger:          fake_logger.NewLogger(),
		ModelRepository: mockRepo,
	}

	input := GetUserInput{UserId: "user-404"}
	underlying := errors.New("user not found")

	mockRepo.EXPECT().GetUser(gomock.Any(), input.UserId).Return(model_mysql.User{}, underlying).Times(1)

	out, err := c.GetUser(context.Background(), input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	sysErr, ok := err.(global.SystemError)
	if !ok {
		t.Fatalf("expected global.SystemError, got %T: %v", err, err)
	}
	if sysErr.Code != global.DatabaseError {
		t.Fatalf("unexpected error code: got %v, want %v", sysErr.Code, global.DatabaseError)
	}
	if sysErr.Message != underlying.Error() {
		t.Fatalf("unexpected error message: got %q, want %q", sysErr.Message, underlying.Error())
	}

	if (out != GetUserOutput{}) {
		t.Fatalf("expected zero GetUserOutput on error, got %+v", out)
	}
}
