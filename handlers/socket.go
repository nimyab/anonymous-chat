package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log/slog"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func SocketConn(w http.ResponseWriter, r *http.Request) {
	id := r.Header.Get("Id")
	if id == "" {
		http.Error(w, "Not found user id", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Error upgrading to websockets", http.StatusBadRequest)
		return
	}

	hub := GetSocketHub()
	client := &SocketClient{
		conn:        conn,
		hub:         hub,
		messageChat: make(chan Message),
	}
	slog.Info(fmt.Sprintf("New connection %s", id))

	hub.register <- &SocketClientWithId{
		Client: client,
		Id:     id,
	}

	go client.ReadPump()
	go client.WritePump()
}

type Message struct {
	MessageName string      `json:"message-name"`
	MessageBody interface{} `json:"message-body"`
	Addressee   string      `json:"addressee"`
}

type SocketClient struct {
	conn *websocket.Conn

	hub *SocketHub

	messageChat chan Message
}

func (sc *SocketClient) ReadPump() {
	defer sc.conn.Close()

	fmt.Println("Reading pump")

	for {
		_, message, err := sc.conn.ReadMessage()
		fmt.Println(message)
		if err != nil {
			slog.Error(err.Error())
			break
		}

		var mess Message
		if err := json.Unmarshal(message, &mess); err != nil {
			slog.Error(err.Error())
			break
		}

		sc.hub.broadcast <- mess
	}
}

func (sc *SocketClient) WritePump() {
	defer sc.conn.Close()

	fmt.Println("Writing pump")

	for {
		select {
		case message, ok := <-sc.messageChat:
			fmt.Println(message)
			if !ok {
				sc.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			bytes, err := json.Marshal(message)
			if err != nil {
				sc.conn.WriteMessage(websocket.TextMessage, []byte("Marshalling error message"))
			} else {
				sc.conn.WriteMessage(websocket.TextMessage, bytes)
			}
		}
	}
}
