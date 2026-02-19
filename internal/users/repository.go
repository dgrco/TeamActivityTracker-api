package users

import "github.com/jackc/pgx/v5/pgxpool"

type Repository struct {
	db *pgxpool.Pool
}

// Construct a new user Repository
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db}
}
