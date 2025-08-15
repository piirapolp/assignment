package model_mysql

import (
	"assignment/datastore/mysql"
	"assignment/entity"
	"assignment/util"
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
)

func (repository *ModelMysqlRepository) GetUserHashedPin(ctx context.Context, userId string) (entity.UserPin, error) {
	var result entity.UserPin
	if err := mysql.DB.WithContext(ctx).Where("user_id = ?", userId).First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.UserPin{}, errors.New("user not found")
		} else {
			return entity.UserPin{}, err
		}
	}
	return result, nil
}

func (repository *ModelMysqlRepository) RevokeExistingTokenAndCreateNewToken(ctx context.Context, userId string) (string, string, error) {
	var token, greeting string

	err := mysql.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		existingToken := entity.Tokens{
			ExpiredAt: time.Now().Add(-(time.Second * 1)),
		}
		if err := tx.Where("user_id = ?", userId).Where("expired_at > ?", time.Now()).Updates(&existingToken).Error; err != nil {
			return err
		}

		token = util.GenerateTokenSessionId(userId)
		newToken := entity.Tokens{
			SessionId: token,
			UserId:    userId,
			ExpiredAt: time.Now().Add(time.Minute * 720),
		}
		if err := tx.Create(&newToken).Error; err != nil {
			return err
		}

		if err := tx.Table(entity.UserGreetings{}.TableName()).
			Select("greeting").Where("user_id = ?", userId).
			Row().Scan(&greeting); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return "", "", err // rollback
	}
	return token, greeting, nil
}
