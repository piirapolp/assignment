package controller

import (
	"assignment/entity"
	"assignment/util"
	"context"
	"errors"
	"testing"

	"assignment/global"
	mock_model "assignment/mocks/model"
	"github.com/golang/mock/gomock"
)

func TestLogin_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_model.NewMockModelRepository(ctrl)
	c := newTestController(repo)
	ctx := context.Background()

	input := LoginInput{
		UserId: "user-123",
		Pin:    "123456",
	}

	pin, _ := util.HashPassword(input.Pin)
	userPin := entity.UserPin{Pin: pin}

	wantToken := "token-xyz"
	wantGreeting := "Welcome!"

	repo.EXPECT().ConfigureRequestId(gomock.AssignableToTypeOf((*string)(nil))).AnyTimes()
	repo.EXPECT().ConfigureUserId(gomock.AssignableToTypeOf((*string)(nil))).AnyTimes()
	repo.EXPECT().GetUserHashedPin(gomock.Any(), input.UserId).Return(userPin, nil).Times(1)
	repo.EXPECT().RevokeExistingTokenAndCreateNewToken(gomock.Any(), input.UserId).Return(wantToken, wantGreeting, nil).Times(1)

	out, err := c.Login(ctx, input)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if out.Token != wantToken {
		t.Fatalf("unexpected token: want %q, got %q", wantToken, out.Token)
	}
	if out.Greeting != wantGreeting {
		t.Fatalf("unexpected greeting: want %q, got %q", wantGreeting, out.Greeting)
	}
}

func TestLogin_GetUserHashedPinError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_model.NewMockModelRepository(ctrl)
	c := newTestController(repo)
	ctx := context.Background()

	input := LoginInput{
		UserId: "user-123",
		Pin:    "123456",
	}

	wantErr := errors.New("db fail on get user pin")

	repo.EXPECT().ConfigureRequestId(gomock.AssignableToTypeOf((*string)(nil))).AnyTimes()
	repo.EXPECT().ConfigureUserId(gomock.AssignableToTypeOf((*string)(nil))).AnyTimes()
	repo.EXPECT().GetUserHashedPin(gomock.Any(), input.UserId).Return(entity.UserPin{}, wantErr).Times(1)

	out, err := c.Login(ctx, input)
	if err == nil {
		t.Fatalf("expected error, got nil (output: %#v)", out)
	}
	var se global.SystemError
	if !errors.As(err, &se) {
		t.Fatalf("expected global.SystemError, got %T", err)
	}
	if se.Code != global.DatabaseError {
		t.Fatalf("expected code %v, got %v", global.DatabaseError, se.Code)
	}
	if se.Message != wantErr.Error() {
		t.Fatalf("expected message %q, got %q", wantErr.Error(), se.Message)
	}
}

func TestLogin_IncorrectPin(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_model.NewMockModelRepository(ctrl)
	c := newTestController(repo)
	ctx := context.Background()

	input := LoginInput{
		UserId: "user-123",
		Pin:    "123456",
	}

	// Hashed pin for a different PIN, ensuring mismatch.
	hashedDifferent, _ := util.HashPassword("000000")
	userPin := entity.UserPin{Pin: hashedDifferent}

	repo.EXPECT().ConfigureRequestId(gomock.AssignableToTypeOf((*string)(nil))).AnyTimes()
	repo.EXPECT().ConfigureUserId(gomock.AssignableToTypeOf((*string)(nil))).AnyTimes()
	repo.EXPECT().GetUserHashedPin(gomock.Any(), input.UserId).Return(userPin, nil).Times(1)

	out, err := c.Login(ctx, input)
	if err == nil {
		t.Fatalf("expected error, got nil (output: %#v)", out)
	}
	var se global.SystemError
	if !errors.As(err, &se) {
		t.Fatalf("expected global.SystemError, got %T", err)
	}
	if se.Code != global.IncorrectPin {
		t.Fatalf("expected code %v, got %v", global.IncorrectPin, se.Code)
	}
	if se.Message != global.GetErrorMessage(global.IncorrectPin) {
		t.Fatalf("expected message %q, got %q", global.GetErrorMessage(global.IncorrectPin), se.Message)
	}
}

func TestLogin_ValidatePinError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_model.NewMockModelRepository(ctrl)
	c := newTestController(repo)
	ctx := context.Background()

	input := LoginInput{
		UserId: "user-123",
		Pin:    "123456",
	}

	userPin := entity.UserPin{Pin: "not-a-valid-hash"}

	repo.EXPECT().ConfigureRequestId(gomock.AssignableToTypeOf((*string)(nil))).AnyTimes()
	repo.EXPECT().ConfigureUserId(gomock.AssignableToTypeOf((*string)(nil))).AnyTimes()
	repo.EXPECT().GetUserHashedPin(gomock.Any(), input.UserId).Return(userPin, nil).Times(1)

	out, err := c.Login(ctx, input)
	if err == nil {
		t.Fatalf("expected error, got nil (output: %#v)", out)
	}
	var se global.SystemError
	if !errors.As(err, &se) {
		t.Fatalf("expected global.SystemError, got %T", err)
	}
	if se.Code != global.DatabaseError {
		t.Fatalf("expected code %v, got %v", global.DatabaseError, se.Code)
	}
}

func TestLogin_CreateTokenError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_model.NewMockModelRepository(ctrl)
	c := newTestController(repo)
	ctx := context.Background()

	input := LoginInput{
		UserId: "user-123",
		Pin:    "123456",
	}

	pin, _ := util.HashPassword(input.Pin)
	userPin := entity.UserPin{Pin: pin}
	wantErr := errors.New("db fail on create token")

	repo.EXPECT().ConfigureRequestId(gomock.AssignableToTypeOf((*string)(nil))).AnyTimes()
	repo.EXPECT().ConfigureUserId(gomock.AssignableToTypeOf((*string)(nil))).AnyTimes()
	repo.EXPECT().GetUserHashedPin(gomock.Any(), input.UserId).Return(userPin, nil).Times(1)
	repo.EXPECT().RevokeExistingTokenAndCreateNewToken(gomock.Any(), input.UserId).Return("", "", wantErr).Times(1)

	out, err := c.Login(ctx, input)
	if err == nil {
		t.Fatalf("expected error, got nil (output: %#v)", out)
	}
	var se global.SystemError
	if !errors.As(err, &se) {
		t.Fatalf("expected global.SystemError, got %T", err)
	}
	if se.Code != global.DatabaseError {
		t.Fatalf("expected code %v, got %v", global.DatabaseError, se.Code)
	}
	if se.Message != wantErr.Error() {
		t.Fatalf("expected message %q, got %q", wantErr.Error(), se.Message)
	}
}
