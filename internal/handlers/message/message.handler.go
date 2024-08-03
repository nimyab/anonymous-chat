package message

import (
	"github.com/labstack/echo/v4"
	"github.com/nimyab/anonymous-chat/internal/handlers/message/dtos"
	"net/http"
	"strconv"
)

type ChatHandler struct {
	messageService *MessageService
}

func NewChatHandler(messageService *MessageService) *ChatHandler {
	return &ChatHandler{messageService: messageService}
}

func (handler *ChatHandler) CreateMessage(c echo.Context) error {
	var dto dtos.MessageCreateDto
	if err := c.Bind(&dto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(&dto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	message, userIds, err := handler.messageService.CreateMessage(&dto)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message":  message,
		"user_ids": userIds,
	})
}

func (handler *ChatHandler) GetAllMessageByChatId(c echo.Context) error {
	chatId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	messages, err := handler.messageService.GetAllMessageByChatId(uint(chatId))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"messages": messages,
	})
}
