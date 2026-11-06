package router

import (
	authhandler "health-tracker-api/internal/auth/handler"

	"github.com/gofiber/fiber/v3"
)

func RegisterAuthRoutes(api fiber.Router, h *authhandler.AuthHandler) {
	group := api.Group("/auth")
	group.Post("/signin", h.SignIn)
	group.Post("/signup", h.SignUp)
	group.Post("/signout", h.SignOut)
	group.Post("/refresh", h.Refresh)
}
