package trace

import (
	"time"

	"assignment/global"
	"assignment/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func New() fiber.Handler {
	return func(context *fiber.Ctx) (err error) {

		requestTime := time.Now()
		path := string(context.Request().URI().Path())
		query := string(context.Request().URI().QueryString())

		err = context.Next()

		// Trace Request
		latency := time.Since(requestTime)
		requestId := context.Locals(global.KEY_REQUEST_ID)
		logger.Logger.Infow(
			path,
			zap.String("request-id", requestId.(string)),
			zap.String("method", string(context.Request().Header.Method())),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("user-agent", string(context.Request().Header.UserAgent())),
			zap.String("ip", context.IP()),
			zap.Int("status", context.Response().StatusCode()),
			zap.Duration("latency", latency),
		)

		return err
	}
}
