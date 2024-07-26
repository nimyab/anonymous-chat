package auth

import (
	"github.com/nimyab/anonymous-chat/internal/database/models"
	"github.com/nimyab/anonymous-chat/internal/handlers/auth/dtos"
	"gorm.io/gorm"
)

type AuthService struct {
	gorm *gorm.DB
}

func NewAuthService(gorm *gorm.DB) *AuthService {
	return &AuthService{gorm: gorm}
}

func (s *AuthService) Login(dto dtos.UserLoginDto) (*models.User, error) {
	user := models.User{
		Login: dto.Login,
	}
	result := s.gorm.Find(&user)
	if result.RowsAffected == 0 {
		return nil, ErrUserNotFound
	}

	return &user, nil
}

func (s *AuthService) Registration(dto *dtos.UserRegistrationDto) (*models.User, error) {
	// TODO: add hashing password

	user := models.User{
		Login:    dto.Login,
		Password: dto.Password,
		Name:     dto.Name,
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

func (s *AuthService) Logout() error { return nil }
