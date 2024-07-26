package models

import "time"

type Message struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Text string `gorm:"not null" json:"text"`

	UserID uint `json:"user_id"`
	ChatID uint `json:"chat_id"`
}
