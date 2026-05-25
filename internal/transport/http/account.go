package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	ucaction "github.com/rlapenok/exchanger/internal/uc/action"
	"github.com/rlapenok/exchanger/internal/uc"
	"github.com/rlapenok/exchanger/internal/uc/user"
)

// AccountResponse is the response for the account endpoint
type accountResponse struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

type sessionActionResponse struct {
	Request   sessionActionRequestResponse `json:"request"`
	CreatedAt time.Time                    `json:"created_at"`
}

type sessionActionRequestResponse struct {
	Method string          `json:"method"`
	Path   string          `json:"path"`
	Query  string          `json:"query,omitempty"`
	Body   json.RawMessage `json:"body,omitempty"`
	Status int             `json:"status"`
}

type sessionActionsResponse []sessionActionResponse

// AccountHandler is the handler for the account endpoint
type AccountHandler struct {
	accountUC      *user.AccountUseCase
	listSessionUC  *ucaction.ListSessionUseCase
}

// NewAccountHandler creates a new AccountHandler
func NewAccountHandler(
	accountUC *user.AccountUseCase,
	listSessionUC *ucaction.ListSessionUseCase,
) *AccountHandler {
	return &AccountHandler{
		accountUC:     accountUC,
		listSessionUC: listSessionUC,
	}
}

// Account returns the current user account from session cookie.
func (h *AccountHandler) Account(c *gin.Context) {
	name, _ := c.Get("name")
	nameStr, ok := name.(string)
	if !ok || nameStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	out, err := h.accountUC.Execute(c.Request.Context(), user.ProfileInput{
		Name: nameStr,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, accountResponse{
		Name: out.Name,
		Role: out.Role,
	})
}

// SessionActions returns actions for the current user session.
func (h *AccountHandler) SessionActions(c *gin.Context) {
	name, _ := c.Get("name")
	nameStr, ok := name.(string)
	if !ok || nameStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	session := sessions.Default(c)
	sessionID := session.ID()
	if sessionID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var query PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	actions, err := h.listSessionUC.Execute(c.Request.Context(), ucaction.ListSessionInput{
		ActorName:  nameStr,
		SessionID:  sessionID,
		Pagination: uc.NewPagination(query.Limit, query.Offset),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make(sessionActionsResponse, len(actions))
	for i, action := range actions {
		response[i] = sessionActionResponse{
			Request: sessionActionRequestResponse{
				Method: action.Request.Method,
				Path:   action.Request.Path,
				Query:  action.Request.Query,
				Body:   action.Request.Body,
				Status: action.Request.Status,
			},
			CreatedAt: action.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}
