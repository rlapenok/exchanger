package currency

import (
	"context"

	domain "github.com/rlapenok/exchanger/internal/domain/currency"
	"github.com/rlapenok/exchanger/internal/uc"
)

type Repository interface {
	// Create creates a new currency
	Create(ctx context.Context, currency domain.Currency) error
	// GetAll returns all currencies with pagination
	GetAll(ctx context.Context, pagination uc.Pagination) ([]domain.Currency, error)
	// GetByCode returns a currency by code
	GetByCode(ctx context.Context, code domain.Code) (domain.Currency, error)
	// Update updates a currency
	Update(ctx context.Context, currency domain.Currency) error
	// Delete deletes a currency
	Delete(ctx context.Context, code domain.Code) error
}
