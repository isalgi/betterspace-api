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

	e.POST("/api/v1/register", cl.AuthController.Register)
	e.POST("/api/v1/login", cl.AuthController.Login)

	users := e.Group("/api/v1/users", middleware.JWTWithConfig(cl.JWTMiddleware))
	users.GET("", cl.AuthController.GetAll).Name = "get-all-user"
	users.GET("/:id", cl.AuthController.GetByID).Name = "get-user-by-id"
	users.DELETE("/:id", cl.AuthController.Delete).Name = "delete-user-account"
	users.PUT("/profile-photo/:id", cl.AuthController.UpdateProfilePhoto).Name = "update-user-profile-photo"
	users.PUT("/:id", cl.AuthController.UpdateProfileData).Name = "update-profile-data"

	auth := e.Group("/api/v1", middleware.JWTWithConfig(cl.JWTMiddleware))
	auth.POST("/logout", cl.AuthController.Logout).Name = "user-logout"
}