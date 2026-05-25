package exchangerate

import (
	"context"

	domain "github.com/rlapenok/exchanger/internal/domain/exchangerate"
	"github.com/rlapenok/exchanger/internal/uc"
)

// ListHistoryInput is the input for listing exchange rate history.
type ListHistoryInput struct {
	ID         string
	Pagination uc.Pagination
}

// ListHistoryUseCase returns history for one exchange rate.
type ListHistoryUseCase struct {
	repo Repository
}

// NewListHistoryUseCase creates a new ListHistoryUseCase.
func NewListHistoryUseCase(repo Repository) *ListHistoryUseCase {
	return &ListHistoryUseCase{repo: repo}
}

// Execute returns history rows for the given exchange rate.
func (uc *ListHistoryUseCase) Execute(
	ctx context.Context,
	input ListHistoryInput,
) ([]domain.History, error) {
	rateID, err := domain.NewID(input.ID)
	if err != nil {
		return nil, err
	}

	return uc.repo.ListHistory(ctx, rateID, input.Pagination)
}
