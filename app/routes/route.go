package routes

import (
	"backend/controllers/facilities"
	officefacilities "backend/controllers/office_facilities"
	officeimage "backend/controllers/office_images"
	"backend/controllers/offices"
	transactions "backend/controllers/transactions"
	"backend/controllers/users"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ControllerList struct {
	LoggerMiddleware         echo.MiddlewareFunc
	JWTMiddleware            middleware.JWTConfig
	AuthController           users.AuthController
	OfficeController         offices.OfficeController
	OfficeImageController    officeimage.OfficeImageController
	FacilityController       facilities.FacilityController
	OfficeFacilityController officefacilities.OfficeFacilityController
	TransactionController    transactions.TransactionController
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

	e.GET("", cl.AuthController.HelloMessage)
  
	v1 := e.Group("/api/v1")
  
	// endpoint login, register, access refresh token
	v1.GET("", cl.AuthController.HelloMessage)
  	v1.POST("/register", cl.AuthController.Register)
	v1.POST("/login", cl.AuthController.Login)
	v1.POST("/refresh", cl.AuthController.Token, middleware.JWTWithConfig(cl.JWTMiddleware))
	v1.POST("/logout", cl.AuthController.Logout, middleware.JWTWithConfig(cl.JWTMiddleware)).Name = "user-logout"

	// endpoint admin
	admin := v1.Group("/admin")

	// endpoint admin : manage user
	userAdmin := admin.Group("/users")
	userAdmin.GET("", cl.AuthController.GetAll, middleware.JWTWithConfig(cl.JWTMiddleware)).Name = "admin-get-all-user"
	userAdmin.GET("/:id", cl.AuthController.GetByID, middleware.JWTWithConfig(cl.JWTMiddleware)).Name = "admin-get-user-by-id"
	userAdmin.DELETE("/:id", cl.AuthController.Delete, middleware.JWTWithConfig(cl.JWTMiddleware)).Name = "admin-delete-user-account"
	userAdmin.PUT("/photo/:id", cl.AuthController.UpdateProfilePhoto, middleware.JWTWithConfig(cl.JWTMiddleware)).Name = "admin-update-user-profile-photo"
	userAdmin.PUT("/:id", cl.AuthController.UpdateProfileData, middleware.JWTWithConfig(cl.JWTMiddleware)).Name = "admin-update-user-profile-data"
	userAdmin.GET("/email", cl.AuthController.SearchByEmail, middleware.JWTWithConfig(cl.JWTMiddleware)).Name = "admin-search-user-by-email"

	// endpoint admin : manage office
	officeAdmin := admin.Group("/offices")
	officeAdmin.POST("/create", cl.OfficeController.Create, middleware.JWTWithConfig(cl.JWTMiddleware)).Name = "[admin]-create-office"
	officeAdmin.PUT("/update/:office_id", cl.OfficeController.Update, middleware.JWTWithConfig(cl.JWTMiddleware)).Name = "[admin]-update-office"
	officeAdmin.DELETE("/delete/:office_id", cl.OfficeController.Delete, middleware.JWTWithConfig(cl.JWTMiddleware)).Name = "[admin]-delete-office"
	officeAdmin.GET("/all", cl.OfficeController.GetAll, middleware.JWTWithConfig(cl.JWTMiddleware)).Name="[admin]-get-all-type-of-offices"

	// endpoint admin : manage facilities
	facilities := admin.Group("/facilities")
	facilities.GET("/all", cl.FacilityController.GetAll, middleware.JWTWithConfig(cl.JWTMiddleware)).Name = "[admin]-get-all-facility"
	facilities.GET("/:id", cl.FacilityController.GetByID, middleware.JWTWithConfig(cl.JWTMiddleware)).Name = "[admin]-get-facility-by-id"
	facilities.POST("/create", cl.FacilityController.Create, middleware.JWTWithConfig(cl.JWTMiddleware)).Name = "[admin]-create-facility"
	facilities.PUT("/update/:id", cl.FacilityController.Update, middleware.JWTWithConfig(cl.JWTMiddleware)).Name = "[admin]-update-facility"
	facilities.DELETE("/delete/:id", cl.FacilityController.Delete, middleware.JWTWithConfig(cl.JWTMiddleware)).Name = "[admin]-delete-facility"

	// endpoint user : profile access
	profile := v1.Group("/profile")
	profile.GET("", cl.AuthController.GetByID, middleware.JWTWithConfig(cl.JWTMiddleware)).Name = "get-user-by-id"
	profile.DELETE("", cl.AuthController.Delete, middleware.JWTWithConfig(cl.JWTMiddleware)).Name = "delete-user-account"
	profile.PUT("/photo", cl.AuthController.UpdateProfilePhoto, middleware.JWTWithConfig(cl.JWTMiddleware)).Name = "update-user-profile-photo"
	profile.PUT("", cl.AuthController.UpdateProfileData, middleware.JWTWithConfig(cl.JWTMiddleware)).Name = "update-user-profile-data"

	// endpoint user : offices access
	offices := v1.Group("/offices")
	offices.GET("/all", cl.OfficeController.GetAll, middleware.JWTWithConfig(cl.JWTMiddleware)).Name = "get-all-type-of-offices"
	offices.GET("/:office_id", cl.OfficeController.GetByID, middleware.JWTWithConfig(cl.JWTMiddleware)).Name = "get-office-by-id"
	offices.GET("/city/:city", cl.OfficeController.SearchByCity, middleware.JWTWithConfig(cl.JWTMiddleware)).Name = "group-office-by-city"
	offices.GET("/rate/:rate", cl.OfficeController.SearchByRate, middleware.JWTWithConfig(cl.JWTMiddleware)).Name = "group-office-by-rate"
	offices.GET("/title", cl.OfficeController.SearchByTitle, middleware.JWTWithConfig(cl.JWTMiddleware)).Name = "search-office-by-title"
	offices.GET("/facilities", cl.FacilityController.GetAll, middleware.JWTWithConfig(cl.JWTMiddleware)).Name = "get-all-office-facility"
	offices.GET("/facilities/:office_id", cl.OfficeFacilityController.GetByOfficeID, middleware.JWTWithConfig(cl.JWTMiddleware)).Name = "get-office-facility-by-id"
	offices.GET("/type/office", cl.OfficeController.GetOffices, middleware.JWTWithConfig(cl.JWTMiddleware)).Name = "get-offices"
	offices.GET("/type/coworking-space", cl.OfficeController.GetCoworkingSpace, middleware.JWTWithConfig(cl.JWTMiddleware)).Name = "get-coworking-spaces"
	offices.GET("/type/meeting-room", cl.OfficeController.GetMeetingRooms, middleware.JWTWithConfig(cl.JWTMiddleware)).Name = "get-meeting-rooms"
	offices.GET("/recommendation", cl.OfficeController.GetRecommendation, middleware.JWTWithConfig(cl.JWTMiddleware)).Name = "recommendation-offices"
	offices.GET("/nearest", cl.OfficeController.GetNearest, middleware.JWTWithConfig(cl.JWTMiddleware)).Name = "get-nearest-building"

	transactions := v1.Group("/transactions", middleware.JWTWithConfig(cl.JWTMiddleware))

	transactions.GET("", cl.TransactionController.GetAll).Name = "get-all-transaction"
	transactions.POST("", cl.TransactionController.Create).Name = "create-transaction"
	transactions.GET("/:id", cl.TransactionController.GetByID).Name = "get-transaction-by-id"
	transactions.PUT("/:id", cl.TransactionController.Update).Name = "update-transaction"
	transactions.DELETE("/:id", cl.TransactionController.Delete).Name = "delete-transaction"
}