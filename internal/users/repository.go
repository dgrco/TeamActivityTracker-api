package users

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	GetById(ctx context.Context, id string) (*User, error)
	Insert(ctx context.Context, email string, username string, password string) error
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
	err := pr.db.QueryRow(ctx, "SELECT email, username FROM users WHERE id=$1", id).Scan(
		&user.Email, &user.Username,
	)

	if err == nil {
		user.ID = id
	}

	return user, err
}

// Insert a new user to the database
func (pr *PostgresRepository) Insert(ctx context.Context, email string, username string, passwordHash string) error {
	_, err := pr.db.Exec(
		ctx,
		`INSERT INTO users (email, username, password_hash) VALUES ($1, $2, $3)`,
		email,
		username,
		passwordHash,
	)

	return err
}
