package controller

import (
	"health-tracker-api/internal/helper"
	"health-tracker-api/internal/model/web"
	userWeb "health-tracker-api/internal/model/web/user"
	"health-tracker-api/internal/service"

	"github.com/gofiber/fiber/v3"
)

type UserController interface {
	Create(ctx fiber.Ctx) 
}

type UserControllerImpl struct {
	UserService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return &UserControllerImpl{
		UserService: userService,
	}
}

func (controller *UserControllerImpl) Create(c fiber.Ctx) {
	userCreateRequest := userWeb.UserCreateRequest{}
	err := c.Bind().Body(userCreateRequest)
	helper.PanicIfError(err)

	createdUser := controller.UserService.Create(c, userCreateRequest)
	webResponse := web.WebResponse{
		Code: 200,
		Status: "OK",
		Data: createdUser,
	}

	helper.WriteResponseBody(c, webResponse)
}
