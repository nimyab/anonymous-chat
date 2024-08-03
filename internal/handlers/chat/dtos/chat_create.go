package dtos

type ChatCreateDto struct {
	UserIds []uint `json:"user_ids" validate:"required"`
}
