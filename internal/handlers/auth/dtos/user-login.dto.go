package dtos

type UserLoginDto struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
