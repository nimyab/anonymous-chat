package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Name     string `gorm:"size:255;not null" json:"name"`
	Login    string `gorm:"size:255;not null;unique" json:"login"`
	Password string `gorm:"size:255;not null" json:"-"`

	Chats    []Chat    `gorm:"many2many:user_chats" json:"chats"`
	Messages []Message `gorm:"foreignKey:UserID;references:ID" json:"messages"`
}
