package exchangerate

import (
	"time"

	domCurrency "github.com/rlapenok/exchanger/internal/domain/currency"
)

// ExchangeRate is the current exchange rate for a currency pair.
type ExchangeRate struct {
	id            ID
	baseCode      domCurrency.Code
	quoteCode     domCurrency.Code
	buyRate       Rate
	sellRate      Rate
	isBuyActive   bool
	isSellActive  bool
	updatedAt     time.Time
}

// NewExchangeRate creates a new exchange rate entity.
func NewExchangeRate(
	id ID,
	baseCode domCurrency.Code,
	quoteCode domCurrency.Code,
	buyRate Rate,
	sellRate Rate,
	isBuyActive bool,
	isSellActive bool,
) (ExchangeRate, error) {
	if baseCode.Value() == quoteCode.Value() {
		return ExchangeRate{}, ErrSameCurrency
	}
	if err := validateRatePair(buyRate, sellRate); err != nil {
		return ExchangeRate{}, err
	}

	return ExchangeRate{
		id:           id,
		baseCode:     baseCode,
		quoteCode:    quoteCode,
		buyRate:      buyRate,
		sellRate:     sellRate,
		isBuyActive:  isBuyActive,
		isSellActive: isSellActive,
	}, nil
}

// NewExchangeRateDraft creates a rate without id for insertion.
func NewExchangeRateDraft(
	baseCode domCurrency.Code,
	quoteCode domCurrency.Code,
	buyRate Rate,
	sellRate Rate,
	isBuyActive bool,
	isSellActive bool,
) (ExchangeRate, error) {
	if baseCode.Value() == quoteCode.Value() {
		return ExchangeRate{}, ErrSameCurrency
	}
	if err := validateRatePair(buyRate, sellRate); err != nil {
		return ExchangeRate{}, err
	}

	return ExchangeRate{
		baseCode:     baseCode,
		quoteCode:    quoteCode,
		buyRate:      buyRate,
		sellRate:     sellRate,
		isBuyActive:  isBuyActive,
		isSellActive: isSellActive,
	}, nil
}

// ID returns the exchange rate identifier.
func (r ExchangeRate) ID() ID {
	return r.id
}

// BaseCode returns the base currency code.
func (r ExchangeRate) BaseCode() domCurrency.Code {
	return r.baseCode
}

// QuoteCode returns the quote currency code.
func (r ExchangeRate) QuoteCode() domCurrency.Code {
	return r.quoteCode
}

// BuyRate returns the buy rate.
func (r ExchangeRate) BuyRate() Rate {
	return r.buyRate
}

// SellRate returns the sell rate.
func (r ExchangeRate) SellRate() Rate {
	return r.sellRate
}

// IsBuyActive returns whether buy is active.
func (r ExchangeRate) IsBuyActive() bool {
	return r.isBuyActive
}

// IsSellActive returns whether sell is active.
func (r ExchangeRate) IsSellActive() bool {
	return r.isSellActive
}

// UpdatedAt returns the last update timestamp.
func (r ExchangeRate) UpdatedAt() time.Time {
	return r.updatedAt
}

// RehydrateExchangeRate restores an exchange rate from storage.
func RehydrateExchangeRate(
	id string,
	baseCode string,
	quoteCode string,
	buyRate string,
	sellRate string,
	isBuyActive bool,
	isSellActive bool,
	updatedAt time.Time,
) ExchangeRate {
	return ExchangeRate{
		id:           RehydrateID(id),
		baseCode:     domCurrency.RehydrateCode(baseCode),
		quoteCode:    domCurrency.RehydrateCode(quoteCode),
		buyRate:      RehydrateRate(buyRate),
		sellRate:     RehydrateRate(sellRate),
		isBuyActive:  isBuyActive,
		isSellActive: isSellActive,
		updatedAt:    updatedAt,
	}
}
