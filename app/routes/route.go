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

	offices := e.Group("/api/v1/offices", middleware.JWTWithConfig(cl.JWTMiddleware))

	offices.GET("/all", cl.OfficeController.GetAll).Name = "get-all-type-of-offices"
	offices.GET("/:id", cl.OfficeController.GetByID).Name = "get-office-by-id"
	offices.POST("/create", cl.OfficeController.Create).Name = "create-office"
	offices.PUT("/update/:office_id", cl.OfficeController.Update).Name = "update-office"
	offices.DELETE("/delete/:office_id", cl.OfficeController.Delete).Name = "delete-office"
	offices.GET("/city/:city", cl.OfficeController.SearchByCity).Name = "group-office-by-city"
	offices.GET("/rate/:rate", cl.OfficeController.SearchByRate).Name = "group-office-by-rate"
	offices.GET("/title", cl.OfficeController.SearchByTitle).Name = "search-office-by-title"
	offices.POST("/images", cl.OfficeImageController.Create).Name = "create-office-image-list"
	offices.GET("/facilities", cl.OfficeFacilityController.GetAll).Name = "get-all-office-facility"
	offices.GET("/facilities/:id", cl.OfficeFacilityController.GetByOfficeID).Name = "get-office-facility-by-id"
	offices.POST("/facilities/create", cl.OfficeFacilityController.Create).Name = "create-office-facility-list"
	offices.GET("/type/office", cl.OfficeController.GetOffices).Name = "get-offices"
	offices.GET("/type/coworking-space", cl.OfficeController.GetCoworkingSpace).Name = "get-coworking-spaces"
	offices.GET("/type/meeting-room", cl.OfficeController.GetMeetingRooms).Name = "get-meeting-rooms"
	offices.GET("/recommendation", cl.OfficeController.GetRecommendation).Name = "recommendation-offices"
	offices.GET("/nearest", cl.OfficeController.GetNearest).Name = "get-nearest-building"

	facilities := e.Group("/api/v1/facilities", middleware.JWTWithConfig(cl.JWTMiddleware))

	facilities.GET("", cl.FacilityController.GetAll).Name = "get-all-facility"
	facilities.GET("/:id", cl.FacilityController.GetByID).Name = "get-facility-by-id"
	facilities.POST("", cl.FacilityController.Create).Name = "create-facility"
	facilities.PUT("/:id", cl.FacilityController.Update).Name = "update-facility"
	facilities.DELETE("/:id", cl.FacilityController.Delete).Name = "delete-facility"

	transactions := e.Group("/api/v1/transactions", middleware.JWTWithConfig(cl.JWTMiddleware))

	transactions.GET("", cl.TransactionController.GetAll).Name = "get-all-transaction"
	transactions.POST("", cl.TransactionController.Create).Name = "create-transaction"
	transactions.GET("/:id", cl.TransactionController.GetByID).Name = "get-transaction-by-id"
	transactions.PUT("/:id", cl.TransactionController.Update).Name = "update-transaction"
	transactions.DELETE("/:id", cl.TransactionController.Delete).Name = "delete-transaction"

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