package exchange

import (
	"context"
	"time"

	domCurrency "github.com/rlapenok/exchanger/internal/domain/currency"
	domain "github.com/rlapenok/exchanger/internal/domain/exchange"
	"github.com/rlapenok/exchanger/internal/uc"
)

// ListReportInput is the input for listing exchanges in a report.
type ListReportInput struct {
	From       string
	To         string
	BaseCode   string
	QuoteCode  string
	Pagination uc.Pagination
}

// ListReportUseCase returns exchanges for admin reports.
type ListReportUseCase struct {
	repo Repository
}

// NewListReportUseCase creates a new ListReportUseCase.
func NewListReportUseCase(repo Repository) *ListReportUseCase {
	return &ListReportUseCase{repo: repo}
}

// Execute returns exchanges filtered by date range and optional pair.
func (uc *ListReportUseCase) Execute(
	ctx context.Context,
	input ListReportInput,
) ([]domain.Exchange, error) {
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
