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

	users := e.Group("/api/v1/users")
	users.POST("/register", cl.AuthController.Register).Name = "user-register"
	users.POST("/login", cl.AuthController.Login).Name = "user-login"

	auth := e.Group("/api/v1/users", middleware.JWTWithConfig(cl.JWTMiddleware))
	auth.POST("/logout", cl.AuthController.Logout).Name = "user-logout"
}