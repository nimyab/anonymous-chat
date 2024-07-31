package websocket

import (
	"github.com/nimyab/anonymous-chat/pkg/validators"
)

type SocketClientWithId struct {
	Client *SocketClient
	ID     uint
}

type SocketHub struct {
	broadcast chan Message

	register   chan *SocketClientWithId
	unregister chan *SocketClient

	getClientById map[uint]*SocketClient // получает сокет по id пользователя
	getIdByClient map[*SocketClient]uint // получает id пользователя по сокету

	serverValidator *validators.ServerValidator
}

func (h *SocketHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.getIdByClient[client.Client] = client.ID
			h.getClientById[client.ID] = client.Client

		case client := <-h.unregister:
			id := h.getIdByClient[client]
			delete(h.getIdByClient, client)
			delete(h.getClientById, id)
			close(client.messageChat)

		case message := <-h.broadcast:
			switch message.MessageName {
			case "search_interlocutor":
				// todo
			case "send_message":
				// todo
			case "save_chat":
				// todo
			}
		}
	}
}
