package websocket

import "github.com/nimyab/anonymous-chat/pkg/validators"

var hub *SocketHub

func GetSocketHub() *SocketHub {
	return hub
}

func init() {
	hub = &SocketHub{
		broadcast:       make(chan Message),
		register:        make(chan *SocketClientWithId),
		unregister:      make(chan *SocketClient),
		getClientById:   make(map[uint]*SocketClient),
		getIdByClient:   make(map[*SocketClient]uint),
		serverValidator: validators.NewServerValidator(),
	}
	go hub.Run()
}
