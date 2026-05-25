package exchangerate

import (
	domain "github.com/rlapenok/exchanger/internal/domain/exchangerate"
	domCurrency "github.com/rlapenok/exchanger/internal/domain/currency"
)

func buildNewExchangeRate(
	baseCode string,
	quoteCode string,
	buyRate string,
	sellRate string,
	isBuyActive bool,
	isSellActive bool,
) (domain.ExchangeRate, error) {
	base, err := domCurrency.NewCode(baseCode)
	if err != nil {
		return domain.ExchangeRate{}, err
	}

	quote, err := domCurrency.NewCode(quoteCode)
	if err != nil {
		return domain.ExchangeRate{}, err
	}

	buy, err := domain.NewRate(buyRate)
	if err != nil {
		return domain.ExchangeRate{}, err
	}

	sell, err := domain.NewRate(sellRate)
	if err != nil {
		return domain.ExchangeRate{}, err
	}

	return domain.NewExchangeRateDraft(
		base,
		quote,
		buy,
		sell,
		isBuyActive,
		isSellActive,
	)
}

func buildExchangeRate(
	id string,
	baseCode string,
	quoteCode string,
	buyRate string,
	sellRate string,
	isBuyActive bool,
	isSellActive bool,
) (domain.ExchangeRate, error) {
	rateID, err := domain.NewID(id)
	if err != nil {
		return domain.ExchangeRate{}, err
	}

	base, err := domCurrency.NewCode(baseCode)
	if err != nil {
		return domain.ExchangeRate{}, err
	}

	quote, err := domCurrency.NewCode(quoteCode)
	if err != nil {
		return domain.ExchangeRate{}, err
	}

	buy, err := domain.NewRate(buyRate)
	if err != nil {
		return domain.ExchangeRate{}, err
	}

	sell, err := domain.NewRate(sellRate)
	if err != nil {
		return domain.ExchangeRate{}, err
	}

	return domain.NewExchangeRate(
		rateID,
		base,
		quote,
		buy,
		sell,
		isBuyActive,
		isSellActive,
	)
}
