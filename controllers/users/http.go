package users

import (
	"backend/app/middlewares"
	"backend/helper"
	"context"
	"fmt"
	"log"

	"backend/businesses/users"

	ctrl "backend/controllers"
	"backend/controllers/users/request"
	"backend/controllers/users/response"

	"net/http"

	passwordvalidator "github.com/wagslane/go-password-validator"
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

	// confirm password input
	if userInput.Password != userInput.ConfirmationPassword {
		return ctrl.NewInfoResponse(c, http.StatusBadRequest, "failed", "password and confirmation password do not match")
	}

	const minEntropyBits = 45
	err = passwordvalidator.Validate(userInput.Password, minEntropyBits)
	
	if err != nil {
		return ctrl.NewInfoResponse(c, http.StatusBadRequest, "failed", fmt.Sprintf("%s", err))
	}

	user := ac.authUsecase.Register(userInput.ToDomainRegister())

	if user.ID == 0 {
		return ctrl.NewInfoResponse(c, http.StatusBadRequest, "failed", "email already taken, please use another email or process to login")
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
		return ctrl.NewInfoResponse(c, http.StatusForbidden, "forbidden", "not allowed to access this info")
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
	
	paramsId := c.Param("id")

	if (role == "user") && (paramsId != userId) {
		return ctrl.NewInfoResponse(c, http.StatusForbidden, "forbidden", "not allowed to access this info")
	}

	user := ac.authUsecase.GetByID(paramsId)

	if user.ID == 0 {
		return ctrl.NewInfoResponse(c, http.StatusNotFound, "failed", "user not found")
	}

	return ctrl.NewResponse(c, http.StatusOK, "success", "user found", response.FromDomain(user))
}

func (ac *AuthController) Delete(c echo.Context) error {
	payload := helper.GetPayloadInfo(c)
	role := payload.Roles
	userId := payload.ID
	
	paramsId := c.Param("id")

	if (role == "user") && (paramsId != userId) {
		return ctrl.NewInfoResponse(c, http.StatusForbidden, "forbidden", "not allowed to access this info")
	}

	isSuccess := ac.authUsecase.Delete(paramsId)

	if !isSuccess {
		return ctrl.NewInfoResponse(c, http.StatusNotFound, "failed", "user not found")
	}

	return ctrl.NewInfoResponse(c, http.StatusOK, "success", "user deleted")
}

func (ac *AuthController) UpdateProfilePhoto(c echo.Context) error {
	payload := helper.GetPayloadInfo(c)
	role := payload.Roles
	userId := payload.ID
	
	paramsId := c.Param("id")

	if (role == "user") && (paramsId != userId) {
		return ctrl.NewInfoResponse(c, http.StatusForbidden, "forbidden", "not allowed to access this info, check user id parameter")
	}

	getUser := ac.authUsecase.GetByID(paramsId)

	if getUser.ID == 0 {
		return ctrl.NewInfoResponse(c, http.StatusNotFound, "failed", "user not found")
	}

	input := request.UserPhoto{}

	if err := c.Bind(&input); err != nil {
		return ctrl.NewInfoResponse(c, http.StatusBadRequest, "failed", "validation failed")
	}

	fileInput, err := c.FormFile("photo")

	// validating input
	switch err {
		case nil:
			// do nothing
		case http.ErrMissingFile:
			return ctrl.NewInfoResponse(c, http.StatusBadRequest, "failed", "no file attached")
		default:
			return ctrl.NewInfoResponse(c, http.StatusBadRequest, "failed", "bind failed")
	}

	isFileAllowed, isFileAllowedMessage := helper.IsFileAllowed(fileInput)

	if !isFileAllowed {
		return ctrl.NewInfoResponse(c, http.StatusBadRequest, "failed", isFileAllowedMessage)
	}

	src, err := fileInput.Open()
	
	if err != nil {
		return ctrl.NewInfoResponse(c, http.StatusBadRequest, "failed", "bind failed")
	}

	defer src.Close()
	
	ctx := context.Background()

	url, err := helper.CloudinaryUpload(ctx, src, paramsId)
	
	if err != nil {
		log.Println(err)
		return ctrl.NewInfoResponse(c, http.StatusConflict, "failed", "upload to cloudinary failed")
	}

	input.Photo = url

	if err != nil {
		return ctrl.NewInfoResponse(c, http.StatusBadRequest, "failed", "validation failed")
	}

	isSuccess := ac.authUsecase.UpdateProfilePhoto(paramsId, input.ToDomainPhoto())

	if !isSuccess {
		return ctrl.NewInfoResponse(c, http.StatusBadRequest, "failed", "failed to update")
	}

	return ctrl.NewInfoResponse(c, http.StatusOK, "success", "profile photo updated")
}

func (ac *AuthController) UpdateProfileData(c echo.Context) error {
	payload := helper.GetPayloadInfo(c)
	role := payload.Roles
	userId := payload.ID
	
	paramsId := c.Param("id")

	// preventing user from updating another user data
	if (role == "user") && (paramsId != userId) {
		return ctrl.NewInfoResponse(c, http.StatusForbidden, "forbidden", "not allowed to access this info, check user id parameter")
	}

	userData := ac.authUsecase.GetByID(paramsId)

	if userData.ID == 0 {
		return ctrl.NewInfoResponse(c, http.StatusNotFound, "failed", "user not found")
	}

	input := request.User{}

	if err := c.Bind(&input); err != nil {
		return ctrl.NewInfoResponse(c, http.StatusBadRequest, "failed", "validation failed")
	}

	// check if body request is filled or not
	if input.FullName == "" && input.Gender == "" && input.Email == "" {
		return ctrl.NewInfoResponse(c, http.StatusBadRequest, "failed", "validation failed, please input data in body request")
	}

	// if full_name in body request is null
	if input.FullName == "" {
		input.FullName = userData.FullName
	}

	// if gender in body request is null
	if input.Gender == "" {
		input.Gender = userData.Gender
	}

	// if email in body request is null
	if input.Email == "" {
		input.Email = userData.Email
	}

	// fill other entity with existed data
	input.Password = userData.Password
	input.ConfirmationPassword = userData.Password
	input.Roles = userData.Roles
	input.Photo = userData.Photo

	err := input.Validate()

	if err != nil {
		return ctrl.NewInfoResponse(c, http.StatusBadRequest, "failed", "validation failed, check body request")
	}

	user := ac.authUsecase.UpdateProfileData(paramsId, input.ToDomainRegister())

	return ctrl.NewResponse(c, http.StatusOK, "success", "profile updated", response.FromDomain(user))
}

func (ac *AuthController) Logout(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)

	isListed := middlewares.CheckToken(user.Raw)

	if !isListed {
		return ctrl.NewInfoResponse(c, http.StatusUnauthorized, "failed", "invalid token")
	}

	middlewares.Logout(user.Raw)

	return ctrl.NewInfoResponse(c, http.StatusOK, "success", "logout success")
}