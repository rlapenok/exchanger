package currency

import (
	"context"

	domain "github.com/rlapenok/exchanger/internal/domain/currency"
	"github.com/rlapenok/exchanger/internal/uc"
)

type GetAllCurrenciesUseCase struct {
	currencyRepository Repository
}

// NewGetAllCurrenciesUseCase creates a new GetAllCurrenciesUseCase
func NewGetAllCurrenciesUseCase(currencyRepository Repository) *GetAllCurrenciesUseCase {
	return &GetAllCurrenciesUseCase{
		currencyRepository: currencyRepository,
	}
}

// Execute executes the GetAllCurrenciesUseCase
func (uc *GetAllCurrenciesUseCase) Execute(ctx context.Context, pagination uc.Pagination) ([]domain.Currency, error) {
	return uc.currencyRepository.GetAll(ctx, pagination)
}
