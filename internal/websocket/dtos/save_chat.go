package dtos

type SaveChatRequest struct {
	ChatID string `json:"chat_id" validate:"required"`
}

type SaveChatResponse struct {
	ChatID   string `json:"chat_id" validate:"required"`
	Response bool   `json:"response" validate:"required"`
}
