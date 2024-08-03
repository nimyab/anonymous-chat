package message

import (
	"github.com/nimyab/anonymous-chat/internal/websocket/dtos"
	"gorm.io/gorm"
	"sync"
)

type MessageService struct {
	gorm *gorm.DB
	mu   sync.Mutex
}

func NewMessageService(gorm *gorm.DB) *MessageService {
	return &MessageService{gorm: gorm}
}

func (ms *MessageService) CreateMessage(message *dtos.SendMessage) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	// models.Chat{}
}

func (ms *MessageService) GetAllMessageByChatId(chatId uint) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
}
