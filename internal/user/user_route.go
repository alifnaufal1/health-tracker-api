package user

import (
	"health-tracker-api/pkg/middleware"

	"github.com/gofiber/fiber/v3"
)

func RegisterUserRoutes(api fiber.Router, handler UserHandler) {
	user := api.Group("/user")
	user.Post("", handler.Create)
	user.Put("/:id", middleware.Protected(), handler.Update)
	user.Delete("/:id", middleware.Protected(), handler.Delete)
	user.Get("/:id", middleware.Protected(), handler.FindByID)
	user.Get("", middleware.Protected(), handler.FindAll)
}