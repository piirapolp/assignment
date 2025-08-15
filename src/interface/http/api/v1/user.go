package v1

import (
	"assignment/controller"
	"assignment/global"
	"assignment/interface/http/response"
	model_mysql "assignment/model/mysql"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func GetUserById(context *fiber.Ctx) error {
	contextLogger := context.Locals(global.KEY_LOGGER)
	apiLogger := contextLogger.(*zap.SugaredLogger)
	apiLogger.Info("GetUserById")

	input := controller.GetUserInput{}
	output := response.ResponseOutput{}

	// Parse Json
	if err := context.BodyParser(&input); err != nil {
		apiLogger.Errorf("could not bind json body to get user because: %s", err.Error())
		output.Code = global.InvalidJSONString
		output.Message = err.Error()
		return context.Status(fiber.ErrBadRequest.Code).JSON(output)
	}

	// Validate
	validate := validator.New()

	// Validate the User struct
	err := validate.Struct(input)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		apiLogger.Errorf("validate json body failed on get user because: %s", errors)
		output.Code = global.InvalidJSONString
		output.Message = err.Error()
		return context.Status(fiber.ErrBadRequest.Code).JSON(output)
	}

	requestId := context.Locals(global.KEY_REQUEST_ID)
	requestIdStr := requestId.(string)

	controllerObj := controller.New(&requestIdStr, &input.UserId, model_mysql.NewModelRepository())

	// Get request-scoped context from Fiber and pass it down
	reqCtx := context.UserContext()

	result, err := controllerObj.GetUser(reqCtx, input)
	if err != nil {
		output.Code = err.(global.SystemError).Code
		output.Message = err.Error()
		if output.Code == global.IncorrectPin {
			return context.Status(fiber.ErrUnauthorized.Code).JSON(output)
		}
		return context.Status(fiber.ErrInternalServerError.Code).JSON(output)
	}

	output.Message = global.RESULT_SUCCESS
	output.Data = result

	return context.JSON(output)
}

func init() {
	RegisterPublicPOST("/get-user-by-id", GetUserById)
}
