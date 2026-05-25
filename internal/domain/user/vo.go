package user

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const (
	// RoleAdmin is the admin role
	RoleAdmin Role = "admin"
	// RoleOperator is the operator role
	RoleOperator Role = "operator"
)

// Name value object
type Name string

// NewName creates a new Name
func NewName(name string) (Name, error) {
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		return "", ErrInvalidName
	}
	return Name(name), nil
}

// Value returns the value of the Name
func (n Name) Value() string {
	return string(n)
}

// RehydrateName rehydrates the name
func RehydrateName(name string) Name {
	return Name(name)
}

// Password value object
type Password string

// NewPassword creates a new Password
func NewPassword(password string) (Password, error) {
	password = strings.TrimSpace(password)
	if len(password) == 0 {
		return "", ErrInvalidPassword
	}

	return Password(password), nil
}

// Hash hashes the password
func (p Password) Hash() (PasswordHash, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(p),
		bcrypt.DefaultCost)
	if err != nil {
		return "", ErrHashPassword
	}
	return PasswordHash(hash), nil
}

// Value returns the value of the Password
func (p Password) Value() string {
	return string(p)
}

// PasswordHash value object
type PasswordHash string

// NewPasswordHash creates a new PasswordHash
func NewPasswordHash(hash string) (PasswordHash, error) {
	hash = strings.TrimSpace(hash)
	if len(hash) == 0 {
		return "", ErrInvalidPasswordHash
	}
	return PasswordHash(hash), nil
}

// Value returns the value of the PasswordHash
func (h PasswordHash) Value() string {
	return string(h)
}

// RehydratePasswordHash rehydrates the password hash
func RehydratePasswordHash(password string) PasswordHash {
	return PasswordHash(password)
}

// Compare compares the hash with the raw password
func (h PasswordHash) Compare(password Password) error {
	if err := bcrypt.CompareHashAndPassword([]byte(h), []byte(password)); err != nil {
		return ErrInvalidCredentials
	}
	return nil
}

// Role value object
type Role string

// NewRole creates a new Role
func NewRole(role string) (Role, error) {
	role = strings.TrimSpace(role)

	switch Role(role) {
	case RoleAdmin, RoleOperator:
		return Role(role), nil
	default:
		return "", ErrInvalidRole
	}
}

// Value returns the value of the Role
func (r Role) Value() string {
	return string(r)
}

// RehydrateRole rehydrates the role
func RehydrateRole(role string) Role {
	return Role(role)
}
