package handlers

import "fmt"

type SocketClientWithId struct {
	Client *SocketClient
	Id     string
}

type SocketHub struct {
	broadcast chan Message

	register   chan *SocketClientWithId
	unregister chan *SocketClient

	getClientById map[string]*SocketClient // получает сокет по id пользователя
	getIdByClient map[*SocketClient]string // получает id пользователя по сокету
}

// изза хаба не работе нужно переделать
var hub *SocketHub

func GetSocketHub() *SocketHub {
	if hub == nil {
		hub = &SocketHub{
			broadcast:     make(chan Message),
			register:      make(chan *SocketClientWithId),
			unregister:    make(chan *SocketClient),
			getClientById: make(map[string]*SocketClient),
			getIdByClient: make(map[*SocketClient]string),
		}
		go hub.Run()
	}
	return hub
}

func (h *SocketHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.getIdByClient[client.Client] = client.Id
			h.getClientById[client.Id] = client.Client

		case client := <-h.unregister:
			id := h.getIdByClient[client]
			delete(h.getIdByClient, client)
			delete(h.getClientById, id)
			close(client.messageChat)

		case message := <-h.broadcast:
			fmt.Println(message.Addressee, h.getClientById)
			client, ok := h.getClientById[message.Addressee]
			fmt.Println(client, ok)

			if ok {
				select {
				case client.messageChat <- message:
				default:
					id := h.getIdByClient[client]
					delete(h.getIdByClient, client)
					delete(h.getClientById, id)
					close(client.messageChat)
				}
			}
		}
	}
}
