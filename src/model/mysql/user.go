package model_mysql

import (
	"assignment/datastore/mysql"
	"assignment/entity"
	"context"
	"errors"
	"gorm.io/gorm"
)

type User struct {
	Name      string `json:"name"`
	DummyCol1 string `json:"dummy_col1"`
}

func (repository *ModelMysqlRepository) GetUser(ctx context.Context, userId string) (User, error) {
	var user entity.Users
	if err := mysql.DB.WithContext(ctx).Where("user_id = ?", userId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return User{}, errors.New("user not found")
		} else {
			return User{}, err
		}
	}

	return User{Name: user.Name, DummyCol1: user.DummyCol1}, nil
}
