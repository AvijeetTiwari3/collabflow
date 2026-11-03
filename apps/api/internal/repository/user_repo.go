package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sde/collabflow/internal/model"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user with this email already exists")
)

// UserRepository handles user persistence
type UserRepository interface {
	Create(ctx context.Context, u *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
}

type postgresUserRepo struct {
	db *sql.DB
}

// NewUserRepository creates a user repository
func NewUserRepository(db *sql.DB) UserRepository {
	return &postgresUserRepo{db: db}
}

func (r *postgresUserRepo) Create(ctx context.Context, u *model.User) error {
	query := `
		INSERT INTO users (id, email, name, avatar_url, password_hash, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	u.ID = uuid.New()
	u.CreatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query, u.ID, u.Email, u.Name, u.AvatarURL, u.PasswordHash, u.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *postgresUserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, email, name, avatar_url, password_hash, created_at
		FROM users
		WHERE email = $1
		LIMIT 1
	`
	var u model.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID, &u.Email, &u.Name, &u.AvatarURL, &u.PasswordHash, &u.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *postgresUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	query := `
		SELECT id, email, name, avatar_url, password_hash, created_at
		FROM users
		WHERE id = $1
		LIMIT 1
	`
	var u model.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.Email, &u.Name, &u.AvatarURL, &u.PasswordHash, &u.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}
