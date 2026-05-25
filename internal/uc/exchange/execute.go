package exchange

import (
	"context"
	"errors"
	"strconv"

	domCurrency "github.com/rlapenok/exchanger/internal/domain/currency"
	domain "github.com/rlapenok/exchanger/internal/domain/exchange"
	domRate "github.com/rlapenok/exchanger/internal/domain/exchangerate"
	ucexrate "github.com/rlapenok/exchanger/internal/uc/exchangerate"
)

// ExecuteInput is the input for performing an exchange.
type ExecuteInput struct {
	OperatorName string
	SessionID    string
	BaseCode     string
	QuoteCode    string
	Side         string
	Amount       string
}

// ExecuteUseCase performs a currency exchange.
type ExecuteUseCase struct {
	rateRepo     ucexrate.Repository
	exchangeRepo Repository
}

// NewExecuteUseCase creates a new ExecuteUseCase.
func NewExecuteUseCase(rateRepo ucexrate.Repository, exchangeRepo Repository) *ExecuteUseCase {
	return &ExecuteUseCase{
		rateRepo:     rateRepo,
		exchangeRepo: exchangeRepo,
	}
}

// Execute performs an exchange using the current live rate.
func (uc *ExecuteUseCase) Execute(ctx context.Context, input ExecuteInput) (domain.Exchange, error) {
	baseCode, err := domCurrency.NewCode(input.BaseCode)
	if err != nil {
		return domain.Exchange{}, err
	}

	quoteCode, err := domCurrency.NewCode(input.QuoteCode)
	if err != nil {
		return domain.Exchange{}, err
	}

	side, err := domain.NewSide(input.Side)
	if err != nil {
		return domain.Exchange{}, err
	}

	amount, err := domain.NewAmount(input.Amount)
	if err != nil {
		return domain.Exchange{}, err
	}

	rateRow, err := uc.rateRepo.GetByPair(ctx, baseCode.Value(), quoteCode.Value())
	if err != nil {
		return domain.Exchange{}, mapRateError(err)
	}

	appliedRate, err := selectRate(rateRow, side)
	if err != nil {
		return domain.Exchange{}, err
	}

	resultAmount, err := multiplyAmount(amount, appliedRate)
	if err != nil {
		return domain.Exchange{}, err
	}

	exchange := domain.NewExchange(
		input.OperatorName,
		input.SessionID,
		baseCode,
		quoteCode,
		side,
		amount,
		appliedRate,
		resultAmount,
	)

	return uc.exchangeRepo.Create(ctx, exchange)
}

func selectRate(rate domRate.ExchangeRate, side domain.Side) (domRate.Rate, error) {
	switch side {
	case domain.SideBuy:
		if !rate.IsSellActive() {
			return "", domain.ErrRateNotActive
		}
		return rate.SellRate(), nil
	case domain.SideSell:
		if !rate.IsBuyActive() {
			return "", domain.ErrRateNotActive
		}
		return rate.BuyRate(), nil
	default:
		return "", domain.ErrInvalidSide
	}
}

func multiplyAmount(amount domain.Amount, rate domRate.Rate) (domain.Amount, error) {
	amountValue, err := strconv.ParseFloat(amount.Value(), 64)
	if err != nil {
		return "", domain.ErrInvalidAmount
	}

	rateValue, err := strconv.ParseFloat(rate.Value(), 64)
	if err != nil {
		return "", domRate.ErrInvalidRate
	}

	result := strconv.FormatFloat(amountValue*rateValue, 'f', 8, 64)
	return domain.NewAmount(result)
}

func mapRateError(err error) error {
	if errors.Is(err, domRate.ErrNotFound) {
		return domain.ErrRateNotFound
	}

	return err
}
