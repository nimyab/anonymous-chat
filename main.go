package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/nimyab/anonymous-chat/handlers"
)

const PORT = ":9999"

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/login", handlers.Login)
	mux.HandleFunc("/logout", handlers.Logout)

	mux.HandleFunc("/ws", handlers.SocketConn)

	slog.Info(fmt.Sprintf("Listening on %s", PORT))
	if err := http.ListenAndServe(PORT, mux); err != nil {
		slog.Error(err.Error())
	}
}
