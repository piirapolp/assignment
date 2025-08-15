package log

import (
	"assignment/global"
	"assignment/logger"
	zaplogger "assignment/logger/zap"
	"github.com/gofiber/fiber/v2"
)

func New() fiber.Handler {
	return func(context *fiber.Ctx) (err error) {
		requestId := ""
		requestIdContext := context.Locals(global.KEY_REQUEST_ID)
		if requestIdContext != nil {
			requestId = requestIdContext.(string)
		}

		zapLogger := logger.Logger.(*zaplogger.ZapLogger).GetLogger()
		context.Locals(
			global.KEY_LOGGER,
			zapLogger.With(global.KEY_REQUEST_ID, requestId, global.KEY_PART, global.PART_INTERFACE),
		)

		return context.Next()
	}
}
