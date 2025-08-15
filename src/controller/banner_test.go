package controller

import (
	"assignment/entity"
	"assignment/global"
	mock_model "assignment/mocks/model"
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestController_GetUserBanners_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_model.NewMockModelRepository(ctrl)

	expected := []entity.Banners{{}, {}}
	mockRepo.EXPECT().GetUserBanners(gomock.Any(), "test-user-id").Return(expected, nil).Times(1)

	c := newTestController(mockRepo)

	out, err := c.GetUserBanners(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(out.Banners, expected) {
		t.Fatalf("unexpected banners: got %+v, want %+v", out.Banners, expected)
	}
}

func TestController_GetUserBanners_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_model.NewMockModelRepository(ctrl)

	mockRepo.EXPECT().GetUserBanners(gomock.Any(), "test-user-id").Return(nil, errors.New("boom")).Times(1)

	c := newTestController(mockRepo)

	out, err := c.GetUserBanners(context.Background())
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
	if sysErr.Message != "boom" {
		t.Fatalf("unexpected error message: got %q, want %q", sysErr.Message, "boom")
	}
	if len(out.Banners) != 0 {
		t.Fatalf("expected no banners on error, got %+v", out.Banners)
	}
}
