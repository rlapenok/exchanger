package user

import (
	"context"

	"github.com/rlapenok/exchanger/internal/domain/user"
)

// AccountUseCase is the use case for the user account
type AccountUseCase struct {
	repo Repo
}

// NewUseCase creates a new UseCase
func NewAccountUseCase(repo Repo) *AccountUseCase {
	return &AccountUseCase{repo: repo}
}

// Execute executes the profile use case
func (uc *AccountUseCase) Execute(ctx context.Context, input ProfileInput) (ProfileOutput, error) {
	name, err := user.NewName(input.Name)
	if err != nil {
		return ProfileOutput{}, err
	}

	user, err := uc.repo.GetByName(ctx, name)
	if err != nil {
		return ProfileOutput{}, err
	}

	return ProfileOutput{
		Name: user.Name().Value(),
		Role: user.Role().Value(),
	}, nil
}

// ProfileInput is the input for the profile use case
type ProfileInput struct {
	Name string
}

// ProfileOutput is the output for the profile use case
type ProfileOutput struct {
	Name string
	Role string
}
