package global

import "github.com/gofiber/fiber/v2"

const (
	BASE_SERVICE_NAME       = "Assignment"
	BASE_SERVICE_SHORT_NAME = "assignment"

	HEADER_REQUEST_ID = "Request-Id"
	KEY_REQUEST_ID    = "request_id"
	KEY_USER_ID       = "user_id"
	KEY_LOGGER        = "logger"
	KEY_PART          = "part"

	PART_INTERFACE  = "interface"
	PART_CONTROLLER = "controller"
	PART_MODEL      = "model"

	METHOD_GET  = "GET"
	METHOD_POST = "POST"

	RESULT_SUCCESS = "success"
)

type HandlerFunc func(*fiber.Ctx) error
