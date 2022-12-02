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

	e.GET("/", cl.AuthController.HelloMessage)
  
  // endpoint login, register, access refresh token
	e.POST("/api/v1/register", cl.AuthController.Register)
	e.POST("/api/v1/login", cl.AuthController.Login)
	e.POST("/api/v1/refresh", cl.AuthController.Token, middleware.JWTWithConfig(cl.JWTMiddleware))

	// endpoint admin
	admin := e.Group("/api/v1/admin", middleware.JWTWithConfig(cl.JWTMiddleware))
	admin.GET("/users", cl.AuthController.GetAll).Name = "admin-get-all-user"
	admin.GET("/user/:id", cl.AuthController.GetByID).Name = "admin-get-user-by-id"
	admin.DELETE("/user/:id", cl.AuthController.Delete).Name = "admin-delete-user-account"
	admin.PUT("/user/profile-photo/:id", cl.AuthController.UpdateProfilePhoto).Name = "admin-update-user-profile-photo"
	admin.PUT("/user/:id", cl.AuthController.UpdateProfileData).Name = "admin-update-user-profile-data"
	admin.GET("/user", cl.AuthController.SearchByEmail).Name = "admin-search-user-by-email"

	// endpoint user
	users := e.Group("/api/v1/users", middleware.JWTWithConfig(cl.JWTMiddleware))
	users.GET("/:id", cl.AuthController.GetByID).Name = "get-user-by-id"
	users.DELETE("/:id", cl.AuthController.Delete).Name = "delete-user-account"
	users.PUT("/profile-photo/:id", cl.AuthController.UpdateProfilePhoto).Name = "update-user-profile-photo"
	users.PUT("/:id", cl.AuthController.UpdateProfileData).Name = "update-user-profile-data"

	// endpoint logout
	auth := e.Group("/api/v1", middleware.JWTWithConfig(cl.JWTMiddleware))
	auth.POST("/logout", cl.AuthController.Logout).Name = "user-logout"
}