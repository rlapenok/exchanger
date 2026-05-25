package action

import (
	"time"

	domUser "github.com/rlapenok/exchanger/internal/domain/user"
)

// Action is a persisted user activity record.
type Action struct {
	id        string
	actorName domUser.Name
	sessionID SessionID
	request   RequestSnapshot
	createdAt time.Time
}

// NewAction creates a new Action ready for persistence.
func NewAction(
	actorName domUser.Name,
	sessionID SessionID,
	request RequestSnapshot,
) Action {
	return Action{
		actorName: actorName,
		sessionID: sessionID,
		request:   request,
	}
}

// ID returns the action identifier.
func (a Action) ID() string {
	return a.id
}

// ActorName returns the actor name.
func (a Action) ActorName() domUser.Name {
	return a.actorName
}

// SessionID returns the session identifier.
func (a Action) SessionID() SessionID {
	return a.sessionID
}

// Request returns the HTTP request snapshot.
func (a Action) Request() RequestSnapshot {
	return a.request
}

// CreatedAt returns the action timestamp.
func (a Action) CreatedAt() time.Time {
	return a.createdAt
}

// RehydrateAction restores an Action from storage.
func RehydrateAction(
	id string,
	actorName string,
	sessionID string,
	request RequestSnapshot,
	createdAt time.Time,
) Action {
	return Action{
		id:        id,
		actorName: domUser.RehydrateName(actorName),
		sessionID: RehydrateSessionID(sessionID),
		request:   request,
		createdAt: createdAt,
	}
}

// WithID returns a copy of the action with the given id.
func (a Action) WithID(id string) Action {
	a.id = id
	return a
}

// WithCreatedAt returns a copy of the action with the given timestamp.
func (a Action) WithCreatedAt(createdAt time.Time) Action {
	a.createdAt = createdAt
	return a
}
