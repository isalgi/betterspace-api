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
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			echo.HeaderServer,
		},
	}))
	
	e.Use(cl.LoggerMiddleware)

	e.GET("/", cl.AuthController.HelloMessage)
  
	v1 := e.Group("/api/v1")
  
	// endpoint login, register, access refresh token
	v1.GET("", cl.AuthController.HelloMessage)
  v1.POST("/register", cl.AuthController.Register)
	v1.POST("/login", cl.AuthController.Login)
	v1.POST("/refresh", cl.AuthController.Token, middleware.JWTWithConfig(cl.JWTMiddleware))
	v1.POST("/logout", cl.AuthController.Logout, middleware.JWTWithConfig(cl.JWTMiddleware)).Name = "user-logout"

	// endpoint admin
	admin := v1.Group("/admin", middleware.JWTWithConfig(cl.JWTMiddleware))

	// endpoint admin : manage user
	userAdmin := admin.Group("/users")
	userAdmin.GET("", cl.AuthController.GetAll).Name = "admin-get-all-user"
	userAdmin.GET("/:id", cl.AuthController.GetByID).Name = "admin-get-user-by-id"
	userAdmin.DELETE("/:id", cl.AuthController.Delete).Name = "admin-delete-user-account"
	userAdmin.PUT("/photo/:id", cl.AuthController.UpdateProfilePhoto).Name = "admin-update-user-profile-photo"
	userAdmin.PUT("/:id", cl.AuthController.UpdateProfileData).Name = "admin-update-user-profile-data"
	userAdmin.GET("/email", cl.AuthController.SearchByEmail).Name = "admin-search-user-by-email"

	// endpoint user
	profile := v1.Group("/profile", middleware.JWTWithConfig(cl.JWTMiddleware))
	profile.GET("", cl.AuthController.GetByID).Name = "get-user-by-id"
	profile.DELETE("", cl.AuthController.Delete).Name = "delete-user-account"
	profile.PUT("/photo", cl.AuthController.UpdateProfilePhoto).Name = "update-user-profile-photo"
	profile.PUT("", cl.AuthController.UpdateProfileData).Name = "update-user-profile-data"
}