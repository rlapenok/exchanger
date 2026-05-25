package action

import (
	"context"
	"encoding/json"

	domain "github.com/rlapenok/exchanger/internal/domain/action"
	domUser "github.com/rlapenok/exchanger/internal/domain/user"
)

// RecordInput is the input for recording a user action.
type RecordInput struct {
	ActorName string
	SessionID string
	Method    string
	Path      string
	Query     string
	Body      json.RawMessage
	Status    int
}

// RecordUseCase records a user action.
type RecordUseCase struct {
	repo Repository
}

// NewRecordUseCase creates a new RecordUseCase.
func NewRecordUseCase(repo Repository) *RecordUseCase {
	return &RecordUseCase{repo: repo}
}

// Execute validates input and persists the action.
func (uc *RecordUseCase) Execute(ctx context.Context, input RecordInput) error {
	actorName, err := domUser.NewName(input.ActorName)
	if err != nil {
		return err
	}

	sessionID, err := domain.NewSessionID(input.SessionID)
	if err != nil {
		return err
	}

	request, err := domain.NewRequestSnapshot(
		input.Method,
		input.Path,
		input.Query,
		input.Body,
		input.Status,
	)
	if err != nil {
		return err
	}

	action := domain.NewAction(actorName, sessionID, request)
	return uc.repo.Record(ctx, action)
}
