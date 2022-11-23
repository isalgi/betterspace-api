package main

import (
	"fmt"
	_middlewares "backend/app/middlewares"
	_routes "backend/app/routes"
	_util "backend/utils"

	_driverFactory "backend/drivers"
	_dbDriver "backend/drivers/mysql"

	_userUseCase "backend/businesses/users"
	_userController "backend/controllers/users"

	"github.com/labstack/echo/v4"
)

const DEFAULT_PORT = "3000"

func main() {
	configDB := _dbDriver.ConfigDB{
		DB_USERNAME: _util.GetConfig("DB_USERNAME"),
		DB_PASSWORD: _util.GetConfig("DB_PASSWORD"),
		DB_HOST: _util.GetConfig("DB_HOST"),
		DB_PORT: _util.GetConfig("DB_PORT"),
		DB_NAME: _util.GetConfig("DB_NAME"),
	}

	db := configDB.InitDB()

	_dbDriver.DBMigrate(db)

	configJWT := _middlewares.ConfigJWT{
		SecretJWT: _util.GetConfig("JWT_SECRET_KEY"),
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

	var appPort string = fmt.Sprintf(":%s", DEFAULT_PORT)

	app.Logger.Fatal(app.Start(appPort))
}
