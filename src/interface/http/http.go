package http

import (
	"assignment/datastore/mysql"
	"assignment/interface/http/middleware/auth"
	"os"

	"assignment/global"
	"assignment/interface/http/api"
	"assignment/interface/http/middleware/log"
	"assignment/interface/http/middleware/trace"
	"assignment/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

var AppServer *fiber.App
var methodRoutes map[string]map[string]global.HandlerFunc

func InitHttpServer() {
	logger.Logger.Infof("http server is initilized")

	// Config Port and Address
	httpPort := viper.GetString("Interface.Http.Port")
	AppServer = fiber.New()

	// Config Middleware
	AppServer.Use(trace.New())
	AppServer.Use(requestid.New(requestid.Config{
		Header:     global.HEADER_REQUEST_ID,
		ContextKey: global.KEY_REQUEST_ID,
		Generator: func() string {
			return uuid.New().String()
		},
	}))
	AppServer.Use(log.New())
	AppServer.Use(recover.New())

	// Config Default Path
	AddRoute()

	// Add Api Path
	apiGroupPublic := AppServer.Group("/api")
	api.AddPublicRoute(&apiGroupPublic)

	apiGroupProtected := AppServer.Group("/api")
	apiGroupProtected.Use("", auth.DBTokenAuth(mysql.DB))
	api.AddProtectedRoute(&apiGroupProtected)

	// Start Server
	logger.Logger.Infof("serving http at http://127.0.0.1:%s", httpPort)
	err := AppServer.Listen(":" + httpPort)
	if err != nil {
		logger.Logger.Infof("http server listen and serves failed")
		os.Exit(1)
	}
	logger.Logger.Infof("http server is started")

}

func ShutdownHttpServer() {
	logger.Logger.Infof("http server is shutting down")
	if err := AppServer.Shutdown(); err != nil {
		logger.Logger.Infof("http server shut down failed: %s", err)
		return
	}
	logger.Logger.Infof("http server shut down completed")
}

func AddRoute() {
	for method, routes := range methodRoutes {
		if method == global.METHOD_GET {
			for routeName, routeFunc := range routes {
				AppServer.Get(routeName, routeFunc)
			}
		} else if method == global.METHOD_POST {
			for routeName, routeFunc := range routes {
				AppServer.Post(routeName, routeFunc)
			}
		}
	}
}

func init() {
	methodRoutes = make(map[string]map[string]global.HandlerFunc)
	methodRoutes[global.METHOD_GET] = make(map[string]global.HandlerFunc)
	methodRoutes[global.METHOD_POST] = make(map[string]global.HandlerFunc)

	methodRoutes[global.METHOD_GET]["/ping"] = Ping
	methodRoutes[global.METHOD_GET]["/version"] = Version
}
