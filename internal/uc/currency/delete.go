package currency

import (
	"context"

	domain "github.com/rlapenok/exchanger/internal/domain/currency"
)

// DeleteUseCase deletes a currency.
type DeleteUseCase struct {
	currencyRepository Repository
}

// NewDeleteUseCase creates a new DeleteUseCase.
func NewDeleteUseCase(currencyRepository Repository) *DeleteUseCase {
	return &DeleteUseCase{currencyRepository: currencyRepository}
}

// DeleteInput is the input for DeleteUseCase.
type DeleteInput struct {
	Code string
}

// Execute executes DeleteUseCase.
func (uc *DeleteUseCase) Execute(ctx context.Context, input DeleteInput) error {
	code, err := domain.NewCode(input.Code)
	if err != nil {
		return err
	}

	return uc.currencyRepository.Delete(ctx, code)
}
