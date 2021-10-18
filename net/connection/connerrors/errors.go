package connerrors

import "errors"

var (
	ErrConnectClosed = errors.New("connect is closed")
)
