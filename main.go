package main

import (
	_middlewares "backend/app/middlewares"
	_routes "backend/app/routes"
	"fmt"
	"os"

	_driverFactory "backend/drivers"
	_dbDriver "backend/drivers/mysql"

	_userUseCase "backend/businesses/users"
	_userController "backend/controllers/users"

	"github.com/labstack/echo/v4"
)

const DEFAULT_PORT = "3000"

func main() {
	configDB := _dbDriver.ConfigDB{
		DB_USERNAME: os.Getenv("DB_USERNAME"),
		DB_PASSWORD: os.Getenv("DB_PASSWORD"),
		DB_HOST: os.Getenv("DB_HOST"),
		DB_PORT: os.Getenv("DB_PORT"),
		DB_NAME: os.Getenv("DB_NAME"),
	}

	db := configDB.InitDB()

	_dbDriver.DBMigrate(db)

	configJWT := _middlewares.ConfigJWT{
		SecretJWT: os.Getenv("JWT_SECRET_KEY"),
		ExpiresDuration: 1,
	}

	configLogger := _middlewares.ConfigLogger{
		Format: "[${time_rfc3339}] ${status} ${method} ${host} ${path} ${latency_human}" + "\n",
	}

	app := echo.New()

	userRepo := _driverFactory.NewUserRepository(db)
	userUseCase := _userUseCase.NewUserUsecase(userRepo, &configJWT)
	userCtrl := _userController.NewAuthController(userUseCase)

	routesInit := _routes.ControllerList{
		LoggerMiddleware: configLogger.Init(),
		JWTMiddleware: configJWT.Init(),
		AuthController: *userCtrl,
	}

	routesInit.RouteRegister(app)

	var port string = os.Getenv("PORT")

	if port == "" {
		port = DEFAULT_PORT
	}

	var appPort string = fmt.Sprintf(":%s", DEFAULT_PORT)

	app.Logger.Fatal(app.Start(appPort))
}
