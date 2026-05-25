package postgres

import (
	"context"
	"database/sql"

	"github.com/rlapenok/exchanger/internal/domain/user"
)

// UserRepo is the authentication repository
type UserRepo struct {
	db *sql.DB
}

// NewUserRepo creates a new user repository
func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

// GetByName gets a user by name
func (r *UserRepo) GetByName(ctx context.Context, name user.Name) (user.User, error) {
	var (
		rawName           string
		rawHashedPassword string
		rawRole           string
	)

	const query = `
	SELECT name,hashed_password,role 
	FROM users 
	WHERE name = $1`

	row := r.db.QueryRowContext(
		ctx,
		query,
		name.Value(),
	)
	if err := row.Scan(&rawName, &rawHashedPassword, &rawRole); err != nil {
		return user.User{}, err
	}

	return user.RehydrateUser(rawName, rawHashedPassword, rawRole), nil
}
