package user

import (
	"context"

	domainuser "github.com/rlapenok/exchanger/internal/domain/user"
)

// Repo is the repository required by user use cases.
type Repo interface {
	GetByName(ctx context.Context, name domainuser.Name) (domainuser.User, error)
}
