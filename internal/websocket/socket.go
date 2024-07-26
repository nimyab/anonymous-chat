package websocket

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

const (
	writeWait = 10 * time.Second

	pongWait = 120 * time.Second

	pingPeriod = (pongWait * 8) / 10

	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func SocketConn(c echo.Context) error {

	id := c.Request().Header.Get("Id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, struct{ message string }{message: "Not found ID"})
	}

	conn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		return c.JSON(http.StatusBadRequest, struct{ message string }{message: "Error upgrading to websockets"})
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

	return nil
}

type Message struct {
	MessageName string      `json:"message_name"`
	MessageBody interface{} `json:"message_body"`
	Addressee   string      `json:"addressee"`
}

type SocketClient struct {
	conn *websocket.Conn

	hub *SocketHub

	messageChat chan Message
}

func (sc *SocketClient) ReadPump() {
	defer func() {
		sc.hub.unregister <- sc
		sc.conn.Close()
	}()

	sc.conn.SetReadLimit(maxMessageSize)
	sc.conn.SetReadDeadline(time.Now().Add(pongWait))
	sc.conn.SetPongHandler(func(appData string) error {
		sc.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := sc.conn.ReadMessage()
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
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		sc.conn.Close()
	}()

	for {
		select {
		case message, ok := <-sc.messageChat:
			if !ok {
				sc.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			bytes, err := json.Marshal(message)
			sc.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				sc.conn.WriteMessage(websocket.TextMessage, []byte("Marshalling error message"))
			} else {
				sc.conn.WriteMessage(websocket.TextMessage, bytes)
			}
		case <-ticker.C:
			sc.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := sc.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
