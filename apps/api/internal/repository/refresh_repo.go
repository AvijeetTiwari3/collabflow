package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrTokenRevokedOrNotFound = errors.New("refresh token is invalid or revoked")
)

// RefreshTokenRepository handles secure refresh token tracking
type RefreshTokenRepository interface {
	Store(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error
	ValidateAndRevoke(ctx context.Context, tokenHash string) (uuid.UUID, error)
	RevokeAllForUser(ctx context.Context, userID uuid.UUID) error
}

type postgresRefreshTokenRepo struct {
	db *sql.DB
}

// NewRefreshTokenRepository creates a refresh token repository
func NewRefreshTokenRepository(db *sql.DB) RefreshTokenRepository {
	return &postgresRefreshTokenRepo{db: db}
}

func (r *postgresRefreshTokenRepo) Store(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, revoked, created_at)
		VALUES ($1, $2, $3, $4, FALSE, NOW())
	`
	_, err := r.db.ExecContext(ctx, query, uuid.New(), userID, tokenHash, expiresAt)
	return err
}

func (r *postgresRefreshTokenRepo) ValidateAndRevoke(ctx context.Context, tokenHash string) (uuid.UUID, error) {
	// Atomic validate and immediate revocation to prevent replay attacks
	query := `
		UPDATE refresh_tokens
		SET revoked = TRUE
		WHERE token_hash = $1 AND revoked = FALSE AND expires_at > NOW()
		RETURNING user_id
	`
	var userID uuid.UUID
	err := r.db.QueryRowContext(ctx, query, tokenHash).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, ErrTokenRevokedOrNotFound
		}
		return uuid.Nil, err
	}
	return userID, nil
}

func (r *postgresRefreshTokenRepo) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE refresh_tokens SET revoked = TRUE WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}
