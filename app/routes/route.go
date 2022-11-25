package routes

import (
	"backend/controllers/users"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ControllerList struct {
	LoggerMiddleware echo.MiddlewareFunc
	JWTMiddleware    middleware.JWTConfig
	AuthController   users.AuthController
}

func (cl *ControllerList) RouteRegister(e *echo.Echo) {
	e.Use(cl.LoggerMiddleware)

	beginning := e.Group("/api/v1")
	beginning.POST("/register", cl.AuthController.Register).Name = "user-register"
	beginning.POST("/login", cl.AuthController.Login).Name = "user-login"

	admins := e.Group("/api/v1/admin", middleware.JWTWithConfig(cl.JWTMiddleware))
	admins.GET("/users", cl.AuthController.GetAll).Name = "get-all-user"

	auth := e.Group("/api/v1/users", middleware.JWTWithConfig(cl.JWTMiddleware))
	auth.POST("/logout", cl.AuthController.Logout).Name = "user-logout"
}