package repository

import (
	"health-tracker-api/internal/model/domain"

	"github.com/gofiber/fiber/v3"
)

type DeviceRepository interface {
	Save(ctx fiber.Ctx, device domain.Device) domain.Device
	Update(ctx fiber.Ctx, device domain.Device) domain.Device
	Delete(ctx fiber.Ctx, deviceId string)
	FindById(ctx fiber.Ctx, deviceId string) domain.Device
	FindAll(ctx fiber.Ctx) []domain.Device
}