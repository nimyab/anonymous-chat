package websocket

import (
	"github.com/mitchellh/mapstructure"
	"github.com/nimyab/anonymous-chat/internal/websocket/dtos"
	"github.com/nimyab/anonymous-chat/pkg/validators"
)

type SocketHub struct {
	broadcast chan *MessageWithSocketClient

	register   chan *SocketClientWithId
	unregister chan *SocketClient

	getClientById map[uint]*SocketClient
	getIdByClient map[*SocketClient]uint

	serverValidator *validators.ServerValidator
}

func (h *SocketHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)
		case client := <-h.unregister:
			h.unregisterClient(client)
		case messageWithSocket := <-h.broadcast:
			h.handleMessage(messageWithSocket)
		}
	}
}

func (h *SocketHub) registerClient(client *SocketClientWithId) {
	h.getIdByClient[client.Client] = client.ID
	h.getClientById[client.ID] = client.Client
}

func (h *SocketHub) unregisterClient(client *SocketClient) {
	id := h.getIdByClient[client]
	delete(h.getIdByClient, client)
	delete(h.getClientById, id)
	close(client.messageChat)
}

func (h *SocketHub) handleMessage(messageWithSocket *MessageWithSocketClient) {
	switch messageWithSocket.Message.MessageName {
	case "search_interlocutor":
		h.handleSearchInterlocutor(messageWithSocket)
	case "send_message":
		h.handleSendMessage(messageWithSocket)
	case "request_save_chat":
		h.handleRequestSaveChat(messageWithSocket)
	case "accept_save_chat":
		h.handleAcceptSaveChat(messageWithSocket)
	case "reject_save_chat":
		h.handleRejectSaveChat(messageWithSocket)
	default:
		messageWithSocket.SocketClient.SendError(ErrSuchMessageNameNoExist)
	}
}

func (h *SocketHub) handleSearchInterlocutor(messageWithSocket *MessageWithSocketClient) {
	// todo
}

func (h *SocketHub) handleSendMessage(messageWithSocket *MessageWithSocketClient) {
	var messageBody dtos.SendMessage
	if err := mapstructure.Decode(messageWithSocket.Message.MessageBody, &messageBody); err != nil {
		messageWithSocket.SocketClient.SendError(ErrInvalidMessageBodyFormat)
		return
	}
	if err := h.serverValidator.Validate(&messageBody); err != nil {
		messageWithSocket.SocketClient.SendError(err)
		return
	}
	// todo
}

func (h *SocketHub) handleRequestSaveChat(messageWithSocket *MessageWithSocketClient) {
	var messageBody dtos.SaveChatRequest
	if err := mapstructure.Decode(messageWithSocket.Message.MessageBody, &messageBody); err != nil {
		messageWithSocket.SocketClient.SendError(ErrInvalidMessageBodyFormat)
		return
	}
	if err := h.serverValidator.Validate(messageBody); err != nil {
		messageWithSocket.SocketClient.SendError(err)
		return
	}
	// todo
}

func (h *SocketHub) handleAcceptSaveChat(messageWithSocket *MessageWithSocketClient) {
	var messageBody dtos.SaveChatResponse
	if err := mapstructure.Decode(messageWithSocket.Message.MessageBody, &messageBody); err != nil {
		messageWithSocket.SocketClient.SendError(ErrInvalidMessageBodyFormat)
		return
	}
	if err := h.serverValidator.Validate(messageBody); err != nil {
		messageWithSocket.SocketClient.SendError(err)
		return
	}
	// todo
}

func (h *SocketHub) handleRejectSaveChat(messageWithSocket *MessageWithSocketClient) {
	var messageBody dtos.SaveChatResponse
	if err := mapstructure.Decode(messageWithSocket.Message.MessageBody, &messageBody); err != nil {
		messageWithSocket.SocketClient.SendError(ErrInvalidMessageBodyFormat)
		return
	}
	if err := h.serverValidator.Validate(messageBody); err != nil {
		messageWithSocket.SocketClient.SendError(err)
		return
	}
	// todo
}
