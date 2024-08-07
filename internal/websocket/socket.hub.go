package websocket

import (
	"github.com/mitchellh/mapstructure"
	messageDtos "github.com/nimyab/anonymous-chat/internal/handlers/message/dtos"
	"github.com/nimyab/anonymous-chat/internal/websocket/dtos"
	"github.com/nimyab/anonymous-chat/pkg/validators"
)

// message names

// here, there
// search_interlocutor, found_interlocutor
// send_message, send_message
// delete_chat, delete_chat
// stop_interlocutor, stop_interlocutor

const (
	searchInterlocutor = "search_interlocutor"
	foundInterlocutor  = "found_interlocutor"
	sendMessage        = "send_message"
	deleteChat         = "delete_chat"
	searchStop         = "search_stop"
)

type SocketHub struct {
	broadcast chan *MessageWithSocketClient

	register    chan *SocketClientWithId
	unregister  chan *SocketClient
	searchUsers <-chan []uint

	getClientById map[uint]*SocketClient
	getIdByClient map[*SocketClient]uint

	validator *validators.ServerValidator

	websocketService *WebsocketService
}

func (h *SocketHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)
		case client := <-h.unregister:
			h.unregisterClient(client)
		case userIds := <-h.searchUsers:
			h.handleFoundInterlocutor(userIds)
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
	h.websocketService.userInSearch.DeleteUserId(id)
	delete(h.getIdByClient, client)
	delete(h.getClientById, id)
	close(client.messageChan)
}

func (h *SocketHub) handleMessage(messageWithSocket *MessageWithSocketClient) {
	switch messageWithSocket.Message.MessageName {
	case searchInterlocutor:
		h.handleSearchInterlocutor(messageWithSocket)
	case sendMessage:
		h.handleSendMessage(messageWithSocket)
	case deleteChat:
		h.deleteChat(messageWithSocket)
	case searchStop:
		h.handleStopSearch(messageWithSocket)
	default:
		messageWithSocket.SocketClient.SendError(ErrSuchMessageNameNoExist)
	}
}

func (h *SocketHub) handleFoundInterlocutor(userIds []uint) {
	chat, err := h.websocketService.chatService.CreateChat(userIds)
	if err != nil {
		h.websocketService.userInSearch.Push(userIds[0])
		h.websocketService.userInSearch.Push(userIds[1])
		return
	}
	mess := &Message{MessageName: foundInterlocutor, MessageBody: map[string]interface{}{
		"chat": chat,
	}}
	for _, userId := range userIds {
		client, ok := h.getClientById[userId]
		if ok {
			client.messageChan <- mess
		}
	}
}

func (h *SocketHub) handleSearchInterlocutor(messageWithSocket *MessageWithSocketClient) {
	userId, ok := h.getIdByClient[messageWithSocket.SocketClient]
	if !ok {
		return
	}
	h.websocketService.userInSearch.Push(userId)
	messageWithSocket.SocketClient.messageChan <- &Message{
		searchInterlocutor,
		map[string]interface{}{
			"message": "start",
		},
	}
}

func (h *SocketHub) handleStopSearch(messageWithSocket *MessageWithSocketClient) {
	userId, ok := h.getIdByClient[messageWithSocket.SocketClient]
	if !ok {
		return
	}
	h.websocketService.userInSearch.DeleteUserId(userId)
	messageWithSocket.SocketClient.messageChan <- &Message{
		searchStop,
		map[string]interface{}{
			"message": "stop",
		},
	}
}

func (h *SocketHub) handleSendMessage(messageWithSocket *MessageWithSocketClient) {
	var messageBody dtos.SendMessage
	if err := mapstructure.Decode(messageWithSocket.Message.MessageBody, &messageBody); err != nil {
		messageWithSocket.SocketClient.SendError(ErrInvalidMessageBodyFormat)
		return
	}
	if err := h.validator.Validate(&messageBody); err != nil {
		messageWithSocket.SocketClient.SendError(err)
		return
	}

	userId, ok := h.getIdByClient[messageWithSocket.SocketClient]
	if !ok {
		return
	}
	messageDto := &messageDtos.MessageCreateDto{
		Text:   messageBody.Text,
		ChatId: messageBody.ChatID,
		UserId: userId,
	}
	message, userIds, err := h.websocketService.messageService.CreateMessage(messageDto)
	if err != nil {
		messageWithSocket.SocketClient.SendError(err)
	}

	mess := &Message{MessageName: sendMessage, MessageBody: map[string]interface{}{
		"message": message,
	}}
	for _, userId := range userIds {
		client, ok := h.getClientById[userId]
		if ok {
			client.messageChan <- mess
		}
	}
}
func (h *SocketHub) deleteChat(messageWithSocket *MessageWithSocketClient) {
	var messageBody dtos.DeleteChat
	if err := mapstructure.Decode(messageWithSocket.Message.MessageBody, &messageBody); err != nil {
		messageWithSocket.SocketClient.SendError(ErrInvalidMessageBodyFormat)
		return
	}
	if err := h.validator.Validate(&messageBody); err != nil {
		messageWithSocket.SocketClient.SendError(err)
		return
	}

	userIds, err := h.websocketService.chatService.DeleteChat(messageBody.ChatID)
	if err != nil {
		messageWithSocket.SocketClient.SendError(err)
		return
	}

	mess := &Message{MessageName: deleteChat, MessageBody: map[string]interface{}{
		"chat_id": messageBody.ChatID,
	}}

	for _, userId := range userIds {
		client, ok := h.getClientById[userId]
		if ok {
			client.messageChan <- mess
		}
	}
}
