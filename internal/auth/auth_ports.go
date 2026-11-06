package auth

import (
	"context"
	"health-tracker-api/internal/refresh_token"
	"health-tracker-api/internal/user"
	"health-tracker-api/pkg/logging"
	"health-tracker-api/pkg/transaction"
	"time"

	"github.com/google/uuid"
)

type UserReader interface {
	FindByUsername(ctx context.Context, username string) (*user.User, error)
}

type AccessTokenIssuer interface {
	Issue(userID uuid.UUID) (string, error)
}

type RefreshTokenCodec interface {
	Generate() (raw, hash string, err error)
	Hash(raw string) string
}

type Clock interface {
	Now() time.Time
}

type Config struct {
	RefreshTokenTTL time.Duration
}

type Dependencies struct {
	Users         user.UserService
	RefreshTokens refresh_token.RefreshTokenRepository
	Tx            transaction.Manager
	AccessTokens  AccessTokenIssuer
	RefreshCodec  RefreshTokenCodec
	Clock         Clock
	Log           logging.Logger
}
