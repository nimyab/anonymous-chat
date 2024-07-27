package auth

import (
	"github.com/labstack/echo/v4"
	"github.com/nimyab/anonymous-chat/internal/handlers/auth/dtos"
	"net/http"
)

type AuthHandler struct {
	authService *AuthService
}

func NewAuthHandler(authService *AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (handler *AuthHandler) Login(c echo.Context) error {
	var dto dtos.UserLoginDto

	if err := c.Bind(&dto); err != nil {
		// todo: add logger
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// todo: validation data
	if err := dto.Validation(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := handler.authService.Login(&dto)
	if err != nil {
		// todo: add logger
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, user)
}

func (handler *AuthHandler) Registration(c echo.Context) error {
	var dto dtos.UserRegistrationDto

	if err := c.Bind(&dto); err != nil {
		// todo: add logger
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// todo: validation data
	if err := dto.Validation(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := handler.authService.Registration(&dto)
	if err != nil {
		// todo: add logger
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, user)
}

func (handler *AuthHandler) Logout(c echo.Context) error {
	return nil
}
