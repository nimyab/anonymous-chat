package dtos

type UserRegistrationDto struct {
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
}
