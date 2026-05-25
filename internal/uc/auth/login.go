package auth

import (
	"context"
	"database/sql"
	"errors"

	domUser "github.com/rlapenok/exchanger/internal/domain/user"
	"github.com/rlapenok/exchanger/internal/uc/user"
)

type LoginUseCase struct {
	repo user.Repo
}

// NewLoginUseCase creates a new LoginUseCase
func NewLoginUseCase(repo user.Repo) *LoginUseCase {
	return &LoginUseCase{repo: repo}
}

func (uc *LoginUseCase) Login(ctx context.Context, req LoginInput) (LoginOutput, error) {
	name, err := domUser.NewName(req.Name)
	if err != nil {
		return LoginOutput{}, err
	}

	password, err := domUser.NewPassword(req.Password)
	if err != nil {
		return LoginOutput{}, err
	}

	u, err := uc.repo.GetByName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return LoginOutput{}, domUser.ErrInvalidCredentials
		}

		return LoginOutput{}, err
	}

	if err := u.PasswordHash().Compare(password); err != nil {
		return LoginOutput{}, err
	}

	return LoginOutput{
		Name: u.Name(),
		Role: u.Role(),
	}, nil
}

// LoginInput is the request body for the login use case
type LoginInput struct {
	Name     string
	Password string
}

// LoginOutput is the output for the login use case
type LoginOutput struct {
	Name domUser.Name
	Role domUser.Role
}
