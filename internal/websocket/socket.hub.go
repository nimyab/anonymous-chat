package websocket

import (
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
	}
}

func (h *SocketHub) handleSearchInterlocutor(messageWithSocket *MessageWithSocketClient) {
	// todo
}

func (h *SocketHub) handleSendMessage(messageWithSocket *MessageWithSocketClient) {
	data := messageWithSocket.Message.MessageBody.(*dtos.SendMessage)
	if err := h.serverValidator.Validate(data); err != nil {
		messageWithSocket.SocketClient.SendError(err)
		return
	}
	// todo
}

func (h *SocketHub) handleRequestSaveChat(messageWithSocket *MessageWithSocketClient) {
	data := messageWithSocket.Message.MessageBody.(*dtos.SaveChatRequest)
	if err := h.serverValidator.Validate(data); err != nil {
		messageWithSocket.SocketClient.SendError(err)
		return
	}
	// todo
}

func (h *SocketHub) handleAcceptSaveChat(messageWithSocket *MessageWithSocketClient) {
	data := messageWithSocket.Message.MessageBody.(*dtos.SaveChatResponse)
	if err := h.serverValidator.Validate(data); err != nil {
		messageWithSocket.SocketClient.SendError(err)
		return
	}
	// todo
}

func (h *SocketHub) handleRejectSaveChat(messageWithSocket *MessageWithSocketClient) {
	data := messageWithSocket.Message.MessageBody.(*dtos.SaveChatResponse)
	if err := h.serverValidator.Validate(data); err != nil {
		messageWithSocket.SocketClient.SendError(err)
		return
	}
	// todo
}
