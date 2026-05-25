package exchangerate

import (
	"strconv"
	"strings"
)

// Rate is a positive decimal exchange rate stored as string.
type Rate string

// NewRate validates and creates a rate.
func NewRate(value string) (Rate, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", ErrInvalidRate
	}

	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil || parsed <= 0 {
		return "", ErrInvalidRate
	}

	return Rate(value), nil
}

// Value returns the raw rate string.
func (r Rate) Value() string {
	return string(r)
}

// RehydrateRate restores a rate from storage.
func RehydrateRate(value string) Rate {
	return Rate(value)
}

func validateRatePair(buy Rate, sell Rate) error {
	buyValue, err := strconv.ParseFloat(buy.Value(), 64)
	if err != nil {
		return ErrInvalidRate
	}

	sellValue, err := strconv.ParseFloat(sell.Value(), 64)
	if err != nil {
		return ErrInvalidRate
	}

	if sellValue < buyValue {
		return ErrSellLessThanBuy
	}

	return nil
}
