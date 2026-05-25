package action

import (
	"context"

	domain "github.com/rlapenok/exchanger/internal/domain/action"
	"github.com/rlapenok/exchanger/internal/uc"
)

// Repository persists and reads user actions.
type Repository interface {
	Record(ctx context.Context, action domain.Action) error
	ListBySession(
		ctx context.Context,
		actorName string,
		sessionID string,
		pagination uc.Pagination,
	) ([]domain.Action, error)
}
