package websocket

type Message struct {
	MessageName string                 `json:"message_name" validate:"required"`
	MessageBody map[string]interface{} `json:"message_body" validate:"required"`
}

type MessageWithSocketClient struct {
	Message      *Message
	SocketClient *SocketClient
}
