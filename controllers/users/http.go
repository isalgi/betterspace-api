package users

import (
	"backend/app/middlewares"
	
	"backend/businesses/users"

	ctrl "backend/controllers"
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
		return ctrl.NewInfoResponse(c, http.StatusBadRequest, "failed", "invalid request")
	}

	err := userInput.Validate()

	if err != nil {
		return ctrl.NewInfoResponse(c, http.StatusBadRequest, "failed", "validation failed")
	}

	if userInput.Password != userInput.ConfirmationPassword {
		return ctrl.NewInfoResponse(c, http.StatusBadRequest, "failed", "password and confirmation password do not match")
	}

	user := ac.authUsecase.Register(userInput.ToDomainRegister())

	if user.ID == 0 {
		return ctrl.NewInfoResponse(c, http.StatusBadRequest, "failed", "email already taken. please use another email or process to login.")
	}

	return ctrl.NewResponse(c, http.StatusCreated, "success", "account created", response.FromDomain(user))
}

func (ac *AuthController) Login(c echo.Context) error {
	userInput := request.UserLogin{}

	if err := c.Bind(&userInput); err != nil {
		return ctrl.NewInfoResponse(c, http.StatusBadRequest, "failed", "invalid request")
	}

	err := userInput.Validate()

	if err != nil {
		return ctrl.NewInfoResponse(c, http.StatusBadRequest, "failed", "validation failed")
	}

	token := ac.authUsecase.Login(userInput.ToDomainLogin())

	if token == "" {
		return ctrl.NewInfoResponse(c, http.StatusUnauthorized, "failed", "invalid email or password")
	}

	return c.JSON(http.StatusOK, map[string]any{
		"token": token,
	})
}

func (ac *AuthController) Logout(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)

	isListed := middlewares.CheckToken(user.Raw)

	if !isListed {
		return ctrl.NewInfoResponse(c, http.StatusUnauthorized, "failed", "invalid token")
	}

	middlewares.Logout(user.Raw)

	return c.JSON(http.StatusOK, map[string]any{
		"message": "logout success",
	})
}