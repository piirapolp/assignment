package controller

import (
	"assignment/global"
	"assignment/logger"
	"assignment/model"
)

type Controller struct {
	RequestId       string
	UserId          string
	Logger          logger.LoggerIface
	ModelRepository model.ModelRepository
}

func New(requestId, userId *string, modelRepository model.ModelRepository) Controller {
	controllerObj := Controller{
		RequestId:       *requestId,
		UserId:          *userId,
		ModelRepository: modelRepository,
	}

	controllerObj.Logger = logger.Logger.
		With(global.KEY_REQUEST_ID, *requestId).
		With(global.KEY_PART, global.PART_CONTROLLER)

	controllerObj.ModelRepository.ConfigureRequestId(requestId)
	controllerObj.ModelRepository.ConfigureUserId(userId)

	return controllerObj
}
