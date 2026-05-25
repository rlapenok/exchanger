package currency

import (
	"context"

	domain "github.com/rlapenok/exchanger/internal/domain/currency"
)

// UpdateUseCase updates an existing currency.
type UpdateUseCase struct {
	currencyRepository Repository
}

// NewUpdateUseCase creates a new UpdateUseCase.
func NewUpdateUseCase(currencyRepository Repository) *UpdateUseCase {
	return &UpdateUseCase{currencyRepository: currencyRepository}
}

// UpdateInput is the input for UpdateUseCase.
type UpdateInput struct {
	Code      string
	Name      string
	Symbol    string
	MinorUnit uint8
}

// Execute executes UpdateUseCase.
func (uc *UpdateUseCase) Execute(ctx context.Context, input UpdateInput) (domain.Currency, error) {
	currency, err := buildCurrency(input.Code, input.Name, input.Symbol, input.MinorUnit)
	if err != nil {
		return domain.Currency{}, err
	}

	if err := uc.currencyRepository.Update(ctx, currency); err != nil {
		return domain.Currency{}, err
	}

	return currency, nil
}
