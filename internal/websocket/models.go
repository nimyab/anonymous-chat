package websocket

import "github.com/gorilla/websocket"

type Message struct {
	MessageName string      `json:"message_name" validate:"required"`
	MessageBody interface{} `json:"message_body" validate:"required"`
}

type SocketClient struct {
	conn *websocket.Conn

	hub *SocketHub

	messageChat chan Message
}
