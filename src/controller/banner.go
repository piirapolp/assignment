package controller

import (
	"assignment/entity"
	"assignment/global"
	"context"
)

type GetBannersOutput struct {
	Banners []entity.Banners `json:"banners"`
}

func (controller Controller) GetUserBanners(ctx context.Context) (GetBannersOutput, error) {
	controller.Logger.Info("getting user banners")
	output := GetBannersOutput{}

	var err error
	output.Banners, err = controller.ModelRepository.GetUserBanners(ctx, controller.UserId)
	if err != nil {
		controller.Logger.Errorf("get user banners failed because: %s", err.Error())
		return output, global.SystemError{
			Code:    global.DatabaseError,
			Message: err.Error(),
		}
	}

	controller.Logger.Info("get user banners completed")

	return output, nil
}
