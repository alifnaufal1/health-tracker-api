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
	Update(ctx fiber.Ctx) error
	Delete(ctx fiber.Ctx) error
	FindByID(ctx fiber.Ctx) error
	FindAll(ctx fiber.Ctx) error
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

func (c *UserControllerImpl) Update(ctx fiber.Ctx) error {
	userUppdateRequest := new(userWeb.UserUpdateRequest)
	err := ctx.Bind().Body(userUppdateRequest)
	if err != nil {
		return err
	}
	
	userUppdateRequest.UserID = ctx.Params("id")

	updatedUser, err := c.UserService.Update(ctx, *userUppdateRequest)
	if err != nil {
		return err
	}

	return helper.ToWebResponse(ctx, updatedUser)
}

func (c *UserControllerImpl) Delete(ctx fiber.Ctx) error {
	err := c.UserService.Delete(ctx, ctx.Params("id"))
	if err != nil {
		return err
	}

	return helper.ToWebResponse(ctx, nil)
}

func (c *UserControllerImpl) FindByID(ctx fiber.Ctx) error {	
	user, err := c.UserService.FindByID(ctx, ctx.Params("id"))
	if err != nil {
		return err
	}

	return helper.ToWebResponse(ctx, user)
}

func (c *UserControllerImpl) FindAll(ctx fiber.Ctx) error {	
	users, err := c.UserService.FindAll(ctx)
	if err != nil {
		return err
	}

	return helper.ToWebResponse(ctx, users)
}
