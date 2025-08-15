package http

import (
	"assignment/interface/http/response"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func Ping(context *fiber.Ctx) error {
	return context.JSON(
		response.ResponseOutput{
			Code:    0,
			Message: "Success",
			Data:    "pong",
		},
	)
}

func Version(context *fiber.Ctx) error {
	version := viper.GetString("Version")
	return context.JSON(
		response.ResponseOutput{
			Code:    0,
			Message: "Success",
			Data:    map[string]string{"version": version},
		},
	)
}
