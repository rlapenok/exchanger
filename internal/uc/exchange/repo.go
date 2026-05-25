package exchange

import (
	"context"
	"time"

	domain "github.com/rlapenok/exchanger/internal/domain/exchange"
	"github.com/rlapenok/exchanger/internal/uc"
)

// ReportFilter filters exchange operations for reports.
type ReportFilter struct {
	From      time.Time
	To        time.Time
	BaseCode  string
	QuoteCode string
}

// Repository persists and reads exchanges.
type Repository interface {
	Create(ctx context.Context, exchange domain.Exchange) (domain.Exchange, error)
	ListReport(
		ctx context.Context,
		filter ReportFilter,
		pagination uc.Pagination,
	) ([]domain.Exchange, error)
}
