package websocket

import "errors"

var (
	ErrUpgradingToWebsocket     = errors.New("error upgrading to websockets")
	ErrSuchMessageNameNoExist   = errors.New("such message name does not exist")
	ErrInvalidMessageBodyFormat = errors.New("invalid message body format")
)
