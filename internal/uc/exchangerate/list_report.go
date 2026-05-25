package exchangerate

import (
	"context"
	"time"

	domCurrency "github.com/rlapenok/exchanger/internal/domain/currency"
	domain "github.com/rlapenok/exchanger/internal/domain/exchangerate"
	"github.com/rlapenok/exchanger/internal/uc"
)

// ListReportInput is the input for exchange rate reports.
type ListReportInput struct {
	From       string
	To         string
	BaseCode   string
	QuoteCode  string
	Pagination uc.Pagination
}

// ListReportUseCase returns exchange rate history for reports.
type ListReportUseCase struct {
	repo Repository
}

// NewListReportUseCase creates a new ListReportUseCase.
func NewListReportUseCase(repo Repository) *ListReportUseCase {
	return &ListReportUseCase{repo: repo}
}

// Execute returns report rows filtered by date range and optional pair.
func (uc *ListReportUseCase) Execute(
	ctx context.Context,
	input ListReportInput,
) ([]domain.History, error) {
	from, err := time.Parse("2006-01-02", input.From)
	if err != nil {
		return nil, domain.ErrInvalidDateRange
	}

	to, err := time.Parse("2006-01-02", input.To)
	if err != nil {
		return nil, domain.ErrInvalidDateRange
	}

	if to.Before(from) {
		return nil, domain.ErrInvalidDateRange
	}

	filter := ReportFilter{
		From: from,
		To:   to,
	}

	if input.BaseCode != "" {
		code, err := domCurrency.NewCode(input.BaseCode)
		if err != nil {
			return nil, err
		}
		filter.BaseCode = code.Value()
	}

	if input.QuoteCode != "" {
		code, err := domCurrency.NewCode(input.QuoteCode)
		if err != nil {
			return nil, err
		}
		filter.QuoteCode = code.Value()
	}

	return uc.repo.ListReport(ctx, filter, input.Pagination)
}
