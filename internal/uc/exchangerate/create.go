package exchangerate

import (
	"context"

	domain "github.com/rlapenok/exchanger/internal/domain/exchangerate"
)

// CreateInput is the input for creating an exchange rate.
type CreateInput struct {
	BaseCode     string
	QuoteCode    string
	BuyRate      string
	SellRate     string
	IsBuyActive  bool
	IsSellActive bool
}

// CreateUseCase creates a new exchange rate.
type CreateUseCase struct {
	repo Repository
}

// NewCreateUseCase creates a new CreateUseCase.
func NewCreateUseCase(repo Repository) *CreateUseCase {
	return &CreateUseCase{repo: repo}
}

// Execute creates a new exchange rate and initial history snapshot.
func (uc *CreateUseCase) Execute(ctx context.Context, input CreateInput) (domain.ExchangeRate, error) {
	rate, err := buildNewExchangeRate(
		input.BaseCode,
		input.QuoteCode,
		input.BuyRate,
		input.SellRate,
		input.IsBuyActive,
		input.IsSellActive,
	)
	if err != nil {
		return domain.ExchangeRate{}, err
	}

	created, err := uc.repo.Create(ctx, rate)
	if err != nil {
		return domain.ExchangeRate{}, err
	}

	return created, nil
}
