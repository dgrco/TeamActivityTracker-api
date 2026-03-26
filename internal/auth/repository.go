package auth

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	// Insert a refresh token into the `tokens` table.
	// Returns nil on success, and an error otherwise.
	InsertRefreshTokenHash(ctx context.Context, userID string, tokenHash string, expiresAt time.Time) error
	// Find a token entry given a hashed token string.
	// Returns the associated `user_id` on success, and an error otherwise.
	GetUserIDFromRefreshTokenHash(ctx context.Context, refreshTokenHash string) (string, error)
	// Revoke a refresh token such that it is no longer valid.
	RevokeToken(ctx context.Context, refreshTokenHash string) error
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
func (pr *PostgresRepository) GetUserIDFromRefreshTokenHash(ctx context.Context, refreshTokenHash string) (string, error) {
	var userID string
	err := pr.db.QueryRow(
		ctx,
		`SELECT user_id
		FROM tokens
		WHERE token_hash = $1
			AND expires_at > NOW()
			AND revoked_at IS NULL`,
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

// Revoke a refresh token such that it is no longer valid.
// This fails if there is either no matching token entry or if the token has already been revoked.
func (pr *PostgresRepository) RevokeToken(ctx context.Context, refreshTokenHash string) error {
	var revokedAt *time.Time
	err := pr.db.QueryRow(
		ctx,
		`SELECT revoked_at
		FROM tokens
		WHERE token_hash = $1`,
		refreshTokenHash,
	).Scan(&revokedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || revokedAt == nil {
			return ErrTokenInvalidOrExpired
		}
	}

	_, err = pr.db.Exec(
		ctx,
		`UPDATE tokens
		SET revoked_at = NOW()
		WHERE token_hash = $1`,
		refreshTokenHash,
	)

	return err
}
