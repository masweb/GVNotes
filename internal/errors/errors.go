package errors

import "errors"

var (
	ErrNotFound       = errors.New("not found")
	ErrConflict       = errors.New("conflict")
	ErrForbidden      = errors.New("forbidden")
	ErrUnauthenticated = errors.New("unauthenticated")
	ErrInvalidInput   = errors.New("invalid input")
)
