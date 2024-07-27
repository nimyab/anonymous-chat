package auth

import (
	"github.com/nimyab/anonymous-chat/internal/database/models"
	"github.com/nimyab/anonymous-chat/internal/handlers/auth/dtos"
	"golang.org/x/crypto/bcrypt"
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

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password))
	if err != nil {
		return nil, ErrWrongPassword
	}

	return &user, nil
}

func (s *AuthService) Registration(dto *dtos.UserRegistrationDto) (*models.User, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Login:    dto.Login,
		Password: string(hashPassword),
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
