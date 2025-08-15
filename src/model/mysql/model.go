package model_mysql

import (
	"assignment/global"
	"assignment/logger"
)

type ModelMysqlRepository struct {
	RequestId string
	UserId    string
	Logger    logger.LoggerIface
}

func NewModelRepository() *ModelMysqlRepository {
	modelObj := ModelMysqlRepository{
		UserId: "system",
		Logger: logger.Logger,
	}
	return &modelObj
}

func (repository *ModelMysqlRepository) ConfigureRequestId(requestId *string) {
	repository.RequestId = *requestId
	repository.Logger = logger.Logger.
		With(global.KEY_REQUEST_ID, *requestId).
		With(global.KEY_PART, global.PART_MODEL)
}

func (repository *ModelMysqlRepository) ConfigureUserId(userId *string) {
	repository.UserId = *userId
}
