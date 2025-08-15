package v1

import (
	"assignment/controller"
	"assignment/global"
	"assignment/interface/http/response"
	model_mysql "assignment/model/mysql"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func GetBanners(context *fiber.Ctx) error {
	contextLogger := context.Locals(global.KEY_LOGGER)
	apiLogger := contextLogger.(*zap.SugaredLogger)
	apiLogger.Info("GetBanners")

	output := response.ResponseOutput{}

	// Get user_id from token
	userId, _ := context.Locals(global.KEY_USER_ID).(string)

	// Validate User
	if userId == "" {
		apiLogger.Errorf("validate user failed on get banners because user_id is empty")
		output.Code = global.InvalidJSONString
		output.Message = global.GetErrorMessage(global.InvalidUserToken)
		return context.Status(fiber.ErrBadRequest.Code).JSON(output)
	}

	requestId := context.Locals(global.KEY_REQUEST_ID)
	requestIdStr := requestId.(string)

	controllerObj := controller.New(&requestIdStr, &userId, model_mysql.NewModelRepository())

	// Get request-scoped context from Fiber and pass it down
	reqCtx := context.UserContext()

	result, err := controllerObj.GetUserBanners(reqCtx)
	if err != nil {
		output.Code = err.(global.SystemError).Code
		output.Message = err.Error()
		return context.Status(fiber.ErrInternalServerError.Code).JSON(output)
	}

	output.Message = global.RESULT_SUCCESS
	output.Data = result

	return context.JSON(output)
}

func init() {
	RegisterProtectedGET("/get-user-banners", GetBanners)
}
