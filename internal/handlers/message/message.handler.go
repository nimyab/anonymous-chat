package message

import "github.com/labstack/echo/v4"

type ChatHandler struct {
	messageService *MessageService
}

func NewChatHandler(messageService *MessageService) *ChatHandler {
	return &ChatHandler{messageService: messageService}
}

func (chatHandler *ChatHandler) HandleMessage(c echo.Context) error {
	return nil
}
