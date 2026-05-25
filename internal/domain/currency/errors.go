package currency

import "errors"

var (
	// ErrInvalidCurrencyCode is returned when a currency code is invalid
	ErrInvalidCode = errors.New("invalid currency code")
	// ErrInvalidName is returned when a currency name is invalid
	ErrInvalidName = errors.New("name must be not empty")
	// ErrInvalidSymbol is returned when a currency symbol is invalid
	ErrInvalidSymbol = errors.New("symbol must be not empty")
	// ErrInvalidMinorUnit is returned when a currency minor unit is invalid
	ErrInvalidMinorUnit = errors.New("minor unit must be between 0 and 4")
	// ErrNotFound is returned when a currency is not found
	ErrNotFound = errors.New("currency not found")
	// ErrAlreadyExists is returned when a currency already exists
	ErrAlreadyExists = errors.New("currency already exists")
	// ErrInUse is returned when a currency is referenced by exchange rates
	ErrInUse = errors.New("currency is in use")
)
