package auth

import (
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService *AuthService
}

func NewAuthHandler(authService *AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (handler *AuthHandler) Login(c echo.Context) error {
	return nil
}

func (handler *AuthHandler) Registration(c echo.Context) error {
	return nil
}

func (handler *AuthHandler) Logout(c echo.Context) error {
	return nil
}
