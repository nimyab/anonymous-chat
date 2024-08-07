package chat

import (
	"errors"
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

	subQuery := s.gorm.Table("user_chats").Select("chat_id").
		Where("user_id = ?", userIds[0]).
		Or("user_id = ?", userIds[1]).
		Group("chat_id").
		Having("COUNT(DISTINCT user_id) = 2")

	chat := models.Chat{}
	err := s.gorm.Preload("Users").Preload("Messages").Where("id IN (?)", subQuery).First(&chat).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if err == nil {
		return &chat, nil
	}

	chat = models.Chat{
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

func (s *ChatService) DeleteChat(chatId uint) (userIds []uint, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	chat := &models.Chat{
		ID: chatId,
	}

	if err := s.gorm.Preload("Users").First(&chat).Error; err != nil {
		slog.Error(err.Error())
		return nil, ErrInternal
	}
	for _, user := range chat.Users {
		userIds = append(userIds, user.ID)
	}

	if err := s.gorm.Model(&chat).Association("Users").Delete(chat.Users); err != nil {
		slog.Error(err.Error())
		return nil, ErrInternal
	}

	if err := s.gorm.Where("chat_id = ?", chatId).Delete(&models.Message{}).Error; err != nil {
		slog.Error(err.Error())
		return nil, ErrInternal
	}

	if err := s.gorm.Delete(&chat).Error; err != nil {
		slog.Error(err.Error())
		return nil, ErrInternal
	}

	return userIds, nil
}
