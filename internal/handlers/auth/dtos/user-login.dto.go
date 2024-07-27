package dtos

type UserLoginDto struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (u *UserLoginDto) Validation() error {
	return nil
}
