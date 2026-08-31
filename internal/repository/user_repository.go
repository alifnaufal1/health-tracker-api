package repository

import (
	"health-tracker-api/internal/model/domain"

	"github.com/gofiber/fiber/v3"
)

type UserRepository interface {
	Save(ctx fiber.Ctx, user domain.User) domain.User
	Update(ctx fiber.Ctx, user domain.User) domain.User
	Delete(ctx fiber.Ctx, userId string)
	FindById(ctx fiber.Ctx, userId string) domain.User
	FindAll(ctx fiber.Ctx) []domain.User
}

