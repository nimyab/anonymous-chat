package validators

import "github.com/go-playground/validator/v10"

type ServerValidator struct {
	validator *validator.Validate
}

func NewServerValidator() *ServerValidator {
	return &ServerValidator{validator: validator.New()}
}

func (v *ServerValidator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}
