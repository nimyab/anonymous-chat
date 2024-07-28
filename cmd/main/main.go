package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nimyab/anonymous-chat/internal/config"
	"github.com/nimyab/anonymous-chat/internal/database"
	"github.com/nimyab/anonymous-chat/internal/handlers/auth"
	"github.com/nimyab/anonymous-chat/internal/websocket"
	"github.com/nimyab/anonymous-chat/pkg/validators"
)

const PORT = ":9999"

func main() {
	e := echo.New()
	cfg := config.GetEnvConfig()
	db := database.ConnectAndMigrateDatabase(cfg)

	// services
	authService := auth.NewAuthService(db)

	// handlers
	authHandler := auth.NewAuthHandler(authService)

	e.Validator = validators.NewApiValidator()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339}\t${method}\t${uri}\tstatus ${status}\tuser agent ${user_agent}\n",
	}))
	e.Use(middleware.Recover())
	api := e.Group("/api")

	// auth routes
	api.POST("/login", authHandler.Login)
	api.POST("/registration", authHandler.Registration)
	api.POST("/logout", authHandler.Logout)

	// socket routes
	api.GET("/ws", websocket.SocketConn)

	e.Logger.Infof("Server start on %s port", cfg.PORT)
	e.Logger.Fatal(e.Start(cfg.PORT))
}
