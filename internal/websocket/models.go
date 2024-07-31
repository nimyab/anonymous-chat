package websocket

type Message struct {
	MessageName string      `json:"message_name" validate:"required"`
	MessageBody interface{} `json:"message_body" validate:"required"`
}

type MessageWithSocketClient struct {
	Message      *Message
	SocketClient *SocketClient
}
