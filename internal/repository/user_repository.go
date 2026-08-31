package repository

import (
	"health-tracker-api/internal/helper"
	"health-tracker-api/internal/model/domain"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type UserRepository interface {
	Save(ctx fiber.Ctx, tx *gorm.DB, user domain.User) domain.User
	Update(ctx fiber.Ctx, user domain.User) domain.User
	Delete(ctx fiber.Ctx, userId string)
	FindById(ctx fiber.Ctx, userId string) domain.User
	FindAll(ctx fiber.Ctx) []domain.User
}

type UserRepositoryImpl struct {}

func (repo *UserRepositoryImpl) Save(ctx fiber.Ctx, tx *gorm.DB, user domain.User) domain.User {
	result := gorm.WithResult()
	err := gorm.G[domain.User](tx, result).Create(ctx, &user)
	helper.PanicIfError(err)
	return user
}
