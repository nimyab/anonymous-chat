package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nimyab/anonymous-chat/config"
	"github.com/nimyab/anonymous-chat/internal/database"
	"github.com/nimyab/anonymous-chat/internal/handlers/auth"
	"github.com/nimyab/anonymous-chat/internal/websocket"
)

const PORT = ":9999"

func main() {
	e := echo.New()
	cfg := config.GetEnvConfig()
	db := database.ConnectAndMigrateDatabase(cfg)

	auth.NewAuthService(db)

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339}\t${method}\t${uri}\tstatus ${status}\tuser agent ${user_agent}\n",
	}))
	e.Use(middleware.Recover())

	e.POST("/login", func(c echo.Context) error {
		return nil
	})
	e.POST("/logout", func(c echo.Context) error {
		return nil
	})

	e.GET("/ws", websocket.SocketConn)

	e.Logger.Infof("Server start on %s port", cfg.PORT)
	e.Logger.Fatal(e.Start(cfg.PORT))
}
