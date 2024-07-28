package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/nimyab/anonymous-chat/internal/config"
)

func Middleware() echo.MiddlewareFunc {
	cfg := config.GetEnvConfig()

	return echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(cfg.Secret),
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
	})
}
