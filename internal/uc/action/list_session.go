package action

import (
	"context"
	"encoding/json"
	"time"

	domain "github.com/rlapenok/exchanger/internal/domain/action"
	domUser "github.com/rlapenok/exchanger/internal/domain/user"
	"github.com/rlapenok/exchanger/internal/uc"
)

// ListSessionInput is the input for listing session actions.
type ListSessionInput struct {
	ActorName string
	SessionID string
	Pagination uc.Pagination
}

// SessionActionOutput is a single session action for API responses.
type SessionActionOutput struct {
	Request   RequestOutput `json:"request"`
	CreatedAt time.Time     `json:"created_at"`
}

// RequestOutput is the HTTP request snapshot in API responses.
type RequestOutput struct {
	Method string          `json:"method"`
	Path   string          `json:"path"`
	Query  string          `json:"query,omitempty"`
	Body   json.RawMessage `json:"body,omitempty"`
	Status int             `json:"status"`
}

// ListSessionUseCase returns actions for the current session.
type ListSessionUseCase struct {
	repo Repository
}

// NewListSessionUseCase creates a new ListSessionUseCase.
func NewListSessionUseCase(repo Repository) *ListSessionUseCase {
	return &ListSessionUseCase{repo: repo}
}

// Execute returns actions for the given actor and session.
func (uc *ListSessionUseCase) Execute(
	ctx context.Context,
	input ListSessionInput,
) ([]SessionActionOutput, error) {
	if _, err := domUser.NewName(input.ActorName); err != nil {
		return nil, err
	}
	if _, err := domain.NewSessionID(input.SessionID); err != nil {
		return nil, err
	}

	actions, err := uc.repo.ListBySession(
		ctx,
		input.ActorName,
		input.SessionID,
		input.Pagination,
	)
	if err != nil {
		return nil, err
	}

	output := make([]SessionActionOutput, len(actions))
	for i, action := range actions {
		request := action.Request()
		output[i] = SessionActionOutput{
			Request: RequestOutput{
				Method: request.Method,
				Path:   request.Path,
				Query:  request.Query,
				Body:   request.Body,
				Status: request.Status,
			},
			CreatedAt: action.CreatedAt(),
		}
	}

	return output, nil
}
