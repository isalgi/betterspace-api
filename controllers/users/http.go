package users

import (
	"backend/app/middlewares"
	"backend/businesses/users"
	"backend/controllers/users/request"
	"backend/controllers/users/response"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type AuthController struct {
	authUsecase users.Usecase
}

func NewAuthController(authUC users.Usecase) *AuthController {
	return &AuthController{
		authUsecase: authUC,
	}
}

func (ac *AuthController) Register(c echo.Context) error {
	userInput := request.User{}

	if err := c.Bind(&userInput); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "invalid request",
		})
	}

	err := userInput.Validate()

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "validation failed",
		})
	}

	if userInput.Password != userInput.ConfirmationPassword {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status":  "Bad request",
			"message": "The password and confirmation password do not match",
		})
	}

	user := ac.authUsecase.Register(userInput.ToDomainRegister())

	if user.ID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "email already taken",
		})
	}

	return c.JSON(http.StatusCreated, response.FromDomain(user))
}

func (ac *AuthController) Login(c echo.Context) error {
	userInput := request.UserLogin{}

	if err := c.Bind(&userInput); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "invalid request",
		})
	}

	err := userInput.Validate()

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "validation failed",
		})
	}

	token := ac.authUsecase.Login(userInput.ToDomainLogin())

	if token == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "invalid email or password",
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"token": token,
	})
}

func (ac *AuthController) Logout(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)

	isListed := middlewares.CheckToken(user.Raw)

	if !isListed {
		return c.JSON(http.StatusUnauthorized, map[string]any{
			"message": "invalid token",
		})
	}

	middlewares.Logout(user.Raw)

	return c.JSON(http.StatusOK, map[string]any{
		"message": "logout success",
	})
}