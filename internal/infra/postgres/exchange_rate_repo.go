package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	domain "github.com/rlapenok/exchanger/internal/domain/exchangerate"
	"github.com/rlapenok/exchanger/internal/uc"
	ucexrate "github.com/rlapenok/exchanger/internal/uc/exchangerate"
)

// ExchangeRateRepo is the repository for exchange rates.
type ExchangeRateRepo struct {
	db *sql.DB
}

// NewExchangeRateRepo creates a new ExchangeRateRepo.
func NewExchangeRateRepo(db *sql.DB) *ExchangeRateRepo {
	return &ExchangeRateRepo{db: db}
}

// ListLive returns current exchange rates.
func (r *ExchangeRateRepo) ListLive(ctx context.Context) ([]domain.ExchangeRate, error) {
	query := `
		SELECT id, base_currency_code, quote_currency_code,
		       buy_rate::text, sell_rate::text,
		       is_buy_active, is_sell_active, updated_at
		FROM exchange_rates
		ORDER BY base_currency_code ASC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanExchangeRates(rows)
}

// GetByID returns an exchange rate by id.
func (r *ExchangeRateRepo) GetByID(ctx context.Context, id domain.ID) (domain.ExchangeRate, error) {
	query := `
		SELECT id, base_currency_code, quote_currency_code,
		       buy_rate::text, sell_rate::text,
		       is_buy_active, is_sell_active, updated_at
		FROM exchange_rates
		WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id.Value())
	rate, err := scanExchangeRateRow(row)
	if err != nil {
		return domain.ExchangeRate{}, mapExchangeRateRepoError(err)
	}

	return rate, nil
}

// GetByPair returns an exchange rate by currency pair.
func (r *ExchangeRateRepo) GetByPair(
	ctx context.Context,
	baseCode string,
	quoteCode string,
) (domain.ExchangeRate, error) {
	query := `
		SELECT id, base_currency_code, quote_currency_code,
		       buy_rate::text, sell_rate::text,
		       is_buy_active, is_sell_active, updated_at
		FROM exchange_rates
		WHERE base_currency_code = $1 AND quote_currency_code = $2
	`

	row := r.db.QueryRowContext(ctx, query, baseCode, quoteCode)
	rate, err := scanExchangeRateRow(row)
	if err != nil {
		return domain.ExchangeRate{}, mapExchangeRateRepoError(err)
	}

	return rate, nil
}

func (r *ExchangeRateRepo) getByIDForUpdate(
	ctx context.Context,
	tx *sql.Tx,
	id domain.ID,
) (domain.ExchangeRate, error) {
	query := `
		SELECT id, base_currency_code, quote_currency_code,
		       buy_rate::text, sell_rate::text,
		       is_buy_active, is_sell_active, updated_at
		FROM exchange_rates
		WHERE id = $1
		FOR UPDATE
	`

	row := tx.QueryRowContext(ctx, query, id.Value())
	rate, err := scanExchangeRateRow(row)
	if err != nil {
		return domain.ExchangeRate{}, mapExchangeRateRepoError(err)
	}

	return rate, nil
}

// Create inserts a new exchange rate and initial history snapshot.
func (r *ExchangeRateRepo) Create(ctx context.Context, rate domain.ExchangeRate) (domain.ExchangeRate, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.ExchangeRate{}, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	insertQuery := `
		INSERT INTO exchange_rates (
			base_currency_code, quote_currency_code,
			buy_rate, sell_rate, is_buy_active, is_sell_active
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, updated_at
	`

	var (
		id        string
		updatedAt time.Time
	)
	err = tx.QueryRowContext(
		ctx,
		insertQuery,
		rate.BaseCode().Value(),
		rate.QuoteCode().Value(),
		rate.BuyRate().Value(),
		rate.SellRate().Value(),
		rate.IsBuyActive(),
		rate.IsSellActive(),
	).Scan(&id, &updatedAt)
	if err != nil {
		return domain.ExchangeRate{}, mapExchangeRateRepoError(err)
	}

	historyQuery := `
		INSERT INTO exchange_rates_history (
			exchange_rate_id, buy_rate, sell_rate, valid_from, valid_to
		)
		VALUES ($1, $2, $3, $4, $4)
	`
	_, err = tx.ExecContext(
		ctx,
		historyQuery,
		id,
		rate.BuyRate().Value(),
		rate.SellRate().Value(),
		updatedAt,
	)
	if err != nil {
		return domain.ExchangeRate{}, err
	}

	if err := tx.Commit(); err != nil {
		return domain.ExchangeRate{}, err
	}

	return domain.RehydrateExchangeRate(
		id,
		rate.BaseCode().Value(),
		rate.QuoteCode().Value(),
		rate.BuyRate().Value(),
		rate.SellRate().Value(),
		rate.IsBuyActive(),
		rate.IsSellActive(),
		updatedAt,
	), nil
}

// Update updates an exchange rate and appends a history snapshot.
func (r *ExchangeRateRepo) Update(ctx context.Context, rate domain.ExchangeRate) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	current, err := r.getByIDForUpdate(ctx, tx, rate.ID())
	if err != nil {
		return err
	}

	now := time.Now()

	historyQuery := `
		INSERT INTO exchange_rates_history (
			exchange_rate_id, buy_rate, sell_rate, valid_from, valid_to
		)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = tx.ExecContext(
		ctx,
		historyQuery,
		current.ID().Value(),
		current.BuyRate().Value(),
		current.SellRate().Value(),
		now,
		now,
	)
	if err != nil {
		return err
	}

	updateQuery := `
		UPDATE exchange_rates
		SET buy_rate = $1,
		    sell_rate = $2,
		    is_buy_active = $3,
		    is_sell_active = $4,
		    updated_at = $5
		WHERE id = $6
	`
	result, err := tx.ExecContext(
		ctx,
		updateQuery,
		rate.BuyRate().Value(),
		rate.SellRate().Value(),
		rate.IsBuyActive(),
		rate.IsSellActive(),
		now,
		rate.ID().Value(),
	)
	if err != nil {
		return mapExchangeRateRepoError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return tx.Commit()
}

// ListHistory returns history rows for one exchange rate.
func (r *ExchangeRateRepo) ListHistory(
	ctx context.Context,
	id domain.ID,
	pagination uc.Pagination,
) ([]domain.History, error) {
	query := `
		SELECT h.exchange_rate_id,
		       er.base_currency_code,
		       er.quote_currency_code,
		       h.buy_rate::text,
		       h.sell_rate::text,
		       h.valid_from,
		       h.valid_to
		FROM exchange_rates_history h
		INNER JOIN exchange_rates er ON er.id = h.exchange_rate_id
		WHERE h.exchange_rate_id = $1
		ORDER BY h.valid_from DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, id.Value(), pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanHistories(rows)
}

// ListReport returns history rows filtered by date range and optional pair.
func (r *ExchangeRateRepo) ListReport(
	ctx context.Context,
	filter ucexrate.ReportFilter,
	pagination uc.Pagination,
) ([]domain.History, error) {
	query := `
		SELECT h.exchange_rate_id,
		       er.base_currency_code,
		       er.quote_currency_code,
		       h.buy_rate::text,
		       h.sell_rate::text,
		       h.valid_from,
		       h.valid_to
		FROM exchange_rates_history h
		INNER JOIN exchange_rates er ON er.id = h.exchange_rate_id
		WHERE h.valid_from >= $1
		  AND h.valid_from < $2
	`
	args := []any{filter.From, filter.To.Add(24 * time.Hour)}
	argIndex := 3

	if filter.BaseCode != "" {
		query += ` AND er.base_currency_code = $` + strconv.Itoa(argIndex)
		args = append(args, filter.BaseCode)
		argIndex++
	}

	if filter.QuoteCode != "" {
		query += ` AND er.quote_currency_code = $` + strconv.Itoa(argIndex)
		args = append(args, filter.QuoteCode)
		argIndex++
	}

	query += ` ORDER BY h.valid_from DESC LIMIT $` + strconv.Itoa(argIndex) + ` OFFSET $` + strconv.Itoa(argIndex+1)
	args = append(args, pagination.Limit, pagination.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanHistories(rows)
}

func scanExchangeRates(rows *sql.Rows) ([]domain.ExchangeRate, error) {
	rates := make([]domain.ExchangeRate, 0)
	for rows.Next() {
		rate, err := scanExchangeRateRows(rows)
		if err != nil {
			return nil, err
		}
		rates = append(rates, rate)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rates, nil
}

func scanExchangeRateRow(row *sql.Row) (domain.ExchangeRate, error) {
	var (
		id           string
		baseCode     string
		quoteCode    string
		buyRate      string
		sellRate     string
		isBuyActive  bool
		isSellActive bool
		updatedAt    time.Time
	)
	if err := row.Scan(
		&id,
		&baseCode,
		&quoteCode,
		&buyRate,
		&sellRate,
		&isBuyActive,
		&isSellActive,
		&updatedAt,
	); err != nil {
		return domain.ExchangeRate{}, err
	}

	return domain.RehydrateExchangeRate(
		id,
		baseCode,
		quoteCode,
		buyRate,
		sellRate,
		isBuyActive,
		isSellActive,
		updatedAt,
	), nil
}

func scanExchangeRateRows(rows *sql.Rows) (domain.ExchangeRate, error) {
	var (
		id           string
		baseCode     string
		quoteCode    string
		buyRate      string
		sellRate     string
		isBuyActive  bool
		isSellActive bool
		updatedAt    time.Time
	)
	if err := rows.Scan(
		&id,
		&baseCode,
		&quoteCode,
		&buyRate,
		&sellRate,
		&isBuyActive,
		&isSellActive,
		&updatedAt,
	); err != nil {
		return domain.ExchangeRate{}, err
	}

	return domain.RehydrateExchangeRate(
		id,
		baseCode,
		quoteCode,
		buyRate,
		sellRate,
		isBuyActive,
		isSellActive,
		updatedAt,
	), nil
}

func scanHistories(rows *sql.Rows) ([]domain.History, error) {
	histories := make([]domain.History, 0)
	for rows.Next() {
		var (
			exchangeRateID string
			baseCode       string
			quoteCode      string
			buyRate        string
			sellRate       string
			validFrom      time.Time
			validTo        time.Time
		)
		if err := rows.Scan(
			&exchangeRateID,
			&baseCode,
			&quoteCode,
			&buyRate,
			&sellRate,
			&validFrom,
			&validTo,
		); err != nil {
			return nil, err
		}

		histories = append(histories, domain.RehydrateHistory(
			exchangeRateID,
			baseCode,
			quoteCode,
			buyRate,
			sellRate,
			validFrom,
			validTo,
		))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return histories, nil
}

func mapExchangeRateRepoError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return domain.ErrAlreadyExists
	}

	return err
}
