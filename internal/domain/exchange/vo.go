package exchange

import (
	"strconv"
	"strings"
)

const (
	SideBuy  Side = "buy"
	SideSell Side = "sell"
)

// Side is the customer exchange direction.
type Side string

// NewSide creates a Side.
func NewSide(value string) (Side, error) {
	switch Side(strings.TrimSpace(value)) {
	case SideBuy, SideSell:
		return Side(value), nil
	default:
		return "", ErrInvalidSide
	}
}

// Value returns the raw side.
func (s Side) Value() string {
	return string(s)
}

// Amount is a positive exchange amount.
type Amount string

// NewAmount validates and creates an amount.
func NewAmount(value string) (Amount, error) {
	value = strings.TrimSpace(value)
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil || parsed <= 0 {
		return "", ErrInvalidAmount
	}

	return Amount(value), nil
}

// Value returns the raw amount string.
func (a Amount) Value() string {
	return string(a)
}

// RehydrateAmount restores an amount from storage.
func RehydrateAmount(value string) Amount {
	return Amount(value)
}
