package message

import (
	"github.com/nimyab/anonymous-chat/internal/database/models"
	"github.com/nimyab/anonymous-chat/internal/handlers/message/dtos"
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

func (ms *MessageService) CreateMessage(dto *dtos.MessageCreateDto) (*models.Message, []uint, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	message := models.Message{
		Text:   dto.Text,
		ChatID: dto.ChatId,
		UserID: dto.UserId,
	}
	result := ms.gorm.Create(&message)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	chat := models.Chat{
		ID: message.ChatID,
	}
	result = ms.gorm.Preload("Users").First(&chat)
	if result.Error != nil {
		ms.gorm.Delete(&message)
		return nil, nil, result.Error
	}

	var userIds []uint
	for _, user := range chat.Users {
		userIds = append(userIds, user.ID)
	}

	return &message, userIds, nil
}

func (ms *MessageService) GetAllMessageByChatId(chatId uint) ([]*models.Message, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	var messages []*models.Message
	result := ms.gorm.
		Preload("Chats").
		Joins("JOIN chats ON messages.chat_id = chats.id").
		Where("chats.id = ?", chatId).
		Find(&messages)
	if result.Error != nil {
		return nil, result.Error
	}
	return messages, nil
}
