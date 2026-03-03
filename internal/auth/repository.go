package auth

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	InsertRefreshTokenHash(ctx context.Context, userID string, tokenHash string, expiresAt time.Time) error
}

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// Insert a refresh token into the `tokens` table.
// Returns nil on success, and an error otherwise.
func (pr *PostgresRepository) InsertRefreshTokenHash(ctx context.Context, userID string, tokenHash string, expiresAt time.Time) error {
	_, err := pr.db.Exec(
		ctx,
		`INSERT INTO tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`,
		userID,
		tokenHash,
		expiresAt,
	)

	return err
}
