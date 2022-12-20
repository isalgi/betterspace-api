package utils

import (
	"backend/app/middlewares"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func GetPayloadInfo(c echo.Context) *middlewares.JwtCustomClaims {
	token := c.Get("user").(*jwt.Token)
	payload := middlewares.GetPayload(token)

	return payload
}