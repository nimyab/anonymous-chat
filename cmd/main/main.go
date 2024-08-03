package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nimyab/anonymous-chat/internal/config"
	"github.com/nimyab/anonymous-chat/internal/database"
	"github.com/nimyab/anonymous-chat/internal/handlers/auth"
	"github.com/nimyab/anonymous-chat/internal/handlers/chat"
	"github.com/nimyab/anonymous-chat/internal/handlers/message"
	"github.com/nimyab/anonymous-chat/internal/jwt"
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
	chatService := chat.NewChatService(db)
	messageService := message.NewMessageService(db)

	// handlers
	authHandler := auth.NewAuthHandler(authService)
	chatHandler := chat.NewChatHandler(chatService)
	messageHandler := message.NewChatHandler(messageService)

	e.Validator = validators.NewServerValidator()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339}\t${method}\t${uri}\tstatus ${status}\tuser agent ${user_agent}\n",
	}))
	e.Use(middleware.Recover())
	api := e.Group("/api")

	// auth routes
	authRoutes := api.Group("/auth")
	authRoutes.POST("/login", authHandler.Login)
	authRoutes.POST("/registration", authHandler.Registration)
	authRoutes.POST("/logout", authHandler.Logout)
	authRoutes.GET("/refresh", authHandler.Refresh)
	authRoutes.GET("/user-info", authHandler.UserInfo, jwt.Middleware())

	// chat routes
	chatRoutes := api.Group("/chat", jwt.Middleware())
	chatRoutes.GET("", chatHandler.GetAllChats)
	chatRoutes.GET("/:id", chatHandler.GetChatById)
	//chatRoutes.POST("", chatHandler.CreateChat)

	// message routes
	messageRoutes := api.Group("/message", jwt.Middleware())
	messageRoutes.GET("", messageHandler.HandleMessage)

	// socket routes
	websocket.StartSocketHub(chatService, messageService)
	api.GET("/ws", websocket.SocketConn, jwt.Middleware())

	e.Logger.Infof("Server start on %s port", cfg.Port)
	e.Logger.Fatal(e.Start(cfg.Port))
}
