package device

import (
	"health-tracker-api/pkg/middleware"

	"github.com/gofiber/fiber/v3"
)

func RegisterDeviceRoutes(api fiber.Router, handler DeviceHandler) {
	device := api.Group("/devices")
	device.Post("", middleware.JwtProtected(), handler.Create)
	device.Put("/:id", middleware.JwtProtected(), handler.Update)
	device.Delete("/:id", middleware.JwtProtected(), handler.Delete)
	device.Get("/me", middleware.JwtProtected(), handler.GetByUserID)
	device.Get("/:id", middleware.JwtProtected(), handler.GetByID)
	device.Get("", middleware.JwtProtected(), handler.GetAll)
}