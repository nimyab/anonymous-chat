package models

import "time"

type Chat struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Messages []Message `gorm:"foreignKey:ChatID;references:ID" json:"messages"`
}
