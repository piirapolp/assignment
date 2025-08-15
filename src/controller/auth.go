package controller

import (
	"assignment/global"
	"assignment/util"
	"context"
)

type LoginInput struct {
	UserId string `json:"user_id" validate:"required"`
	Pin    string `json:"pin" validate:"required"`
}

type LoginOutput struct {
	Greeting string `json:"greeting"`
	Token    string `json:"token"`
}

func (controller Controller) Login(ctx context.Context, input LoginInput) (LoginOutput, error) {
	controller.Logger.Info("start logging in")
	output := LoginOutput{}

	userPin, err := controller.ModelRepository.GetUserHashedPin(ctx, input.UserId)
	if err != nil {
		controller.Logger.Errorf("get user hashed pin failed because: %s", err.Error())
		return output, global.SystemError{
			Code:    global.DatabaseError,
			Message: err.Error(),
		}
	}

	if same, err := util.ValidatePin(input.Pin, userPin.Pin); !same {
		if err != nil {
			controller.Logger.Errorf("cannot validate pin for user %s", input.UserId)
			return output, global.SystemError{
				Code:    global.DatabaseError,
				Message: err.Error(),
			}
		}
		controller.Logger.Errorf("user %s just input an incorrect password", input.UserId)
		return output, global.SystemError{
			Code:    global.IncorrectPin,
			Message: global.GetErrorMessage(global.IncorrectPin),
		}
	}

	output.Token, output.Greeting, err = controller.ModelRepository.RevokeExistingTokenAndCreateNewToken(ctx, input.UserId)
	if err != nil {
		controller.Logger.Errorf("create token failed because: %s", err.Error())
		return output, global.SystemError{
			Code:    global.DatabaseError,
			Message: err.Error(),
		}
	}

	controller.Logger.Info("login completed")
	return output, nil
}
