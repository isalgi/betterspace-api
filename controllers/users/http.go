package users

import (
	"backend/app/middlewares"
	"backend/helper"

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

func (ac *AuthController) GetAll(c echo.Context) error {
	users := []response.User{}

	payload := helper.GetPayloadInfo(c)
	role := payload.Roles
	
	if role != "admin" {
		return ctrl.NewInfoResponse(c, http.StatusForbidden, "failed", "forbidden")
	}

	usersData := ac.authUsecase.GetAll()

	for _, user := range usersData {
		users = append(users, response.FromDomain(user))
	}

	return ctrl.NewResponse(c, http.StatusOK, "success", "all users", users)
}

func (ac *AuthController) GetByID(c echo.Context) error {
	payload := helper.GetPayloadInfo(c)
	role := payload.Roles
	userId := payload.ID
	
	id := c.Param("id")

	if (role == "user") && (id != userId) {
		return ctrl.NewInfoResponse(c, http.StatusForbidden, "failed", "forbidden")
	}

	user := ac.authUsecase.GetByID(id)

	if user.ID == 0 {
		return ctrl.NewInfoResponse(c, http.StatusNotFound, "failed", "user not found")
	}

	return ctrl.NewResponse(c, http.StatusOK, "success", "user found", response.FromDomain(user))
}

func (ac *AuthController) Delete(c echo.Context) error {
	payload := helper.GetPayloadInfo(c)
	role := payload.Roles
	userId := payload.ID
	
	id := c.Param("id")

	if (role == "user") && (id != userId) {
		return ctrl.NewInfoResponse(c, http.StatusForbidden, "failed", "forbidden")
	}

	isSuccess := ac.authUsecase.Delete(id)

	if !isSuccess {
		return ctrl.NewInfoResponse(c, http.StatusNotFound, "failed", "user not found")
	}

	return ctrl.NewInfoResponse(c, http.StatusOK, "success", "user deleted")
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