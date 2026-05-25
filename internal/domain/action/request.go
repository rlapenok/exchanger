package action

import (
	"encoding/json"
	"net/http"
	"strings"
)

// RequestSnapshot stores HTTP request metadata as a single JSON document.
type RequestSnapshot struct {
	Method string          `json:"method"`
	Path   string          `json:"path"`
	Query  string          `json:"query,omitempty"`
	Body   json.RawMessage `json:"body,omitempty"`
	Status int             `json:"status"`
}

// NewRequestSnapshot validates and creates a request snapshot.
func NewRequestSnapshot(
	method string,
	path string,
	query string,
	body json.RawMessage,
	status int,
) (RequestSnapshot, error) {
	method = strings.ToUpper(strings.TrimSpace(method))
	if !isAllowedMethod(method) {
		return RequestSnapshot{}, ErrInvalidMethod
	}

	path = strings.TrimSpace(path)
	if path == "" {
		return RequestSnapshot{}, ErrInvalidPath
	}

	if status < http.StatusContinue || status > 599 {
		return RequestSnapshot{}, ErrInvalidStatus
	}

	if len(body) == 0 {
		body = nil
	}

	return RequestSnapshot{
		Method: method,
		Path:   path,
		Query:  strings.TrimSpace(query),
		Body:   body,
		Status: status,
	}, nil
}

// RehydrateRequestSnapshot restores a snapshot from storage.
func RehydrateRequestSnapshot(
	method string,
	path string,
	query string,
	body json.RawMessage,
	status int,
) RequestSnapshot {
	return RequestSnapshot{
		Method: method,
		Path:   path,
		Query:  query,
		Body:   body,
		Status: status,
	}
}

func isAllowedMethod(method string) bool {
	switch method {
	case http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodHead,
		http.MethodOptions:
		return true
	default:
		return false
	}
}
