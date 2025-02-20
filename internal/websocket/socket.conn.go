package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/nimyab/anonymous-chat/internal/jwt"
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

type SocketClient struct {
	conn        *websocket.Conn
	hub         *SocketHub
	messageChan chan *Message
}

type SocketClientWithId struct {
	Client *SocketClient
	ID     uint
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func SocketConn(c echo.Context) error {
	userId := jwt.GetUserId(c)

	conn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		slog.Error(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, ErrUpgradingToWebsocket)
	}

	hub := GetSocketHub()
	client := &SocketClient{
		conn:        conn,
		hub:         hub,
		messageChan: make(chan *Message),
	}

	slog.Info(fmt.Sprintf("New connection %d", userId))

	hub.register <- &SocketClientWithId{
		Client: client,
		ID:     userId,
	}

	go client.ReadPump()
	go client.WritePump()

	return nil
}

func (sc *SocketClient) ReadPump() {
	defer func() {
		sc.hub.unregister <- sc
		if err := sc.conn.Close(); err != nil {
			slog.Error(err.Error())
		}
	}()

	sc.conn.SetReadLimit(maxMessageSize)
	if err := sc.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		slog.Error(err.Error())
	}
	sc.conn.SetPongHandler(func(appData string) error {
		return sc.conn.SetReadDeadline(time.Now().Add(pongWait))
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
			sc.SendError(err)
			continue
		}

		if err := sc.hub.validator.Validate(&mess); err != nil {
			slog.Error(err.Error())
			sc.SendError(err)
			continue
		}

		sc.hub.broadcast <- &MessageWithSocketClient{Message: &mess, SocketClient: sc}
	}
}

func (sc *SocketClient) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case message, ok := <-sc.messageChan:
			if !ok {
				sc.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			bytes, err := json.Marshal(message)
			if err := sc.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				slog.Error(err.Error())
			}

			if err != nil {
				sc.conn.WriteMessage(websocket.TextMessage, []byte("Marshalling error message"))
			} else {
				sc.conn.WriteMessage(websocket.TextMessage, bytes)
			}
		case <-ticker.C:
			if err := sc.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				slog.Error(err.Error())
			}
			if err := sc.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (sc *SocketClient) SendError(err error) {
	sc.messageChan <- &Message{
		MessageName: "error",
		MessageBody: map[string]interface{}{
			"message": err.Error(),
		},
	}
}
