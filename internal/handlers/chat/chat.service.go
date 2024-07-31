package chat

import (
	"github.com/nimyab/anonymous-chat/internal/database/models"
	"gorm.io/gorm"
	"log/slog"
	"sync"
)

type ChatService struct {
	gorm *gorm.DB
	mu   sync.Mutex
}

func NewChatService(gorm *gorm.DB) *ChatService {
	return &ChatService{gorm: gorm}
}

func (s *ChatService) CreateChat(userIds []uint) (*models.Chat, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var users []models.User
	result := s.gorm.Find(&users, userIds)

	if result.Error != nil {
		return nil, result.Error
	}
	if len(users) < 2 {
		return nil, ErrNotFoundTwoUsers
	}

	chat := models.Chat{
		Users: users,
	}
	result = s.gorm.Create(&chat)
	if result.Error != nil {
		return nil, result.Error
	}

	return &chat, nil
}

func (s *ChatService) GetAllChats(userId uint) ([]*models.Chat, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var chats []*models.Chat
	result := s.gorm.
		Preload("Users").
		Preload("Messages").
		Joins("JOIN user_chats ON chats.id = user_chats.chat_id").
		Where("user_chats.user_id = ?", userId).
		Find(&chats)

	if result.Error != nil {
		slog.Error(result.Error.Error())
		return nil, ErrInternal
	}

	return chats, nil
}

func (s *ChatService) GetChatById(chatId uint) (*models.Chat, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	chat := &models.Chat{
		ID: chatId,
	}

	result := s.gorm.Joins("users").First(chat)
	if result.Error != nil {
		slog.Error(result.Error.Error())
		return nil, ErrInternal
	}

	if result.RowsAffected == 0 {
		return nil, ErrChatNotFound
	}

	return chat, nil
}
