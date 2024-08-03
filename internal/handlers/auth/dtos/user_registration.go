package dtos

type UserRegistrationDto struct {
	Name     string `json:"name" validate:"required,min=1,max=255"`
	Login    string `json:"login" validate:"required,min=1,max=255"`
	Password string `json:"password" validate:"required,min=1,max=255"`
}
