package currency

import (
	"context"

	domain "github.com/rlapenok/exchanger/internal/domain/currency"
)

// CreateUseCase creates a new currency.
type CreateUseCase struct {
	currencyRepository Repository
}

// NewCreateUseCase creates a new CreateUseCase.
func NewCreateUseCase(currencyRepository Repository) *CreateUseCase {
	return &CreateUseCase{currencyRepository: currencyRepository}
}

// CreateInput is the input for CreateUseCase.
type CreateInput struct {
	Code      string
	Name      string
	Symbol    string
	MinorUnit uint8
}

// Execute executes CreateUseCase.
func (uc *CreateUseCase) Execute(ctx context.Context, input CreateInput) (domain.Currency, error) {
	currency, err := buildCurrency(input.Code, input.Name, input.Symbol, input.MinorUnit)
	if err != nil {
		return domain.Currency{}, err
	}

	if err := uc.currencyRepository.Create(ctx, currency); err != nil {
		return domain.Currency{}, err
	}

	return currency, nil
}
