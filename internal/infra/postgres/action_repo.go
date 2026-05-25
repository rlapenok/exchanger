package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	domain "github.com/rlapenok/exchanger/internal/domain/action"
	"github.com/rlapenok/exchanger/internal/uc"
)

// ActionRepo persists user actions in Postgres.
type ActionRepo struct {
	db *sql.DB
}

// NewActionRepo creates a new ActionRepo.
func NewActionRepo(db *sql.DB) *ActionRepo {
	return &ActionRepo{db: db}
}

type requestSnapshotRow struct {
	Method string          `json:"method"`
	Path   string          `json:"path"`
	Query  string          `json:"query,omitempty"`
	Body   json.RawMessage `json:"body,omitempty"`
	Status int             `json:"status"`
}

// Record inserts a user action.
func (r *ActionRepo) Record(ctx context.Context, action domain.Action) error {
	requestJSON, err := json.Marshal(requestSnapshotRow{
		Method: action.Request().Method,
		Path:   action.Request().Path,
		Query:  action.Request().Query,
		Body:   action.Request().Body,
		Status: action.Request().Status,
	})
	if err != nil {
		return err
	}

	query := `
		INSERT INTO user_actions (actor_name, session_id, request)
		VALUES ($1, $2, $3::jsonb)
	`
	_, err = r.db.ExecContext(
		ctx,
		query,
		action.ActorName().Value(),
		action.SessionID().Value(),
		string(requestJSON),
	)

	return err
}

// ListBySession returns actions for the given actor and session.
func (r *ActionRepo) ListBySession(
	ctx context.Context,
	actorName string,
	sessionID string,
	pagination uc.Pagination,
) ([]domain.Action, error) {
	query := `
		SELECT id, actor_name, session_id, request, created_at
		FROM user_actions
		WHERE actor_name = $1 AND session_id = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`
	rows, err := r.db.QueryContext(
		ctx,
		query,
		actorName,
		sessionID,
		pagination.Limit,
		pagination.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	actions := make([]domain.Action, 0, pagination.Limit)
	for rows.Next() {
		var (
			id           string
			rawActorName string
			rawSessionID string
			rawRequest   []byte
			createdAt    time.Time
		)
		if err := rows.Scan(&id, &rawActorName, &rawSessionID, &rawRequest, &createdAt); err != nil {
			return nil, err
		}

		var snapshot requestSnapshotRow
		if err := json.Unmarshal(rawRequest, &snapshot); err != nil {
			return nil, err
		}

		action := domain.RehydrateAction(
			id,
			rawActorName,
			rawSessionID,
			domain.RehydrateRequestSnapshot(
				snapshot.Method,
				snapshot.Path,
				snapshot.Query,
				snapshot.Body,
				snapshot.Status,
			),
			createdAt,
		)
		actions = append(actions, action)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return actions, nil
}
