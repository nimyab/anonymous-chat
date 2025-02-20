package auth

import (
	"github.com/labstack/echo/v4"
	"github.com/nimyab/anonymous-chat/internal/handlers/auth/dtos"
	"github.com/nimyab/anonymous-chat/internal/jwt"
	"log/slog"
	"net/http"
	"time"
)

const CookiesMaxAge = 30 * 24 * 60 * 60 * 1000 // 30 дней;

type AuthHandler struct {
	authService *AuthService
}

func NewAuthHandler(authService *AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (handler *AuthHandler) Login(c echo.Context) error {
	var dto dtos.UserLoginDto

	if err := c.Bind(&dto); err != nil {
		slog.Error(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(&dto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, accessToken, refreshToken, err := handler.authService.Login(&dto)
	if err != nil {
		slog.Error(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	refreshCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		MaxAge:   CookiesMaxAge,
	}
	c.SetCookie(refreshCookie)

	return c.JSON(http.StatusOK, echo.Map{
		"user":         user,
		"access_token": accessToken,
	})
}

func (handler *AuthHandler) Registration(c echo.Context) error {
	var dto dtos.UserRegistrationDto

	if err := c.Bind(&dto); err != nil {
		slog.Error(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(&dto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := handler.authService.Registration(&dto)
	if err != nil {
		slog.Error(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, user)
}

func (handler *AuthHandler) Logout(c echo.Context) error {
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}
	c.SetCookie(cookie)
	return c.JSON(http.StatusOK, echo.Map{})
}

func (handler *AuthHandler) Refresh(c echo.Context) error {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	newAccessToken, newRefreshToken, err := handler.authService.Refresh(refreshToken.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	refreshCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		HttpOnly: true,
		MaxAge:   CookiesMaxAge,
	}
	c.SetCookie(refreshCookie)

	return c.JSON(http.StatusOK, echo.Map{
		"access_token": newAccessToken,
	})
}

func (handler *AuthHandler) UserInfo(c echo.Context) error {
	userId := jwt.GetUserId(c)
	user, err := handler.authService.UserInfo(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, user)
}
