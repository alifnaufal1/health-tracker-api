package context

import (
	"errors"

	"github.com/gofiber/fiber/v3"
)

type contextKey string

const AuthUserKey contextKey = "authUser"

func SetAuthUser(c fiber.Ctx, user *AuthenticatedUser) {
	c.Locals(AuthUserKey, user)
}

func GetAuthUser(c fiber.Ctx) (*AuthenticatedUser, error) {
	user, ok := c.Locals(AuthUserKey).(*AuthenticatedUser)
	if !ok || user == nil {
		return nil, errors.New("User not authenticated")
	}
	return user, nil
}