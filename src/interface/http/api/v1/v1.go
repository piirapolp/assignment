package v1

import (
	"assignment/global"
	"github.com/gofiber/fiber/v2"
)

var methodRoutesPublic = map[string]map[string]global.HandlerFunc{
	global.METHOD_GET:  make(map[string]global.HandlerFunc),
	global.METHOD_POST: make(map[string]global.HandlerFunc),
}

var methodRoutesProtected = map[string]map[string]global.HandlerFunc{
	global.METHOD_GET:  make(map[string]global.HandlerFunc),
	global.METHOD_POST: make(map[string]global.HandlerFunc),
}

func RegisterPublicGET(path string, h global.HandlerFunc) {
	methodRoutesPublic[global.METHOD_GET][path] = h
}
func RegisterPublicPOST(path string, h global.HandlerFunc) {
	methodRoutesPublic[global.METHOD_POST][path] = h
}
func RegisterProtectedGET(path string, h global.HandlerFunc) {
	methodRoutesProtected[global.METHOD_GET][path] = h
}
func RegisterProtectedPOST(path string, h global.HandlerFunc) {
	methodRoutesProtected[global.METHOD_POST][path] = h
}

func AddPublicRoutes(router *fiber.Router) {
	for route, h := range methodRoutesPublic[global.METHOD_GET] {
		(*router).Get(route, h)
	}
	for route, h := range methodRoutesPublic[global.METHOD_POST] {
		(*router).Post(route, h)
	}
}

func AddProtectedRoutes(router *fiber.Router) {
	for route, h := range methodRoutesProtected[global.METHOD_GET] {
		(*router).Get(route, h)
	}
	for route, h := range methodRoutesProtected[global.METHOD_POST] {
		(*router).Post(route, h)
	}
}
