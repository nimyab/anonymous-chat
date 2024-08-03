package websocket

import (
	"github.com/nimyab/anonymous-chat/internal/handlers/chat"
	"github.com/nimyab/anonymous-chat/internal/handlers/message"
)

type WebsocketService struct {
	chatService    *chat.ChatService
	messageService *message.MessageService
}
