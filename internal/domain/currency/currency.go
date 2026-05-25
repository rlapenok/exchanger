package currency

// Currency entity
type Currency struct {
	code      Code
	name      Name
	symbol    Symbol
	minorUnit MinorUnit
}

// NewCurrency creates a new Currency entity
func NewCurrency(
	code Code,
	name Name,
	symbol Symbol,
	minorUnit MinorUnit,
) Currency {
	return Currency{
		code:      code,
		name:      name,
		symbol:    symbol,
		minorUnit: minorUnit,
	}
}

// Code returns the code of the Currency
func (c Currency) Code() Code {
	return c.code
}

// Name returns the name of the Currency
func (c Currency) Name() Name {
	return c.name
}

// Symbol returns the symbol of the Currency
func (c Currency) Symbol() Symbol {
	return c.symbol
}

// MinorUnit returns the minor unit of the Currency
func (c Currency) MinorUnit() MinorUnit {
	return c.minorUnit
}

// RehydrateCurrency rehydrates the currency
func RehydrateCurrency(code, name, symbol string, minorUnit int16) Currency {
	return NewCurrency(
		RehydrateCode(code),
		RehydrateName(name),
		RehydrateSymbol(symbol),
		RehydrateMinorUnit(minorUnit),
	)
}
