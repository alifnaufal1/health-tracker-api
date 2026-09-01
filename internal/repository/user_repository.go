package repository

import (
	"health-tracker-api/internal/helper"
	"health-tracker-api/internal/model/domain"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type UserRepository interface {
	Save(c fiber.Ctx, tx *gorm.DB, user domain.User) domain.User
	// Update(c fiber.Ctx, tx *gorm.DB, user domain.User) domain.User
	// Delete(c fiber.Ctx, tx *gorm.DB, userId string)
	// FindById(c fiber.Ctx, tx *gorm.DB, userId string) domain.User
	// FindAll(c fiber.Ctx, tx *gorm.DB) []domain.User
}

type UserRepositoryImpl struct {}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}

func (repo *UserRepositoryImpl) Save(c fiber.Ctx, tx *gorm.DB, user domain.User) domain.User {
	result := gorm.WithResult()
	err := gorm.G[domain.User](tx, result).Create(c, &user)
	helper.PanicIfError(err)
	return user
}

// func (repo *UserRepositoryImpl) Update(c fiber.Ctx, tx *gorm.DB, user domain.User) domain.User {
// 	result := gorm.WithResult()
// 	err := gorm.G[domain.User](tx, result).Update(c, &user)
// 	helper.PanicIfError(err)
// 	return user
// }
