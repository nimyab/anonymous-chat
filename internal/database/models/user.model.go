package models

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Name     string `gorm:"size:255;not null" json:"name"`
	Login    string `gorm:"size:255;not null;unique" json:"login"`
	Password string `gorm:"size:255;not null" json:"-"`
	IsAdmin  bool   `gorm:"default:false" json:"is_admin"`

	Chats    []Chat    `gorm:"many2many:user_chats" json:"chats"`
	Messages []Message `gorm:"foreignKey:ID;references:ID" json:"messages"`
}

func (u *User) GeneratePassword() error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashPassword)
	return nil
}

func (u *User) ComparePasswords(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err
}
