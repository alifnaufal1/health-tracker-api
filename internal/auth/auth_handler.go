package auth

import (
	"health-tracker-api/internal/user"
	"health-tracker-api/pkg/response"

	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

type AuthHandler interface {
	Login(c fiber.Ctx) error
	Register(c fiber.Ctx) error
	// Logout(c fiber.Ctx) error
}

type AuthHandlerImpl struct {
	AuthService AuthService
	log *logrus.Logger
}

func NewAuthHandler(authService AuthService, log *logrus.Logger) AuthHandler {
	return &AuthHandlerImpl{
		AuthService: authService,
		log: log,
	}
}

func (h *AuthHandlerImpl) Login(c fiber.Ctx) error {
	loginRequest := new(AuthLoginRequest)
	err := c.Bind().Body(loginRequest)
	if err != nil {
		return err
	}
	
	user, err := h.AuthService.Login(c, *loginRequest)
	if err != nil {
		return err
	}

	return response.Success(c, user, "success login")
}

func (h *AuthHandlerImpl) Register(c fiber.Ctx) error {
	registerRequest := new(user.UserCreateRequest)
	err := c.Bind().Body(registerRequest)
	if err != nil {
		return err
	}
	
	user, err := h.AuthService.Register(c, *registerRequest)
	if err != nil {
		return err
	}

	return response.Success(c, user, "success login")
}