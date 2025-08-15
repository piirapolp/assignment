package controller

import (
	"assignment/global"
	model_mysql "assignment/model/mysql"
	"context"
)

type GetUserInput struct {
	UserId string `json:"user_id" validate:"required"`
}

type GetUserOutput struct {
	UserInfo model_mysql.User `json:"user_info"`
}

func (controller Controller) GetUser(ctx context.Context, input GetUserInput) (GetUserOutput, error) {
	controller.Logger.Info("start get user")
	output := GetUserOutput{}

	var err error
	output.UserInfo, err = controller.ModelRepository.GetUser(ctx, input.UserId)
	if err != nil {
		controller.Logger.Errorf("get user failed because: %s", err.Error())
		return output, global.SystemError{
			Code:    global.DatabaseError,
			Message: err.Error(),
		}
	}

	controller.Logger.Info("get user completed")
	return output, nil
}
