package users

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

// Construct a new user Repository
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db}
}

// Get a user by their ID
func (r *Repository) GetById(ctx context.Context, id string) (*User, error) {
	user := &User{}
	err := r.db.QueryRow(ctx, "SELECT id, email, username FROM users WHERE id=$1", id).Scan(
		&user.id, &user.email, &user.username,
	)

	return user, err
}
