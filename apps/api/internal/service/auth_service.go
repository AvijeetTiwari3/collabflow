package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/sde/collabflow/internal/auth"
	"github.com/sde/collabflow/internal/model"
	"github.com/sde/collabflow/internal/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
)

// AuthService defines authentication business operations
type AuthService interface {
	Register(ctx context.Context, req model.RegisterRequest) (*model.AuthResponse, error)
	Login(ctx context.Context, req model.LoginRequest) (*model.AuthResponse, error)
	RefreshToken(ctx context.Context, rawRefreshToken string) (*model.TokenPair, error)
	Logout(ctx context.Context, userID uuid.UUID) error
}

type authService struct {
	userRepo    repository.UserRepository
	refreshRepo repository.RefreshTokenRepository
	tokenMgr    *auth.TokenManager
}

// NewAuthService creates an auth service instance
func NewAuthService(
	userRepo repository.UserRepository,
	refreshRepo repository.RefreshTokenRepository,
	tokenMgr *auth.TokenManager,
) AuthService {
	return &authService{
		userRepo:    userRepo,
		refreshRepo: refreshRepo,
		tokenMgr:    tokenMgr,
	}
}

func (s *authService) Register(ctx context.Context, req model.RegisterRequest) (*model.AuthResponse, error) {
	// Normalize email
	email := strings.ToLower(strings.TrimSpace(req.Email))

	// Check if already registered
	existing, _ := s.userRepo.GetByEmail(ctx, email)
	if existing != nil {
		return nil, repository.ErrUserAlreadyExists
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := model.User{
		Email:        email,
		Name:         strings.TrimSpace(req.Name),
		AvatarURL:    fmt.Sprintf("https://api.dicebear.com/7.x/initials/svg?seed=%s", req.Name),
		PasswordHash: hash,
	}

	if err := s.userRepo.Create(ctx, &user); err != nil {
		return nil, err
	}

	return s.generateAuthResponse(ctx, user)
}

func (s *authService) Login(ctx context.Context, req model.LoginRequest) (*model.AuthResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !auth.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	return s.generateAuthResponse(ctx, *user)
}

func (s *authService) RefreshToken(ctx context.Context, rawRefreshToken string) (*model.TokenPair, error) {
	tokenHash := auth.HashRefreshToken(rawRefreshToken)

	// Atomic validation and revocation (Token Rotation)
	userID, err := s.refreshRepo.ValidateAndRevoke(ctx, tokenHash)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Issue new token pair
	accessToken, err := s.tokenMgr.GenerateAccessToken(user.ID, user.Email, user.Name)
	if err != nil {
		return nil, err
	}

	newRawRefresh, newRefreshHash, expiresAt, err := s.tokenMgr.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	if err := s.refreshRepo.Store(ctx, user.ID, newRefreshHash, expiresAt); err != nil {
		return nil, err
	}

	return &model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRawRefresh,
		ExpiresIn:    15 * 60, // 15 mins
	}, nil
}

func (s *authService) Logout(ctx context.Context, userID uuid.UUID) error {
	return s.refreshRepo.RevokeAllForUser(ctx, userID)
}

func (s *authService) generateAuthResponse(ctx context.Context, user model.User) (*model.AuthResponse, error) {
	accessToken, err := s.tokenMgr.GenerateAccessToken(user.ID, user.Email, user.Name)
	if err != nil {
		return nil, err
	}

	rawRefresh, refreshHash, expiresAt, err := s.tokenMgr.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	if err := s.refreshRepo.Store(ctx, user.ID, refreshHash, expiresAt); err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		User: user,
		Tokens: model.TokenPair{
			AccessToken:  accessToken,
			RefreshToken: rawRefresh,
			ExpiresIn:    15 * 60,
		},
	}, nil
}
