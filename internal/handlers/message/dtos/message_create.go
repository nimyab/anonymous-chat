package dtos

type MessageCreateDto struct {
	Text   string `json:"text" validate:"required"`
	ChatId uint   `json:"chat_id" validate:"required"`
	UserId uint   `json:"user_id" validate:"required"`
}
