package exchange

import (
	"time"

	domCurrency "github.com/rlapenok/exchanger/internal/domain/currency"
	domRate "github.com/rlapenok/exchanger/internal/domain/exchangerate"
)

// Exchange is a completed currency exchange operation.
type Exchange struct {
	id           string
	operatorName string
	sessionID    string
	baseCode     domCurrency.Code
	quoteCode    domCurrency.Code
	side         Side
	amount       Amount
	rate         domRate.Rate
	resultAmount Amount
	createdAt    time.Time
}

// NewExchange creates a new exchange record.
func NewExchange(
	operatorName string,
	sessionID string,
	baseCode domCurrency.Code,
	quoteCode domCurrency.Code,
	side Side,
	amount Amount,
	rate domRate.Rate,
	resultAmount Amount,
) Exchange {
	return Exchange{
		operatorName: operatorName,
		sessionID:    sessionID,
		baseCode:     baseCode,
		quoteCode:    quoteCode,
		side:         side,
		amount:       amount,
		rate:         rate,
		resultAmount: resultAmount,
	}
}

// ID returns the exchange identifier.
func (e Exchange) ID() string {
	return e.id
}

// OperatorName returns the operator name.
func (e Exchange) OperatorName() string {
	return e.operatorName
}

// SessionID returns the session identifier.
func (e Exchange) SessionID() string {
	return e.sessionID
}

// BaseCode returns the base currency code.
func (e Exchange) BaseCode() domCurrency.Code {
	return e.baseCode
}

// QuoteCode returns the quote currency code.
func (e Exchange) QuoteCode() domCurrency.Code {
	return e.quoteCode
}

// Side returns the exchange side.
func (e Exchange) Side() Side {
	return e.side
}

// Amount returns the exchanged amount in base currency.
func (e Exchange) Amount() Amount {
	return e.amount
}

// Rate returns the applied rate.
func (e Exchange) Rate() domRate.Rate {
	return e.rate
}

// ResultAmount returns the result amount in quote currency.
func (e Exchange) ResultAmount() Amount {
	return e.resultAmount
}

// CreatedAt returns the exchange timestamp.
func (e Exchange) CreatedAt() time.Time {
	return e.createdAt
}

// RehydrateExchange restores an exchange from storage.
func RehydrateExchange(
	id string,
	operatorName string,
	sessionID string,
	baseCode string,
	quoteCode string,
	side string,
	amount string,
	rate string,
	resultAmount string,
	createdAt time.Time,
) Exchange {
	return Exchange{
		id:           id,
		operatorName: operatorName,
		sessionID:    sessionID,
		baseCode:     domCurrency.RehydrateCode(baseCode),
		quoteCode:    domCurrency.RehydrateCode(quoteCode),
		side:         Side(side),
		amount:       RehydrateAmount(amount),
		rate:         domRate.RehydrateRate(rate),
		resultAmount: RehydrateAmount(resultAmount),
		createdAt:    createdAt,
	}
}
