package exchangerate

import (
	"time"

	domCurrency "github.com/rlapenok/exchanger/internal/domain/currency"
)

// History is a historical exchange rate snapshot.
type History struct {
	exchangeRateID ID
	baseCode       domCurrency.Code
	quoteCode      domCurrency.Code
	buyRate        Rate
	sellRate       Rate
	validFrom      time.Time
	validTo        time.Time
}

// ExchangeRateID returns the related exchange rate id.
func (h History) ExchangeRateID() ID {
	return h.exchangeRateID
}

// BaseCode returns the base currency code.
func (h History) BaseCode() domCurrency.Code {
	return h.baseCode
}

// QuoteCode returns the quote currency code.
func (h History) QuoteCode() domCurrency.Code {
	return h.quoteCode
}

// BuyRate returns the buy rate.
func (h History) BuyRate() Rate {
	return h.buyRate
}

// SellRate returns the sell rate.
func (h History) SellRate() Rate {
	return h.sellRate
}

// ValidFrom returns when the snapshot started.
func (h History) ValidFrom() time.Time {
	return h.validFrom
}

// ValidTo returns when the snapshot ended.
func (h History) ValidTo() time.Time {
	return h.validTo
}

// RehydrateHistory restores a history row from storage.
func RehydrateHistory(
	exchangeRateID string,
	baseCode string,
	quoteCode string,
	buyRate string,
	sellRate string,
	validFrom time.Time,
	validTo time.Time,
) History {
	return History{
		exchangeRateID: RehydrateID(exchangeRateID),
		baseCode:       domCurrency.RehydrateCode(baseCode),
		quoteCode:      domCurrency.RehydrateCode(quoteCode),
		buyRate:        RehydrateRate(buyRate),
		sellRate:       RehydrateRate(sellRate),
		validFrom:      validFrom,
		validTo:        validTo,
	}
}
