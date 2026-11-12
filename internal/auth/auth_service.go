package auth

import (
	"context"
	"errors"
	"fmt"
	"health-tracker-api/internal/refresh_token"
	"health-tracker-api/internal/user"
	"health-tracker-api/pkg/apperror"

	"github.com/google/uuid"
)

type AuthService interface {
	SignIn(ctx context.Context, request AuthSignInRequest, userAgent string) (*AuthSignInResponse, error)
	SignUp(ctx context.Context, request user.UserCreateRequest) (*user.UserResponse, error)
	SignOut(ctx context.Context, request AuthSignOutRequest) error
	Refresh(ctx context.Context, request AuthRefreshRequest, userAgent string) (*TokenPair, error)
}

type AuthServiceImpl struct {
	Dependencies
	cfg Config
}

func NewAuthService(deps Dependencies, cfg Config) AuthService {
	return &AuthServiceImpl{Dependencies: deps, cfg: cfg}
}

func errInvalidRefreshToken() error {
	return &apperror.UnauthorizedError{Message: "invalid refresh token"}
}

func (s *AuthServiceImpl) SignUp(ctx context.Context, request user.UserCreateRequest) (*user.UserResponse, error) {
	s.Log.WithContext(ctx).Info("Received sign up request")
	return s.Users.Create(ctx, request)
}

func (s *AuthServiceImpl) SignIn(ctx context.Context, request AuthSignInRequest, userAgent string) (*AuthSignInResponse, error) {
	log := s.Log.WithContext(ctx).With("username", request.Username)
	log.Info("Received sign in request")

	u, err := s.Users.VerifyCredentials(ctx, request.Username, request.Password)
	if err != nil {
		return nil, err
	}

	tokens, err := s.issueTokens(ctx, u.UserID, userAgent)
	if err != nil {
		return nil, err
	}

	log.With("user_id", u.UserID).Info("Sign in successful")
	return &AuthSignInResponse{Tokens: *tokens, User: *u}, nil
}

func (s *AuthServiceImpl) SignOut(ctx context.Context, request AuthSignOutRequest) error {
	log := s.Log.WithContext(ctx)
	log.Info("Received sign out request")

	hash := s.RefreshCodec.Hash(request.RefreshToken)
	if err := s.RefreshTokens.RevokeByHash(ctx, hash, s.Clock.Now()); err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}

	log.Info("Sign out successful")
	return nil
}

func (s *AuthServiceImpl) Refresh(ctx context.Context, request AuthRefreshRequest, userAgent string) (*TokenPair, error) {
	log := s.Log.WithContext(ctx)
	log.Info("Received refresh token request")

	hash, now := s.RefreshCodec.Hash(request.RefreshToken), s.Clock.Now()

	var (
		pair  *TokenPair
		reuse bool
	)

	err := s.Tx.WithinTransaction(ctx, func(ctx context.Context) error {
		token, err := s.RefreshTokens.FindByHash(ctx, hash)
		if errors.Is(err, refresh_token.ErrNotFound) {
			log.Warn("Refresh denied: token not recognized")
			return errInvalidRefreshToken()
		}
		if err != nil {
			return fmt.Errorf("find refresh token: %w", err)
		}

		log := log.With("user_id", token.UserID).With("token_id", token.ID)

		switch {
		case token.IsExpired(now):
			log.Warn("Refresh denied: token expired")
			return errInvalidRefreshToken()

		case token.IsRevoked():
			log.Warn("Refresh token reuse detected, revoking all sessions of the user")
			if err := s.RefreshTokens.RevokeAllByUser(ctx, token.UserID, now); err != nil {
				return fmt.Errorf("revoke all refresh tokens: %w", err)
			}
			reuse = true
			return nil
		}

		revoked, err := s.RefreshTokens.RevokeByID(ctx, token.ID, now)
		if err != nil {
			return fmt.Errorf("revoke refresh token: %w", err)
		}
		if !revoked {
			log.Warn("Refresh denied: token was revoked concurrently")
			return errInvalidRefreshToken()
		}

		pair, err = s.issueTokens(ctx, token.UserID, userAgent)
		return err
	})
	if err != nil {
		return nil, err
	}
	if reuse {
		return nil, errInvalidRefreshToken()
	}

	log.Info("Refresh token successful")
	return pair, nil
}

func (s *AuthServiceImpl) issueTokens(ctx context.Context, userID uuid.UUID, userAgent string) (*TokenPair, error) {
	access, err := s.AccessTokens.Issue(userID)
	if err != nil {
		return nil, fmt.Errorf("issue access token: %w", err)
	}

	raw, hash, err := s.RefreshCodec.Generate()
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	now := s.Clock.Now()
	err = s.RefreshTokens.Create(ctx, &refresh_token.RefreshToken{
		ID:          uuid.New(),
		UserID:      userID,
		TokenHash:   hash,
		DeviceAgent: userAgent,
		ExpiresAt:   now.Add(s.cfg.RefreshTokenTTL),
		CreatedAt:   now,
	})
	if err != nil {
		return nil, fmt.Errorf("save refresh token: %w", err)
	}

	return &TokenPair{AccessToken: access, RefreshToken: raw}, nil
}
