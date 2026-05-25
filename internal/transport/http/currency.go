package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	domain "github.com/rlapenok/exchanger/internal/domain/currency"
	"github.com/rlapenok/exchanger/internal/uc"
	"github.com/rlapenok/exchanger/internal/uc/currency"
)

type currencyCodeParam struct {
	Code string `uri:"code"`
}

type currencyRequest struct {
	Code      string `json:"code"`
	Name      string `json:"name"`
	Symbol    string `json:"symbol"`
	MinorUnit uint8  `json:"minor_unit"`
}

type updateCurrencyRequest struct {
	Name      string `json:"name"`
	Symbol    string `json:"symbol"`
	MinorUnit uint8  `json:"minor_unit"`
}

type currencyResponse struct {
	Code      string `json:"code"`
	Name      string `json:"name"`
	Symbol    string `json:"symbol"`
	MinorUnit uint8  `json:"minor_unit"`
}

type getAllCurrenciesResponse []currencyResponse

// CurrencyHandler is the handler for the currency endpoints
type CurrencyHandler struct {
	getAllUseCase    *currency.GetAllCurrenciesUseCase
	getByCodeUseCase *currency.GetByCodeUseCase
	createUseCase    *currency.CreateUseCase
	updateUseCase    *currency.UpdateUseCase
	deleteUseCase    *currency.DeleteUseCase
}

// NewCurrencyHandler creates a new CurrencyHandler
func NewCurrencyHandler(
	getAllUseCase *currency.GetAllCurrenciesUseCase,
	getByCodeUseCase *currency.GetByCodeUseCase,
	createUseCase *currency.CreateUseCase,
	updateUseCase *currency.UpdateUseCase,
	deleteUseCase *currency.DeleteUseCase,
) *CurrencyHandler {
	return &CurrencyHandler{
		getAllUseCase:    getAllUseCase,
		getByCodeUseCase: getByCodeUseCase,
		createUseCase:    createUseCase,
		updateUseCase:    updateUseCase,
		deleteUseCase:    deleteUseCase,
	}
}

// GetAllCurrencies handles the GET request for all currencies
func (h *CurrencyHandler) GetAllCurrencies(c *gin.Context) {
	var query PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pagination := uc.NewPagination(query.Limit, query.Offset)

	currencies, err := h.getAllUseCase.Execute(c.Request.Context(), pagination)
	if err != nil {
		writeCurrencyError(c, err)
		return
	}

	c.JSON(http.StatusOK, toGetAllCurrenciesResponse(currencies))
}

// GetCurrencyByCode handles the GET request for a currency by code
func (h *CurrencyHandler) GetCurrencyByCode(c *gin.Context) {
	var params currencyCodeParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.getByCodeUseCase.Execute(c.Request.Context(), currency.GetByCodeInput{
		Code: params.Code,
	})
	if err != nil {
		writeCurrencyError(c, err)
		return
	}

	c.JSON(http.StatusOK, toCurrencyResponse(result))
}

// CreateCurrency handles the POST request for creating a currency
func (h *CurrencyHandler) CreateCurrency(c *gin.Context) {
	var req currencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.createUseCase.Execute(c.Request.Context(), currency.CreateInput{
		Code:      req.Code,
		Name:      req.Name,
		Symbol:    req.Symbol,
		MinorUnit: req.MinorUnit,
	})
	if err != nil {
		writeCurrencyError(c, err)
		return
	}

	c.JSON(http.StatusCreated, toCurrencyResponse(result))
}

// UpdateCurrency handles the PUT request for updating a currency
func (h *CurrencyHandler) UpdateCurrency(c *gin.Context) {
	var params currencyCodeParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req updateCurrencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.updateUseCase.Execute(c.Request.Context(), currency.UpdateInput{
		Code:      params.Code,
		Name:      req.Name,
		Symbol:    req.Symbol,
		MinorUnit: req.MinorUnit,
	})
	if err != nil {
		writeCurrencyError(c, err)
		return
	}

	c.JSON(http.StatusOK, toCurrencyResponse(result))
}

// DeleteCurrency handles the DELETE request for deleting a currency
func (h *CurrencyHandler) DeleteCurrency(c *gin.Context) {
	var params currencyCodeParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.deleteUseCase.Execute(c.Request.Context(), currency.DeleteInput{
		Code: params.Code,
	}); err != nil {
		writeCurrencyError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func writeCurrencyError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidCode),
		errors.Is(err, domain.ErrInvalidName),
		errors.Is(err, domain.ErrInvalidSymbol),
		errors.Is(err, domain.ErrInvalidMinorUnit):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrAlreadyExists),
		errors.Is(err, domain.ErrInUse):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

// toGetAllCurrenciesResponse converts a list of currencies to a JSON array response
func toGetAllCurrenciesResponse(currencies []domain.Currency) getAllCurrenciesResponse {
	response := make(getAllCurrenciesResponse, len(currencies))
	for i, currency := range currencies {
		response[i] = toCurrencyResponse(currency)
	}
	return response
}

// toCurrencyResponse converts a currency to a currencyResponse
func toCurrencyResponse(currency domain.Currency) currencyResponse {
	return currencyResponse{
		Code:      currency.Code().Value(),
		Name:      currency.Name().Value(),
		Symbol:    currency.Symbol().Value(),
		MinorUnit: currency.MinorUnit().Value(),
	}
}
