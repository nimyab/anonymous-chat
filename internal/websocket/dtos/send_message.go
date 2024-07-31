package dtos

type SendMessage struct {
	ChatID uint   `json:"chat_id" validate:"required"`
	Text   string `json:"text" validate:"required"`
}
