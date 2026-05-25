package user

import "errors"

var (
	// ErrInvalidName is returned when a name is invalid
	ErrInvalidName = errors.New("name must be not empty")
	// ErrInvalidPassword is returned when a password is invalid
	ErrInvalidPassword    = errors.New("password must be not empty")
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrHashPassword is returned when a password hash is invalid
	ErrHashPassword = errors.New("failed to hash password")
	// ErrInvalidPasswordHash is returned when a password hash is invalid
	ErrInvalidPasswordHash = errors.New("password hash must be not empty")
	// ErrInvalidRole is returned when a role is invalid
	ErrInvalidRole = errors.New("role must be not empty")
)
