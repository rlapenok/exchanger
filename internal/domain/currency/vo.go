package currency

import "strings"

// CurrencyCode value object
type Code string

// NewCurrencyCode creates a new CurrencyCode value object
func NewCode(code string) (Code, error) {
	code = strings.TrimSpace(code)
	if len(code) != 3 {
		return "", ErrInvalidCode
	}

	code = strings.ToUpper(code)
	return Code(code), nil
}

// Value returns the value of the CurrencyCode
func (c Code) Value() string {
	return string(c)
}

// RehydrateCode rehydrates the code
func RehydrateCode(code string) Code {
	return Code(code)
}

// Name value object
type Name string

// NewName creates a new Name value object
func NewName(name string) (Name, error) {
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		return "", ErrInvalidName
	}

	return Name(name), nil
}

// Value returns the value of the Name
func (c Name) Value() string {
	return string(c)
}

// RehydrateName rehydrates the name
func RehydrateName(name string) Name {
	return Name(name)
}

// Symbol value object
type Symbol string

// NewSymbol creates a new Symbol value object
func NewSymbol(symbol string) (Symbol, error) {
	symbol = strings.TrimSpace(symbol)
	if len(symbol) == 0 {
		return "", ErrInvalidSymbol
	}

	return Symbol(symbol), nil
}

// Value returns the value of the Symbol
func (c Symbol) Value() string {
	return string(c)
}

// RehydrateSymbol rehydrates the symbol
func RehydrateSymbol(symbol string) Symbol {
	return Symbol(symbol)
}

// MinorUnit value object
type MinorUnit uint8

// NewMinorUnit creates a new MinorUnit value object
func NewMinorUnit(minorUnit uint8) (MinorUnit, error) {
	if minorUnit > 4 {
		return 0, ErrInvalidMinorUnit
	}

	return MinorUnit(minorUnit), nil
}

// Value returns the value of the MinorUnit
func (m MinorUnit) Value() uint8 {
	return uint8(m)
}

// Int8 returns the value of the MinorUnit as an int8
func (m MinorUnit) Int16() int16 {
	return int16(m)
}

// RehydrateMinorUnit rehydrates the minor unit
func RehydrateMinorUnit(minorUnit int16) MinorUnit {
	return MinorUnit(uint8(minorUnit))
}
