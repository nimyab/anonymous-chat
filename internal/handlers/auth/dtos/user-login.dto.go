package dtos

type UserLoginDto struct {
	Login    string `json:"login" validate:"required,min=1,max=255"`
	Password string `json:"password" validate:"required,min=1,max=255"`
}
