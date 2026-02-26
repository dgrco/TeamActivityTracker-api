package users

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	GetById(ctx context.Context, id string) (*User, error)
}

// Concrete PostgresSQL Repository
type PostgresRepository struct {
	db *pgxpool.Pool
}

// Construct a new user Repository
func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db}
}

// Get a user by their ID
func (pr *PostgresRepository) GetById(ctx context.Context, id string) (*User, error) {
	user := &User{}
	err := pr.db.QueryRow(ctx, "SELECT id, email, username FROM users WHERE id=$1", id).Scan(
		&user.id, &user.email, &user.username,
	)

	return user, err
}
