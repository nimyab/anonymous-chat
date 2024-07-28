package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/nimyab/anonymous-chat/internal/config"
	"strconv"
	"time"
)

type jwtCustomClaims struct {
	ID uint `json:"id"`
	jwt.RegisteredClaims
}

func GetUserId(c echo.Context) uint {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	return claims.ID
}

func CreateTokens(userId uint) (accessToken string, refreshToken string, err error) {
	cfg := config.GetEnvConfig()

	accessTime, err := strconv.Atoi(cfg.AccessTime)
	if err != nil {
		panic(err)
	}
	refreshTime, err := strconv.Atoi(cfg.RefreshTime)
	if err != nil {
		panic(err)
	}

	accessClaims := &jwtCustomClaims{
		userId,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(accessTime))),
		},
	}
	refreshClaims := &jwtCustomClaims{
		userId,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(refreshTime))),
		},
	}

	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", "", err
	}
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}