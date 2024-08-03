package models

import "time"

type Chat struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Users    []User    `gorm:"many2many:user_chats" json:"users"`
	Messages []Message `gorm:"foreignKey:ID;references:ID" json:"messages"`
}
