package websocket

import (
	"github.com/mitchellh/mapstructure"
	"github.com/nimyab/anonymous-chat/internal/handlers/chat"
	"github.com/nimyab/anonymous-chat/internal/handlers/message"
	"github.com/nimyab/anonymous-chat/pkg/validators"
)

type WebsocketService struct {
	chatService    *chat.ChatService
	messageService *message.MessageService
	validator      *validators.ServerValidator
}

func (handler *WebsocketService) decodeAndValidate(input interface{}, messageBody interface{}) error {
	if err := mapstructure.Decode(input, &messageBody); err != nil {
		return ErrInvalidMessageBodyFormat
	}
	if err := handler.validator.Validate(&messageBody); err != nil {
		return err
	}
	return nil
}
