package user

import (
	"health-tracker-api/pkg/middleware"

	"github.com/gofiber/fiber/v3"
)

func RegisterUserRoutes(api fiber.Router, handler UserHandler) {
	user := api.Group("/user")
	user.Post("", handler.Create)
	user.Put("/:id", middleware.JwtProtected(), handler.Update)
	user.Delete("/:id", middleware.JwtProtected(), handler.Delete)
	user.Get("/me", middleware.JwtProtected(), handler.GetByMe)
	user.Get("/:id", middleware.JwtProtected(), handler.GetByID)
	user.Get("", middleware.JwtProtected(), handler.GetAll)
}