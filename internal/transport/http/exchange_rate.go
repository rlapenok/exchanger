package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	domCurrency "github.com/rlapenok/exchanger/internal/domain/currency"
	domExchange "github.com/rlapenok/exchanger/internal/domain/exchange"
	domain "github.com/rlapenok/exchanger/internal/domain/exchangerate"
	"github.com/rlapenok/exchanger/internal/uc"
	"github.com/rlapenok/exchanger/internal/uc/exchange"
	"github.com/rlapenok/exchanger/internal/uc/exchangerate"
)

type exchangeRateIDParam struct {
	ID string `uri:"id"`
}

type exchangeRateRequest struct {
	BaseCode     string `json:"base_currency_code"`
	QuoteCode    string `json:"quote_currency_code"`
	BuyRate      string `json:"buy_rate"`
	SellRate     string `json:"sell_rate"`
	IsBuyActive  bool   `json:"is_buy_active"`
	IsSellActive bool   `json:"is_sell_active"`
}

type updateExchangeRateRequest struct {
	BuyRate      string `json:"buy_rate"`
	SellRate     string `json:"sell_rate"`
	IsBuyActive  bool   `json:"is_buy_active"`
	IsSellActive bool   `json:"is_sell_active"`
}

type exchangeRateResponse struct {
	ID              string `json:"id"`
	BaseCurrencyCode string `json:"base_currency_code"`
	QuoteCurrencyCode string `json:"quote_currency_code"`
	BuyRate         string `json:"buy_rate"`
	SellRate        string `json:"sell_rate"`
	IsBuyActive     bool   `json:"is_buy_active"`
	IsSellActive    bool   `json:"is_sell_active"`
	UpdatedAt       string `json:"updated_at"`
}

type exchangeRateHistoryResponse struct {
	ExchangeRateID    string `json:"exchange_rate_id"`
	BaseCurrencyCode  string `json:"base_currency_code"`
	QuoteCurrencyCode string `json:"quote_currency_code"`
	BuyRate           string `json:"buy_rate"`
	SellRate          string `json:"sell_rate"`
	ValidFrom         string `json:"valid_from"`
	ValidTo           string `json:"valid_to"`
}

type exchangeRateReportQuery struct {
	PaginationQuery
	From      string `form:"from" binding:"required"`
	To        string `form:"to" binding:"required"`
	BaseCode  string `form:"base"`
	QuoteCode string `form:"quote"`
}

type getExchangeRatesResponse []exchangeRateResponse
type getExchangeRateHistoryResponse []exchangeRateHistoryResponse
type getExchangeReportResponse []exchangeResponse

// ExchangeRateHandler handles exchange rate endpoints.
type ExchangeRateHandler struct {
	listLiveUseCase       *exchangerate.ListLiveUseCase
	createUseCase         *exchangerate.CreateUseCase
	updateUseCase         *exchangerate.UpdateUseCase
	listHistoryUseCase    *exchangerate.ListHistoryUseCase
	listReportUseCase     *exchangerate.ListReportUseCase
	listExchangeReportUseCase *exchange.ListReportUseCase
}

// NewExchangeRateHandler creates a new ExchangeRateHandler.
func NewExchangeRateHandler(
	listLiveUseCase *exchangerate.ListLiveUseCase,
	createUseCase *exchangerate.CreateUseCase,
	updateUseCase *exchangerate.UpdateUseCase,
	listHistoryUseCase *exchangerate.ListHistoryUseCase,
	listReportUseCase *exchangerate.ListReportUseCase,
	listExchangeReportUseCase *exchange.ListReportUseCase,
) *ExchangeRateHandler {
	return &ExchangeRateHandler{
		listLiveUseCase:           listLiveUseCase,
		createUseCase:             createUseCase,
		updateUseCase:             updateUseCase,
		listHistoryUseCase:        listHistoryUseCase,
		listReportUseCase:         listReportUseCase,
		listExchangeReportUseCase: listExchangeReportUseCase,
	}
}

// ListLive handles GET /v1/exchange-rate.
func (h *ExchangeRateHandler) ListLive(c *gin.Context) {
	rates, err := h.listLiveUseCase.Execute(c.Request.Context())
	if err != nil {
		writeExchangeRateError(c, err)
		return
	}

	c.JSON(http.StatusOK, toGetExchangeRatesResponse(rates))
}

// Create handles POST /v1/exchange-rate.
func (h *ExchangeRateHandler) Create(c *gin.Context) {
	var req exchangeRateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.createUseCase.Execute(c.Request.Context(), exchangerate.CreateInput{
		BaseCode:     req.BaseCode,
		QuoteCode:    req.QuoteCode,
		BuyRate:      req.BuyRate,
		SellRate:     req.SellRate,
		IsBuyActive:  req.IsBuyActive,
		IsSellActive: req.IsSellActive,
	})
	if err != nil {
		writeExchangeRateError(c, err)
		return
	}

	c.JSON(http.StatusCreated, toExchangeRateResponse(result))
}

// Update handles PUT /v1/exchange-rate/:id.
func (h *ExchangeRateHandler) Update(c *gin.Context) {
	var params exchangeRateIDParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req updateExchangeRateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.updateUseCase.Execute(c.Request.Context(), exchangerate.UpdateInput{
		ID:           params.ID,
		BuyRate:      req.BuyRate,
		SellRate:     req.SellRate,
		IsBuyActive:  req.IsBuyActive,
		IsSellActive: req.IsSellActive,
	})
	if err != nil {
		writeExchangeRateError(c, err)
		return
	}

	c.JSON(http.StatusOK, toExchangeRateResponse(result))
}

// ListHistory handles GET /v1/exchange-rate/:id/history.
func (h *ExchangeRateHandler) ListHistory(c *gin.Context) {
	var params exchangeRateIDParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var query PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	history, err := h.listHistoryUseCase.Execute(c.Request.Context(), exchangerate.ListHistoryInput{
		ID:         params.ID,
		Pagination: uc.NewPagination(query.Limit, query.Offset),
	})
	if err != nil {
		writeExchangeRateError(c, err)
		return
	}

	c.JSON(http.StatusOK, toGetExchangeRateHistoryResponse(history))
}

// ListRatesByDate handles GET /v1/exchange-rate/rates.
func (h *ExchangeRateHandler) ListRatesByDate(c *gin.Context) {
	var query exchangeRateReportQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	report, err := h.listReportUseCase.Execute(c.Request.Context(), exchangerate.ListReportInput{
		From:       query.From,
		To:         query.To,
		BaseCode:   query.BaseCode,
		QuoteCode:  query.QuoteCode,
		Pagination: uc.NewPagination(query.Limit, query.Offset),
	})
	if err != nil {
		writeExchangeRateError(c, err)
		return
	}

	c.JSON(http.StatusOK, toGetExchangeRateHistoryResponse(report))
}

// ListReport handles GET /v1/exchange-rate/report.
func (h *ExchangeRateHandler) ListReport(c *gin.Context) {
	var query exchangeRateReportQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pagination := uc.NewPagination(query.Limit, query.Offset)

	exchanges, err := h.listExchangeReportUseCase.Execute(c.Request.Context(), exchange.ListReportInput{
		From:       query.From,
		To:         query.To,
		BaseCode:   query.BaseCode,
		QuoteCode:  query.QuoteCode,
		Pagination: pagination,
	})
	if err != nil {
		writeExchangeReportError(c, err)
		return
	}

	c.JSON(http.StatusOK, getExchangeReportResponse(toExchangeResponses(exchanges)))
}

func writeExchangeReportError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domExchange.ErrInvalidDateRange),
		errors.Is(err, domCurrency.ErrInvalidCode):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func writeExchangeRateError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidID),
		errors.Is(err, domain.ErrInvalidRate),
		errors.Is(err, domain.ErrSellLessThanBuy),
		errors.Is(err, domain.ErrSameCurrency),
		errors.Is(err, domain.ErrInvalidDateRange):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrAlreadyExists):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func toGetExchangeRatesResponse(rates []domain.ExchangeRate) getExchangeRatesResponse {
	response := make(getExchangeRatesResponse, len(rates))
	for i, rate := range rates {
		response[i] = toExchangeRateResponse(rate)
	}
	return response
}

func toExchangeRateResponse(rate domain.ExchangeRate) exchangeRateResponse {
	return exchangeRateResponse{
		ID:               rate.ID().Value(),
		BaseCurrencyCode: rate.BaseCode().Value(),
		QuoteCurrencyCode: rate.QuoteCode().Value(),
		BuyRate:          rate.BuyRate().Value(),
		SellRate:         rate.SellRate().Value(),
		IsBuyActive:      rate.IsBuyActive(),
		IsSellActive:     rate.IsSellActive(),
		UpdatedAt:        rate.UpdatedAt().Format(timeRFC3339),
	}
}

func toGetExchangeRateHistoryResponse(history []domain.History) getExchangeRateHistoryResponse {
	response := make(getExchangeRateHistoryResponse, len(history))
	for i, item := range history {
		response[i] = exchangeRateHistoryResponse{
			ExchangeRateID:    item.ExchangeRateID().Value(),
			BaseCurrencyCode:  item.BaseCode().Value(),
			QuoteCurrencyCode: item.QuoteCode().Value(),
			BuyRate:           item.BuyRate().Value(),
			SellRate:          item.SellRate().Value(),
			ValidFrom:         item.ValidFrom().Format(timeRFC3339),
			ValidTo:           item.ValidTo().Format(timeRFC3339),
		}
	}
	return response
}

const timeRFC3339 = "2006-01-02T15:04:05Z07:00"
