package middlewares

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4/middleware"
)

var whiteList []string = make([]string, 5)

type JwtCustomClaims struct {
	ID int `json:"id"`
	jwt.StandardClaims
}

type ConfigJWT struct {
	SecretJWT 		string
	ExpiresDuration int
}

func (jwtConf *ConfigJWT) Init() middleware.JWTConfig {
	return middleware.JWTConfig{
		Claims: &JwtCustomClaims{},
		SigningKey: []byte(jwtConf.SecretJWT),
	}
}

func (jwtConf *ConfigJWT) GenerateToken(userID int) string {
	claims := JwtCustomClaims {
		userID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(int64(jwtConf.ExpiresDuration))).Unix(),
		},
	}

	//create token with claims
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, _ := t.SignedString([]byte(jwtConf.SecretJWT))
	whiteList = append(whiteList, token)

	return token
}

func CheckToken(token string) bool {
	for _, tkn := range whiteList {
		if tkn == token {
			return true
		}
	}

	return false
}

func GetUser(token *jwt.Token) *JwtCustomClaims {
	// user := c.Get("user").(*jwt.Token)
	isListed := CheckToken(token.Raw)
	
	if !isListed {
		return nil
	}

	claims := token.Claims.(*JwtCustomClaims)

	return claims
}

func Logout(token string) bool {
	for i, tkn := range whiteList {
		if tkn == token {
			whiteList = append(whiteList[:i], whiteList[i+1:]...)
		}
	}

	return true
}