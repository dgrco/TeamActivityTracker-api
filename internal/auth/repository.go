package auth

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	InsertRefreshTokenHash(ctx context.Context, userID string, tokenHash string, expiresAt time.Time) error
	GetTokenFromRefreshTokenHash(ctx context.Context, refreshTokenHash string) (string, error)
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

// Find a token entry given a hashed token string.
// Returns the associated `user_id` on success, and an error otherwise
func (pr *PostgresRepository) GetTokenFromRefreshTokenHash(ctx context.Context, refreshTokenHash string) (string, error) {
	var userID string
	err := pr.db.QueryRow(
		ctx,
		`SELECT user_id 
		FROM tokens
		WHERE token_hash=$1
			AND expires_at > NOW()
			AND revoked_at IS NULL
		LIMIT 1`,
		refreshTokenHash,
	).Scan(&userID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrTokenInvalidOrExpired
		}

		return "", err
	}

	return userID, nil
}
