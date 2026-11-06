package refresh_token

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *RefreshToken) error
	FindByHash(ctx context.Context, hash string) (*RefreshToken, error)
	RevokeByID(ctx context.Context, id uuid.UUID, at time.Time) (revoked bool, err error)
	RevokeByHash(ctx context.Context, hash string, at time.Time) error
	RevokeAllByUser(ctx context.Context, userID uuid.UUID, at time.Time) error
}
