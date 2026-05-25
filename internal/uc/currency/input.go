package currency

import (
	domain "github.com/rlapenok/exchanger/internal/domain/currency"
)

func buildCurrency(code, name, symbol string, minorUnit uint8) (domain.Currency, error) {
	domainCode, err := domain.NewCode(code)
	if err != nil {
		return domain.Currency{}, err
	}

	domainName, err := domain.NewName(name)
	if err != nil {
		return domain.Currency{}, err
	}

	domainSymbol, err := domain.NewSymbol(symbol)
	if err != nil {
		return domain.Currency{}, err
	}

	domainMinorUnit, err := domain.NewMinorUnit(minorUnit)
	if err != nil {
		return domain.Currency{}, err
	}

	return domain.NewCurrency(domainCode, domainName, domainSymbol, domainMinorUnit), nil
}
