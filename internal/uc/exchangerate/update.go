package exchangerate

import (
	"context"

	domain "github.com/rlapenok/exchanger/internal/domain/exchangerate"
)

// UpdateInput is the input for updating an exchange rate.
type UpdateInput struct {
	ID           string
	BuyRate      string
	SellRate     string
	IsBuyActive  bool
	IsSellActive bool
}

// UpdateUseCase updates an exchange rate.
type UpdateUseCase struct {
	repo Repository
}

// NewUpdateUseCase creates a new UpdateUseCase.
func NewUpdateUseCase(repo Repository) *UpdateUseCase {
	return &UpdateUseCase{repo: repo}
}

// Execute updates an exchange rate and appends history.
func (uc *UpdateUseCase) Execute(ctx context.Context, input UpdateInput) (domain.ExchangeRate, error) {
	rateID, err := domain.NewID(input.ID)
	if err != nil {
		return domain.ExchangeRate{}, err
	}

	current, err := uc.repo.GetByID(ctx, rateID)
	if err != nil {
		return domain.ExchangeRate{}, err
	}

	updated, err := buildExchangeRate(
		current.ID().Value(),
		current.BaseCode().Value(),
		current.QuoteCode().Value(),
		input.BuyRate,
		input.SellRate,
		input.IsBuyActive,
		input.IsSellActive,
	)
	if err != nil {
		return domain.ExchangeRate{}, err
	}

	if err := uc.repo.Update(ctx, updated); err != nil {
		return domain.ExchangeRate{}, err
	}

	result, err := uc.repo.GetByID(ctx, rateID)
	if err != nil {
		return domain.ExchangeRate{}, err
	}

	return result, nil
}
