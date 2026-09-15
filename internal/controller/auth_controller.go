package controller

import (
	"health-tracker-api/internal/helper"
	web "health-tracker-api/internal/model/web/auth"
	userWeb "health-tracker-api/internal/model/web/user"
	"health-tracker-api/internal/service"

	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

type AuthController interface {
	Login(ctx fiber.Ctx) error
	Register(ctx fiber.Ctx) error
	// Logout(ctx fiber.Ctx) error
}

type AuthControllerImpl struct {
	AuthService service.AuthService
	log *logrus.Logger
}

func NewAuthController(authService service.AuthService, log *logrus.Logger) AuthController {
	return &AuthControllerImpl{
		AuthService: authService,
		log: log,
	}
}

func (c *AuthControllerImpl) Login(ctx fiber.Ctx) error {
	loginRequest := new(web.AuthLoginRequest)
	err := ctx.Bind().Body(loginRequest)
	if err != nil {
		return err
	}
	
	user, err := c.AuthService.Login(ctx, *loginRequest)
	if err != nil {
		return err
	}

	return helper.ToWebResponse(ctx, user, "success login")
}

func (c *AuthControllerImpl) Register(ctx fiber.Ctx) error {
	registerRequest := new(userWeb.UserCreateRequest)
	err := ctx.Bind().Body(registerRequest)
	if err != nil {
		return err
	}
	
	user, err := c.AuthService.Register(ctx, *registerRequest)
	if err != nil {
		return err
	}

	return helper.ToWebResponse(ctx, user, "success login")
}