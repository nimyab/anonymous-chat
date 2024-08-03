package chat

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/nimyab/anonymous-chat/internal/handlers/chat/dtos"
	"github.com/nimyab/anonymous-chat/internal/jwt"
	"net/http"
	"strconv"
)

type ChatHandler struct {
	chatService *ChatService
}

func NewChatHandler(chatService *ChatService) *ChatHandler {
	return &ChatHandler{chatService: chatService}
}

func (handler *ChatHandler) CreateChat(c echo.Context) error {
	var dto dtos.ChatCreateDto

	if err := c.Bind(&dto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(&dto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	chat, err := handler.chatService.CreateChat(dto.UserIds)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{"chat": chat})
}

func (handler *ChatHandler) GetAllChats(c echo.Context) error {
	userId := jwt.GetUserId(c)

	chats, err := handler.chatService.GetAllChats(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{"chats": chats})
}

func (handler *ChatHandler) GetChatById(c echo.Context) error {
	chatId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	chat, err := handler.chatService.GetChatById(uint(chatId))
	if err != nil {
		if errors.Is(err, ErrChatNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if errors.Is(err, ErrInternal) {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, echo.Map{"chat": chat})
}
