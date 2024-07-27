package auth

import "errors"

// errors for auth service
var (
	ErrUserNotFound  = errors.New("user not found")
	ErrWrongPassword = errors.New("wrong password")
)

// errors for auth handler
var (
	ErrValidation = errors.New("validation error")
)
