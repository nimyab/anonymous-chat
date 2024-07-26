package websocket

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

var hub *SocketHub

func init() {
	hub = &SocketHub{
		broadcast:     make(chan Message),
		register:      make(chan *SocketClientWithId),
		unregister:    make(chan *SocketClient),
		getClientById: make(map[string]*SocketClient),
		getIdByClient: make(map[*SocketClient]string),
	}
	go hub.Run()
}

func GetSocketHub() *SocketHub {
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
			client, ok := h.getClientById[message.Addressee]

			if ok {
				select {
				case client.messageChat <- message:
				default:
					h.unregister <- client
				}
			}
		}
	}
}
