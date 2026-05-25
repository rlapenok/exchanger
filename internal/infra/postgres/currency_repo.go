package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	domain "github.com/rlapenok/exchanger/internal/domain/currency"
	"github.com/rlapenok/exchanger/internal/uc"
)

// CurrencyRepo is the repository for the currency
type CurrencyRepo struct {
	db *sql.DB
}

// NewCurrencyRepo creates a new CurrencyRepo
func NewCurrencyRepo(db *sql.DB) *CurrencyRepo {
	return &CurrencyRepo{db: db}
}

// GetAllCurrencies gets all currencies from the database
func (r *CurrencyRepo) GetAll(ctx context.Context, pagination uc.Pagination) ([]domain.Currency, error) {
	query := `
		SELECT *
		FROM currencies
		ORDER BY code ASC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	currencies := make([]domain.Currency, 0, pagination.Limit)
	for rows.Next() {
		var (
			rawCode      string
			rawName      string
			rawSymbol    string
			rawMinorUnit int16
		)
		if err := rows.Scan(&rawCode, &rawName, &rawSymbol, &rawMinorUnit); err != nil {
			return nil, err
		}

		currency := domain.RehydrateCurrency(
			rawCode,
			rawName,
			rawSymbol,
			rawMinorUnit,
		)
		currencies = append(currencies, currency)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return currencies, nil
}

// GetByCode gets a currency by code from the database
func (r *CurrencyRepo) GetByCode(ctx context.Context, code domain.Code) (domain.Currency, error) {
	var (
		rawCode      string
		rawName      string
		rawSymbol    string
		rawMinorUnit int16
	)

	query := `
		SELECT code, name, symbol, minor_unit
		FROM currencies
		WHERE code = $1
	`
	row := r.db.QueryRowContext(ctx, query, code.Value())
	if err := row.Scan(&rawCode, &rawName, &rawSymbol, &rawMinorUnit); err != nil {
		return domain.Currency{}, mapCurrencyRepoError(err)
	}

	return domain.RehydrateCurrency(
		rawCode,
		rawName,
		rawSymbol,
		rawMinorUnit,
	), nil
}

// Create creates a currency in the database
func (r *CurrencyRepo) Create(ctx context.Context, currency domain.Currency) error {
	query := `
		INSERT INTO currencies (code, name, symbol, minor_unit)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.ExecContext(
		ctx, query,
		currency.Code().Value(),
		currency.Name().Value(),
		currency.Symbol().Value(),
		currency.MinorUnit().Value(),
	)
	if err != nil {
		return mapCurrencyRepoError(err)
	}

	return nil
}

// Update updates a currency in the database
func (r *CurrencyRepo) Update(ctx context.Context, currency domain.Currency) error {
	query := `
		UPDATE currencies
		SET name = $1, symbol = $2, minor_unit = $3
		WHERE code = $4
	`
	result, err := r.db.ExecContext(
		ctx, query,
		currency.Name().Value(),
		currency.Symbol().Value(),
		currency.MinorUnit().Value(),
		currency.Code().Value(),
	)
	if err != nil {
		return mapCurrencyRepoError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// Delete deletes a currency from the database
func (r *CurrencyRepo) Delete(ctx context.Context, code domain.Code) error {
	query := `
		DELETE FROM currencies
		WHERE code = $1
	`
	result, err := r.db.ExecContext(ctx, query, code.Value())
	if err != nil {
		return mapCurrencyRepoError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func mapCurrencyRepoError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return domain.ErrAlreadyExists
		case "23503":
			return domain.ErrInUse
		}
	}

	return err
}
