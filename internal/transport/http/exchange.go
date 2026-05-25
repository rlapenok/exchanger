package http

import (
	"errors"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	domCurrency "github.com/rlapenok/exchanger/internal/domain/currency"
	domain "github.com/rlapenok/exchanger/internal/domain/exchange"
	"github.com/rlapenok/exchanger/internal/uc/exchange"
)

type exchangeRequest struct {
	BaseCurrencyCode  string `json:"base_currency_code"`
	QuoteCurrencyCode string `json:"quote_currency_code"`
	Side              string `json:"side"`
	Amount            string `json:"amount"`
}

type exchangeResponse struct {
	ID                string `json:"id"`
	OperatorName      string `json:"operator_name"`
	BaseCurrencyCode  string `json:"base_currency_code"`
	QuoteCurrencyCode string `json:"quote_currency_code"`
	Side              string `json:"side"`
	Amount            string `json:"amount"`
	Rate              string `json:"rate"`
	ResultAmount      string `json:"result_amount"`
	CreatedAt         string `json:"created_at"`
}

// ExchangeHandler handles currency exchange operations.
type ExchangeHandler struct {
	executeUseCase *exchange.ExecuteUseCase
}

// NewExchangeHandler creates a new ExchangeHandler.
func NewExchangeHandler(executeUseCase *exchange.ExecuteUseCase) *ExchangeHandler {
	return &ExchangeHandler{executeUseCase: executeUseCase}
}

// Execute handles POST /v1/exchange.
func (h *ExchangeHandler) Execute(c *gin.Context) {
	var req exchangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	name := c.GetString("name")
	if name == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	session := sessions.Default(c)
	sessionID := session.ID()
	if sessionID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	result, err := h.executeUseCase.Execute(c.Request.Context(), exchange.ExecuteInput{
		OperatorName: name,
		SessionID:    sessionID,
		BaseCode:     req.BaseCurrencyCode,
		QuoteCode:    req.QuoteCurrencyCode,
		Side:         req.Side,
		Amount:       req.Amount,
	})
	if err != nil {
		writeExchangeOperationError(c, err)
		return
	}

	c.JSON(http.StatusCreated, toExchangeResponse(result))
}

func writeExchangeOperationError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidSide),
		errors.Is(err, domain.ErrInvalidAmount),
		errors.Is(err, domCurrency.ErrInvalidCode):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrRateNotActive),
		errors.Is(err, domain.ErrRateNotFound):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func toExchangeResponses(items []domain.Exchange) []exchangeResponse {
	response := make([]exchangeResponse, len(items))
	for i, item := range items {
		response[i] = toExchangeResponse(item)
	}
	return response
}

func toExchangeResponse(item domain.Exchange) exchangeResponse {
	return exchangeResponse{
		ID:                item.ID(),
		OperatorName:      item.OperatorName(),
		BaseCurrencyCode:  item.BaseCode().Value(),
		QuoteCurrencyCode: item.QuoteCode().Value(),
		Side:              item.Side().Value(),
		Amount:            item.Amount().Value(),
		Rate:              item.Rate().Value(),
		ResultAmount:      item.ResultAmount().Value(),
		CreatedAt:         item.CreatedAt().Format(timeRFC3339),
	}
}
