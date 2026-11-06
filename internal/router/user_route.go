package router

import (
	userHandler "health-tracker-api/internal/user/handler"
	"health-tracker-api/pkg/middleware"

	"github.com/gofiber/fiber/v3"
)

func RegisterUserRoutes(api fiber.Router, h *userHandler.UserHandler) {
	user := api.Group("/user")
	user.Put("/me", middleware.JwtProtected(), h.Update)
	user.Get("/me", middleware.JwtProtected(), h.GetByMe)
}
