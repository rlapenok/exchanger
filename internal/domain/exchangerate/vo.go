package exchangerate

import (
	"strings"
)

// ID identifies an exchange rate row.
type ID string

// NewID creates a new ID.
func NewID(value string) (ID, error) {
	value = strings.TrimSpace(value)
	if len(value) != 36 {
		return "", ErrInvalidID
	}

	return ID(value), nil
}

// Value returns the raw id.
func (id ID) Value() string {
	return string(id)
}

// RehydrateID restores an ID from storage.
func RehydrateID(value string) ID {
	return ID(value)
}
