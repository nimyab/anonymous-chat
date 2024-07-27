package validators

import "github.com/go-playground/validator/v10"

type ApiValidator struct {
	validator *validator.Validate
}

func NewApiValidator() *ApiValidator {
	return &ApiValidator{validator: validator.New()}
}

func (v *ApiValidator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}
