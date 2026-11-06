package gormrepo

import (
	"context"
	"errors"
	"health-tracker-api/internal/refresh_token"
	"health-tracker-api/pkg/database/gormtx"
	"time"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

type Model struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID      uuid.UUID `gorm:"type:uuid;index;not null"`
	TokenHash   string    `gorm:"uniqueIndex;not null"`
	DeviceAgent string
	ExpiresAt   time.Time `gorm:"not null"`
	RevokedAt   *time.Time
	CreatedAt   time.Time
}

func (Model) TableName() string { return "refresh_tokens" }

func toModel(t *refresh_token.RefreshToken) *Model {
	return &Model{
		ID:          t.ID,
		UserID:      t.UserID,
		TokenHash:   t.TokenHash,
		DeviceAgent: t.DeviceAgent,
		ExpiresAt:   t.ExpiresAt,
		RevokedAt:   t.RevokedAt,
		CreatedAt:   t.CreatedAt,
	}
}

func (m Model) toEntity() *refresh_token.RefreshToken {
	return &refresh_token.RefreshToken{
		ID:          m.ID,
		UserID:      m.UserID,
		TokenHash:   m.TokenHash,
		DeviceAgent: m.DeviceAgent,
		ExpiresAt:   m.ExpiresAt,
		RevokedAt:   m.RevokedAt,
		CreatedAt:   m.CreatedAt,
	}
}

type Repository struct{ tm *gormtx.Manager }

func New(tm *gormtx.Manager) refresh_token.RefreshTokenRepository { return &Repository{tm: tm} }

func (r *Repository) Create(ctx context.Context, t *refresh_token.RefreshToken) error {
	return gorm.G[Model](r.tm.DB(ctx)).Create(ctx, toModel(t))
}

func (r *Repository) FindByHash(ctx context.Context, hash string) (*refresh_token.RefreshToken, error) {
	m, err := gorm.G[Model](r.tm.DB(ctx)).Where("token_hash = ?", hash).First(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, refresh_token.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return m.toEntity(), nil
}

func (r *Repository) RevokeByID(ctx context.Context, id uuid.UUID, at time.Time) (bool, error) {
	rows, err := gorm.G[Model](r.tm.DB(ctx)).
		Where("id = ? AND revoked_at IS NULL", id).
		Update(ctx, "revoked_at", at)
	return rows > 0, err
}

func (r *Repository) RevokeByHash(ctx context.Context, hash string, at time.Time) error {
	_, err := gorm.G[Model](r.tm.DB(ctx)).
		Where("token_hash = ? AND revoked_at IS NULL", hash).
		Update(ctx, "revoked_at", at)
	return err
}

func (r *Repository) RevokeAllByUser(ctx context.Context, userID uuid.UUID, at time.Time) error {
	_, err := gorm.G[Model](r.tm.DB(ctx)).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update(ctx, "revoked_at", at)
	return err
}
