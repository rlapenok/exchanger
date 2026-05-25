package currency

import (
	"context"

	domain "github.com/rlapenok/exchanger/internal/domain/currency"
)

// GetByCodeUseCase is the use case for getting a currency by code
type GetByCodeUseCase struct {
	currencyRepository Repository
}

// NewGetByCodeUseCase creates a new GetByCodeUseCase
func NewGetByCodeUseCase(currencyRepository Repository) *GetByCodeUseCase {
	return &GetByCodeUseCase{currencyRepository: currencyRepository}
}

// Execute executes the GetByCodeUseCase
func (uc *GetByCodeUseCase) Execute(ctx context.Context, input GetByCodeInput) (domain.Currency, error) {
	code, err := domain.NewCode(input.Code)
	if err != nil {
		return domain.Currency{}, err
	}

	return uc.currencyRepository.GetByCode(ctx, code)
}

// GetByCodeInput is the input for the GetByCodeUseCase
type GetByCodeInput struct {
	Code string
}
