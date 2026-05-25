package exchangerate

import "errors"

var (
	ErrInvalidID        = errors.New("invalid exchange rate id")
	ErrInvalidRate      = errors.New("invalid rate")
	ErrSellLessThanBuy  = errors.New("sell rate must be greater than or equal to buy rate")
	ErrSameCurrency     = errors.New("base and quote currencies must differ")
	ErrNotFound         = errors.New("exchange rate not found")
	ErrAlreadyExists    = errors.New("exchange rate already exists")
	ErrInvalidDateRange = errors.New("invalid date range")
)
