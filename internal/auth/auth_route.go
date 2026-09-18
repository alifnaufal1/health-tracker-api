package auth

import (
	"github.com/gofiber/fiber/v3"
)

func RegisterAuthRoutes(api fiber.Router, handler AuthHandler) {
	auth := api.Group("/auth")
	auth.Post("/login", handler.Login)
	auth.Post("/register", handler.Register)
}