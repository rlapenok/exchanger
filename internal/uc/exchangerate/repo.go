package exchangerate

import (
	"context"
	"time"

	domain "github.com/rlapenok/exchanger/internal/domain/exchangerate"
	"github.com/rlapenok/exchanger/internal/uc"
)

// ReportFilter filters exchange rate history for reports.
type ReportFilter struct {
	From      time.Time
	To        time.Time
	BaseCode  string
	QuoteCode string
}

// Repository persists and reads exchange rates.
type Repository interface {
	ListLive(ctx context.Context) ([]domain.ExchangeRate, error)
	GetByID(ctx context.Context, id domain.ID) (domain.ExchangeRate, error)
	GetByPair(ctx context.Context, baseCode, quoteCode string) (domain.ExchangeRate, error)
	Create(ctx context.Context, rate domain.ExchangeRate) (domain.ExchangeRate, error)
	Update(ctx context.Context, rate domain.ExchangeRate) error
	ListHistory(
		ctx context.Context,
		id domain.ID,
		pagination uc.Pagination,
	) ([]domain.History, error)
	ListReport(
		ctx context.Context,
		filter ReportFilter,
		pagination uc.Pagination,
	) ([]domain.History, error)
}
