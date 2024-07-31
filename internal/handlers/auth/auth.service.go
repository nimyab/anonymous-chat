package auth

import (
	"github.com/nimyab/anonymous-chat/internal/database/models"
	"github.com/nimyab/anonymous-chat/internal/handlers/auth/dtos"
	"github.com/nimyab/anonymous-chat/internal/jwt"
	"gorm.io/gorm"
	"log/slog"
	"sync"
)

type AuthService struct {
	gorm *gorm.DB
	mu   sync.Mutex
}

func NewAuthService(gorm *gorm.DB) *AuthService {
	return &AuthService{gorm: gorm}
}

func (s *AuthService) Login(dto *dtos.UserLoginDto) (u *models.User, accessToken string, refreshToken string, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user := models.User{
		Login: dto.Login,
	}
	result := s.gorm.Preload("Chats").Preload("Messages").First(&user)
	if result.RowsAffected == 0 {
		return nil, "", "", ErrUserNotFound
	}

	if err := user.ComparePasswords(dto.Password); err != nil {
		return nil, "", "", ErrWrongPassword
	}

	accessToken, refreshToken, err = jwt.GenerateTokens(user.ID)
	if err != nil {
		slog.Error(err.Error())
		return nil, "", "", ErrInternal
	}

	return &user, accessToken, refreshToken, nil
}

func (s *AuthService) Registration(dto *dtos.UserRegistrationDto) (*models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user := models.User{
		Login:    dto.Login,
		Password: dto.Password,
		Name:     dto.Name,
	}

	if err := user.GeneratePassword(); err != nil {
		return nil, err
	}

	result := s.gorm.Create(&user)

	if result.Error != nil {
		if result.RowsAffected == 0 {
			return nil, result.Error
		}
		return &user, result.Error
	}
	return &user, nil
}

func (s *AuthService) Refresh(refreshToken string) (newAccessToken string, newRefreshToken string, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	userId, err := jwt.ParseToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	newAccessToken, newRefreshToken, err = jwt.GenerateTokens(userId)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

func (s *AuthService) UserInfo(userId uint) (*models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user := models.User{
		ID: userId,
	}

	result := s.gorm.Preload("Chats").Preload("Messages").First(&user)
	if result.RowsAffected == 0 {
		return nil, ErrUserNotFound
	}

	return &user, nil
}
