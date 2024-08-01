package dtos

type SendMessage struct {
	ChatID uint   `json:"chat_id" validate:"required" mapstructure:"chat_id"`
	Text   string `json:"text" validate:"required" mapstructure:"text"`
}
