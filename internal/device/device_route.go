package device

import (
	"health-tracker-api/pkg/middleware"

	"github.com/gofiber/fiber/v3"
)

func RegisterDeviceRoutes(api fiber.Router, handler DeviceHandler) {
	device := api.Group("/device")
	device.Post("", handler.Create)
	device.Put("/:id", middleware.Protected(), handler.Update)
	device.Delete("/:id", middleware.Protected(), handler.Delete)
	device.Get("/:id", middleware.Protected(), handler.FindByID)
	device.Get("", middleware.Protected(), handler.FindAll)
}