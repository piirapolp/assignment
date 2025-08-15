package controller

import (
	mock_model "assignment/mocks/model"
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"

	"assignment/global"
	model_mysql "assignment/model/mysql"
)

func TestGetUserAccounts_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_model.NewMockModelRepository(ctrl)

	// Arrange
	ctx := context.Background()
	wantAccounts := []model_mysql.AccountWithDetails{
		{}, {}, {},
	}

	// If your controller calls these, allow them freely; safe no-ops if not called.
	repo.EXPECT().ConfigureRequestId(gomock.AssignableToTypeOf((*string)(nil))).AnyTimes()
	repo.EXPECT().ConfigureUserId(gomock.AssignableToTypeOf((*string)(nil))).AnyTimes()

	// Expect the main call.
	repo.EXPECT().
		GetUserAccounts(gomock.Any(), gomock.Any()).
		Return(wantAccounts, nil).
		Times(1)

	c := newTestController(repo)
	out, err := c.GetUserAccounts(ctx)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !reflect.DeepEqual(out.Accounts, wantAccounts) {
		t.Fatalf("unexpected accounts.\nwant: %#v\ngot:  %#v", wantAccounts, out.Accounts)
	}
}

func TestGetUserAccounts_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_model.NewMockModelRepository(ctrl)

	ctx := context.Background()
	wantErr := errors.New("db fail")

	// Allow these if your controller sets them.
	repo.EXPECT().ConfigureRequestId(gomock.AssignableToTypeOf((*string)(nil))).AnyTimes()
	repo.EXPECT().ConfigureUserId(gomock.AssignableToTypeOf((*string)(nil))).AnyTimes()

	// The repository returns an error for GetUserAccounts.
	repo.EXPECT().
		GetUserAccounts(gomock.Any(), gomock.Any()).
		Return(nil, wantErr).
		Times(1)

	c := newTestController(repo)

	out, err := c.GetUserAccounts(ctx)
	if err == nil {
		t.Fatalf("expected error, got nil (output: %#v)", out)
	}

	var se global.SystemError
	if !errors.As(err, &se) {
		t.Fatalf("expected error of type global.SystemError, got %T", err)
	}
	if se.Code != global.DatabaseError {
		t.Fatalf("expected code %v, got %v", global.DatabaseError, se.Code)
	}
	if se.Message != wantErr.Error() {
		t.Fatalf("expected message %q, got %q", wantErr.Error(), se.Message)
	}
}

func TestGetUserDebitCards_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_model.NewMockModelRepository(ctrl)

	ctx := context.Background()

	// Cards returned by repository (unmasked).
	repoCards := []model_mysql.CardsWithDetails{
		{Number: "1234 5678 9012 3456"},
		{Number: "5555 6666 7777 8888 9999"},
		{Number: "11112222"}, // no spaces, unchanged
		{Number: "12 34"},    // only two parts, unchanged
	}

	// Expected numbers after masking in controller.
	wantCards := []model_mysql.CardsWithDetails{
		{Number: "1234 **** **** 3456"},
		{Number: "5555 **** **** **** 9999"},
		{Number: "11112222"},
		{Number: "12 34"},
	}

	repo.EXPECT().ConfigureRequestId(gomock.AssignableToTypeOf((*string)(nil))).AnyTimes()
	repo.EXPECT().ConfigureUserId(gomock.AssignableToTypeOf((*string)(nil))).AnyTimes()
	repo.EXPECT().
		GetUserCards(gomock.Any(), gomock.Any()).
		Return(repoCards, nil).
		Times(1)

	c := newTestController(repo)

	out, err := c.GetUserDebitCards(ctx)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(out.DebitCards) != len(wantCards) {
		t.Fatalf("unexpected number of cards. want %d, got %d", len(wantCards), len(out.DebitCards))
	}

	if !reflect.DeepEqual(out.DebitCards, wantCards) {
		t.Fatalf("unexpected cards after masking.\nwant: %#v\ngot:  %#v", wantCards, out.DebitCards)
	}
}

func TestGetUserDebitCards_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_model.NewMockModelRepository(ctrl)

	ctx := context.Background()
	wantErr := errors.New("db fail")

	repo.EXPECT().ConfigureRequestId(gomock.AssignableToTypeOf((*string)(nil))).AnyTimes()
	repo.EXPECT().ConfigureUserId(gomock.AssignableToTypeOf((*string)(nil))).AnyTimes()
	repo.EXPECT().
		GetUserCards(gomock.Any(), gomock.Any()).
		Return(nil, wantErr).
		Times(1)

	c := newTestController(repo)

	out, err := c.GetUserDebitCards(ctx)
	if err == nil {
		t.Fatalf("expected error, got nil (output: %#v)", out)
	}

	var se global.SystemError
	if !errors.As(err, &se) {
		t.Fatalf("expected error of type global.SystemError, got %T", err)
	}
	if se.Code != global.DatabaseError {
		t.Fatalf("expected code %v, got %v", global.DatabaseError, se.Code)
	}
	if se.Message != wantErr.Error() {
		t.Fatalf("expected message %q, got %q", wantErr.Error(), se.Message)
	}
}

func TestGetUserSavedAccounts_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_model.NewMockModelRepository(ctrl)

	ctx := context.Background()
	wantSaved := []model_mysql.SavedAccounts{
		{}, {}, {},
	}

	repo.EXPECT().ConfigureRequestId(gomock.AssignableToTypeOf((*string)(nil))).AnyTimes()
	repo.EXPECT().ConfigureUserId(gomock.AssignableToTypeOf((*string)(nil))).AnyTimes()
	repo.EXPECT().GetUserSavedAccounts(gomock.Any(), gomock.Any()).Return(wantSaved, nil).Times(1)

	c := newTestController(repo)

	out, err := c.GetUserSavedAccounts(ctx)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !reflect.DeepEqual(out.SavedAccounts, wantSaved) {
		t.Fatalf("unexpected saved accounts.\nwant: %#v\ngot:  %#v", wantSaved, out.SavedAccounts)
	}
}

func TestGetUserSavedAccounts_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_model.NewMockModelRepository(ctrl)

	ctx := context.Background()
	wantErr := errors.New("db fail")

	repo.EXPECT().ConfigureRequestId(gomock.AssignableToTypeOf((*string)(nil))).AnyTimes()
	repo.EXPECT().ConfigureUserId(gomock.AssignableToTypeOf((*string)(nil))).AnyTimes()
	repo.EXPECT().GetUserSavedAccounts(gomock.Any(), gomock.Any()).Return(nil, wantErr).Times(1)

	c := newTestController(repo)

	out, err := c.GetUserSavedAccounts(ctx)
	if err == nil {
		t.Fatalf("expected error, got nil (output: %#v)", out)
	}

	var se global.SystemError
	if !errors.As(err, &se) {
		t.Fatalf("expected error of type global.SystemError, got %T", err)
	}
	if se.Code != global.DatabaseError {
		t.Fatalf("expected code %v, got %v", global.DatabaseError, se.Code)
	}
	if se.Message != wantErr.Error() {
		t.Fatalf("expected message %q, got %q", wantErr.Error(), se.Message)
	}
}
