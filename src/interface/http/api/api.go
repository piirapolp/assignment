package api

import (
	v1 "assignment/interface/http/api/v1"
	"github.com/gofiber/fiber/v2"
)

func AddPublicRoute(router *fiber.Router) {
	v1Route := (*router).Group("/v1")
	v1.AddPublicRoutes(&v1Route)
}

func AddProtectedRoute(router *fiber.Router) {
	v1Route := (*router).Group("/v1")
	v1.AddProtectedRoutes(&v1Route)
}
