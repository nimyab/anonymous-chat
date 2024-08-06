package dtos

type DeleteChat struct {
	ChatID uint `json:"chat_id" validate:"required" mapstructure:"chat_id"`
}
