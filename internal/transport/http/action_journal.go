package http

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	ucaction "github.com/rlapenok/exchanger/internal/uc/action"
)

const (
	maxActionBodySize  = 64 * 1024
	sessionActionsPath = "/v1/account/actions"
	loginPath          = "/v1/auth/login"
)

// ActionJournal records API requests after handlers finish.
func ActionJournal(recordUC *ucaction.RecordUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		if shouldSkipActionJournal(c.Request.URL.Path) {
			c.Next()
			return
		}

		session := sessions.Default(c)
		preActorName := resolveActorName(c, session)
		preSessionID := session.ID()

		bodyBytes := captureRequestBody(c)

		c.Next()

		actorName := resolveActorName(c, session)
		if actorName == "" {
			actorName = preActorName
		}
		if actorName == "" && c.Request.URL.Path == loginPath {
			actorName = parseLoginActorName(bodyBytes)
		}
		if actorName == "" {
			return
		}

		sessionID := session.ID()
		if sessionID == "" {
			sessionID = preSessionID
		}
		if sessionID == "" {
			return
		}

		body := sanitizeRequestBody(c.Request.URL.Path, bodyBytes)
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		err := recordUC.Execute(c.Request.Context(), ucaction.RecordInput{
			ActorName: actorName,
			SessionID: sessionID,
			Method:    c.Request.Method,
			Path:      path,
			Query:     c.Request.URL.RawQuery,
			Body:      body,
			Status:    c.Writer.Status(),
		})
		if err != nil {
			slog.WarnContext(
				c.Request.Context(),
				"failed to record user action",
				slog.String("path", path),
				slog.String("error", err.Error()),
			)
		}
	}
}

func shouldSkipActionJournal(path string) bool {
	return path == sessionActionsPath
}

func resolveActorName(c *gin.Context, session sessions.Session) string {
	if name := c.GetString("name"); name != "" {
		return name
	}

	value, ok := session.Get("name").(string)
	if !ok {
		return ""
	}

	return strings.TrimSpace(value)
}

func captureRequestBody(c *gin.Context) []byte {
	if c.Request.Body == nil {
		return nil
	}
	if !shouldCaptureBody(c.Request.Method) {
		return nil
	}

	limited := io.LimitReader(c.Request.Body, maxActionBodySize+1)
	bodyBytes, err := io.ReadAll(limited)
	if err != nil {
		return nil
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	if len(bodyBytes) > maxActionBodySize {
		return bodyBytes[:maxActionBodySize]
	}

	return bodyBytes
}

func shouldCaptureBody(method string) bool {
	switch strings.ToUpper(method) {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		return true
	default:
		return false
	}
}

func sanitizeRequestBody(path string, bodyBytes []byte) json.RawMessage {
	if len(bodyBytes) == 0 {
		return nil
	}

	if path != loginPath {
		if !json.Valid(bodyBytes) {
			return nil
		}
		return json.RawMessage(bodyBytes)
	}

	var payload map[string]json.RawMessage
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		return nil
	}

	delete(payload, "password")

	sanitized, err := json.Marshal(payload)
	if err != nil {
		return nil
	}

	return json.RawMessage(sanitized)
}

func parseLoginActorName(bodyBytes []byte) string {
	var payload struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		return ""
	}

	return strings.TrimSpace(payload.Name)
}
