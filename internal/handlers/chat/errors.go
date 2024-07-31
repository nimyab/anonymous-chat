package chat

import "errors"

var (
	ErrChatNotFound     = errors.New("chat not found")
	ErrInternal         = errors.New("internal error")
	ErrNotFoundTwoUsers = errors.New("not found two users")
)
