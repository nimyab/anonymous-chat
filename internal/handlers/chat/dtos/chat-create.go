package dtos

type ChatCreate struct {
	UserIds []uint `json:"user_ids" validate:"required"`
}
