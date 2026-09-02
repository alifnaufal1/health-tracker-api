package controller

import (
	"health-tracker-api/internal/helper"
	userWeb "health-tracker-api/internal/model/web/user"
	"health-tracker-api/internal/service"

	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

type UserController interface {
	Create(ctx fiber.Ctx) error
}

type UserControllerImpl struct {
	UserService service.UserService
	log *logrus.Logger
}

func NewUserController(userService service.UserService, log *logrus.Logger) UserController {
	return &UserControllerImpl{
		UserService: userService,
		log: log,
	}
}

func (c *UserControllerImpl) Create(ctx fiber.Ctx) error {
	userCreateRequest := new(userWeb.UserCreateRequest)
	err := ctx.Bind().Body(userCreateRequest)
	if err != nil {
		return err
	}
	
	createdUser, err := c.UserService.Create(ctx, *userCreateRequest)
	if err != nil {
		return err
	}

	return helper.ToWebResponse(ctx, createdUser)
}
