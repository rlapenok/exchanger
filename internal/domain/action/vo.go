package action

import "strings"

// SessionID identifies a user session.
type SessionID string

// NewSessionID creates a new SessionID.
func NewSessionID(value string) (SessionID, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", ErrInvalidSessionID
	}

	return SessionID(value), nil
}

// Value returns the raw session id.
func (id SessionID) Value() string {
	return string(id)
}

// RehydrateSessionID restores a SessionID from storage.
func RehydrateSessionID(value string) SessionID {
	return SessionID(value)
}
