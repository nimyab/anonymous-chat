package websocket

import "github.com/gorilla/websocket"

type SocketClient struct {
	conn *websocket.Conn

	hub *SocketHub

	messageChat chan *Message
}

type SocketClientWithId struct {
	Client *SocketClient
	ID     uint
}

type Message struct {
	MessageName string      `json:"message_name" validate:"required"`
	MessageBody interface{} `json:"message_body" validate:"required"`
}

type MessageWithSocketClient struct {
	Message      *Message
	SocketClient *SocketClient
}
