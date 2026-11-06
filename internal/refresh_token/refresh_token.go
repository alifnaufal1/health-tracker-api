package refresh_token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("refresh token not found")

type RefreshToken struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	TokenHash   string
	DeviceAgent string
	ExpiresAt   time.Time
	RevokedAt   *time.Time
	CreatedAt   time.Time
}

func (t *RefreshToken) IsExpired(now time.Time) bool { return now.After(t.ExpiresAt) }
func (t *RefreshToken) IsRevoked() bool              { return t.RevokedAt != nil }
