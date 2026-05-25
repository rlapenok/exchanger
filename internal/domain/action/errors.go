package action

import "errors"

var (
	ErrInvalidSessionID = errors.New("invalid session id")
	ErrInvalidMethod    = errors.New("invalid http method")
	ErrInvalidPath      = errors.New("invalid path")
	ErrInvalidStatus    = errors.New("invalid http status")
)
