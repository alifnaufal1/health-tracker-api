package device

import (
	"health-tracker-api/pkg/middleware"

	"github.com/gofiber/fiber/v3"
)

func RegisterDeviceRoutes(api fiber.Router, handler DeviceHandler) {
	device := api.Group("/device")
	device.Post("", handler.Create)
	device.Put("/:id", middleware.JwtProtected(), handler.Update)
	device.Delete("/:id", middleware.JwtProtected(), handler.Delete)
	device.Get("/:id", middleware.JwtProtected(), handler.FindByID)
	device.Get("", middleware.JwtProtected(), handler.FindAll)
}