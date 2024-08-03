package websocket

import (
	"github.com/nimyab/anonymous-chat/internal/handlers/chat"
	"github.com/nimyab/anonymous-chat/internal/handlers/message"
	"github.com/nimyab/anonymous-chat/pkg/validators"
	"log/slog"
)

var hub *SocketHub

func GetSocketHub() *SocketHub {
	return hub
}

func StartSocketHub(chatService *chat.ChatService, messageService *message.MessageService) {
	hub = &SocketHub{
		broadcast:     make(chan *MessageWithSocketClient),
		register:      make(chan *SocketClientWithId),
		unregister:    make(chan *SocketClient),
		getClientById: make(map[uint]*SocketClient),
		getIdByClient: make(map[*SocketClient]uint),
		websocketHandler: &WebsocketService{
			chatService:    chatService,
			messageService: messageService,
			validator:      validators.NewServerValidator(),
		},
	}
	go hub.Run()
	slog.Info("socket hub start")
}
