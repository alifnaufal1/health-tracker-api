package redis

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type RefreshTokenRepository interface {
	Create(c fiber.Ctx, tx *gorm.DB, refreshToken *RefreshToken) (*RefreshToken, error)
	RevokeByHash(c fiber.Ctx, tx *gorm.DB, hash *string, at *time.Time) error
}

type RefreshTokenRepositoryImpl struct{}

func NewRefreshTokenRepository() RefreshTokenRepository {
	return &RefreshTokenRepositoryImpl{}
}

func (r *RefreshTokenRepositoryImpl) Create(c fiber.Ctx, tx *gorm.DB, refreshToken *RefreshToken) (*RefreshToken, error) {
	result := gorm.WithResult()
	err := gorm.G[RefreshToken](tx, result).Create(c, refreshToken)
	if err != nil {
		return nil, err
	}
	return refreshToken, nil
}

func (r *RefreshTokenRepositoryImpl) RevokeByHash(c fiber.Ctx, tx *gorm.DB, hash *string, at *time.Time) error {
	_, err := gorm.G[RefreshToken](tx).
		Where("token_hash = ? AND revoked_at IS NULL", hash).
		Update(c, "revoked_at", *at)
	if err != nil {
		return err
	}
	return nil
}
