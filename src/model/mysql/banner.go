package model_mysql

import (
	"assignment/datastore/mysql"
	"assignment/entity"
	"context"
	"errors"
	"gorm.io/gorm"
)

func (repository *ModelMysqlRepository) GetUserBanners(ctx context.Context, userId string) ([]entity.Banners, error) {
	var result []entity.Banners
	if err := mysql.DB.WithContext(ctx).Where("user_id = ?", userId).Find(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []entity.Banners{}, errors.New("banner not found")
		} else {
			return []entity.Banners{}, err
		}
	}
	return result, nil
}
