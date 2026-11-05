package auth

import (
	"github.com/gofiber/fiber/v3"
)

func RegisterAuthRoutes(api fiber.Router, handler AuthHandler) {
	auth := api.Group("/auth")
	auth.Post("/signin", handler.SignIn)
	auth.Post("/signup", handler.SignUp)
	auth.Post("/signout", handler.SignOut)
}
