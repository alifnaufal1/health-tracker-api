package redis

import (
	"health-tracker-api/internal/user"
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	TokenHash   string    `gorm:"uniqueIndex;not null;"`
	UserID      uuid.UUID `gorm:"type:uuid;index;not null"`
	User        user.User
	DeviceAgent string    `gorm:"not null"`
	ExpiresAt   time.Time `gorm:"not null"`
	RevokedAt   *time.Time
	CreatedAt   time.Time
}
