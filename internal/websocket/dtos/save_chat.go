package dtos

type SaveChatRequest struct {
	ChatID string `json:"chat_id" validate:"required" mapstructure:"chat_id"`
}

type SaveChatResponse struct {
	ChatID   string `json:"chat_id" validate:"required" mapstructure:"chat_id"`
	Response bool   `json:"response" validate:"required" mapstructure:"response"`
}
