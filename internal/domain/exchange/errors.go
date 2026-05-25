package exchange

import "errors"

var (
	ErrInvalidSide        = errors.New("invalid exchange side")
	ErrInvalidAmount      = errors.New("invalid amount")
	ErrInvalidDateRange   = errors.New("invalid date range")
	ErrRateNotActive      = errors.New("exchange rate side is not active")
	ErrRateNotFound       = errors.New("exchange rate not found")
)
