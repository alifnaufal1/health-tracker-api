package redis

import (
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type RefreshTokenRepository interface {
	Create(c fiber.Ctx, tx *gorm.DB, refreshToken *RefreshToken) (*RefreshToken, error)
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
