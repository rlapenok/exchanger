package postgres

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	domain "github.com/rlapenok/exchanger/internal/domain/exchange"
	"github.com/rlapenok/exchanger/internal/uc"
	ucexchange "github.com/rlapenok/exchanger/internal/uc/exchange"
)

// ExchangeRepo persists currency exchanges.
type ExchangeRepo struct {
	db *sql.DB
}

// NewExchangeRepo creates a new ExchangeRepo.
func NewExchangeRepo(db *sql.DB) *ExchangeRepo {
	return &ExchangeRepo{db: db}
}

// Create inserts an exchange operation.
func (r *ExchangeRepo) Create(ctx context.Context, exchange domain.Exchange) (domain.Exchange, error) {
	query := `
		INSERT INTO exchanges (
			operator_name, session_id,
			base_currency_code, quote_currency_code,
			side, amount, rate, result_amount
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at
	`

	var (
		id        string
		createdAt time.Time
	)
	err := r.db.QueryRowContext(
		ctx,
		query,
		exchange.OperatorName(),
		exchange.SessionID(),
		exchange.BaseCode().Value(),
		exchange.QuoteCode().Value(),
		exchange.Side().Value(),
		exchange.Amount().Value(),
		exchange.Rate().Value(),
		exchange.ResultAmount().Value(),
	).Scan(&id, &createdAt)
	if err != nil {
		return domain.Exchange{}, err
	}

	return domain.RehydrateExchange(
		id,
		exchange.OperatorName(),
		exchange.SessionID(),
		exchange.BaseCode().Value(),
		exchange.QuoteCode().Value(),
		exchange.Side().Value(),
		exchange.Amount().Value(),
		exchange.Rate().Value(),
		exchange.ResultAmount().Value(),
		createdAt,
	), nil
}

// ListReport returns exchanges filtered by date range and optional pair.
func (r *ExchangeRepo) ListReport(
	ctx context.Context,
	filter ucexchange.ReportFilter,
	pagination uc.Pagination,
) ([]domain.Exchange, error) {
	query := `
		SELECT id, operator_name, session_id,
		       base_currency_code, quote_currency_code,
		       side, amount::text, rate::text, result_amount::text,
		       created_at
		FROM exchanges
		WHERE created_at >= $1
		  AND created_at < $2
	`
	args := []any{filter.From, filter.To.Add(24 * time.Hour)}
	argIndex := 3

	if filter.BaseCode != "" {
		query += ` AND base_currency_code = $` + strconv.Itoa(argIndex)
		args = append(args, filter.BaseCode)
		argIndex++
	}

	if filter.QuoteCode != "" {
		query += ` AND quote_currency_code = $` + strconv.Itoa(argIndex)
		args = append(args, filter.QuoteCode)
		argIndex++
	}

	query += ` ORDER BY created_at DESC LIMIT $` + strconv.Itoa(argIndex) + ` OFFSET $` + strconv.Itoa(argIndex+1)
	args = append(args, pagination.Limit, pagination.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	exchanges := make([]domain.Exchange, 0)
	for rows.Next() {
		var (
			id             string
			operatorName   string
			sessionID      string
			baseCode       string
			quoteCode      string
			side           string
			amount         string
			rate           string
			resultAmount   string
			createdAt      time.Time
		)
		if err := rows.Scan(
			&id,
			&operatorName,
			&sessionID,
			&baseCode,
			&quoteCode,
			&side,
			&amount,
			&rate,
			&resultAmount,
			&createdAt,
		); err != nil {
			return nil, err
		}

		exchanges = append(exchanges, domain.RehydrateExchange(
			id,
			operatorName,
			sessionID,
			baseCode,
			quoteCode,
			side,
			amount,
			rate,
			resultAmount,
			createdAt,
		))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return exchanges, nil
}
