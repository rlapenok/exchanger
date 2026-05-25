package exchangerate

import (
	"context"

	domain "github.com/rlapenok/exchanger/internal/domain/exchangerate"
)

// ListLiveUseCase returns current exchange rates.
type ListLiveUseCase struct {
	repo Repository
}

// NewListLiveUseCase creates a new ListLiveUseCase.
func NewListLiveUseCase(repo Repository) *ListLiveUseCase {
	return &ListLiveUseCase{repo: repo}
}

// Execute returns live exchange rates.
func (uc *ListLiveUseCase) Execute(ctx context.Context) ([]domain.ExchangeRate, error) {
	return uc.repo.ListLive(ctx)
}
