package model

import (
	"assignment/entity"
	model_mysql "assignment/model/mysql"
	"context"
)

type ModelRepository interface {
	ConfigureRequestId(requestId *string)
	ConfigureUserId(userId *string)

	GetUserHashedPin(ctx context.Context, userId string) (entity.UserPin, error)
	RevokeExistingTokenAndCreateNewToken(ctx context.Context, userId string) (string, string, error)

	GetUserBanners(ctx context.Context, userId string) ([]entity.Banners, error)
	GetUserAccounts(ctx context.Context, userId string) ([]model_mysql.AccountWithDetails, error)
	GetUserCards(ctx context.Context, userId string) ([]model_mysql.CardsWithDetails, error)
	GetUserSavedAccounts(ctx context.Context, userId string) ([]model_mysql.SavedAccounts, error)
	GetUser(ctx context.Context, userId string) (model_mysql.User, error)
}
