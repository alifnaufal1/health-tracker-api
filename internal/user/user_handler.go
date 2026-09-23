package user

import (
	"health-tracker-api/pkg/context"
	"health-tracker-api/pkg/response"

	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

type UserHandler interface {
	Create(c fiber.Ctx) error
	Update(c fiber.Ctx) error
	Delete(c fiber.Ctx) error
	GetByID(c fiber.Ctx) error
	GetAll(c fiber.Ctx) error
	GetByMe(c fiber.Ctx) error
}

type UserHandlerImpl struct {
	UserService UserService
	log         *logrus.Logger
}

func NewUserHandler(userService UserService, log *logrus.Logger) UserHandler {
	return &UserHandlerImpl{
		UserService: userService,
		log:         log,
	}
}

func (h *UserHandlerImpl) Create(c fiber.Ctx) error {
	userCreateRequest := new(UserCreateRequest)
	err := c.Bind().Body(userCreateRequest)
	if err != nil {
		return err
	}

	createdUser, err := h.UserService.Create(c, *userCreateRequest)
	if err != nil {
		return err
	}

	return response.Success(c, createdUser, "success create new user")
}

func (h *UserHandlerImpl) Update(c fiber.Ctx) error {
	userUpdateRequest := new(UserUpdateRequest)
	err := c.Bind().Body(userUpdateRequest)
	if err != nil {
		return err
	}

	updatedUser, err := h.UserService.Update(c, *userUpdateRequest, c.Params("id"))
	if err != nil {
		return err
	}

	return response.Success(c, updatedUser, "success update user")
}

func (h *UserHandlerImpl) Delete(c fiber.Ctx) error {
	err := h.UserService.Delete(c, c.Params("id"))
	if err != nil {
		return err
	}

	return response.Success(c, nil, "success delete user")
}

func (h *UserHandlerImpl) GetByID(c fiber.Ctx) error {
	user, err := h.UserService.FindByID(c, c.Params("id"))
	if err != nil {
		return err
	}

	return response.Success(c, user, "success find user")
}

func (h *UserHandlerImpl) GetAll(c fiber.Ctx) error {
	users, err := h.UserService.FindAll(c)
	if err != nil {
		return err
	}

	return response.Success(c, users, "success find all users")
}

func (h *UserHandlerImpl) GetByMe(c fiber.Ctx) error {
	authUser, err := context.GetAuthUser(c)
	if err != nil {
		return err
	}
	user, err := h.UserService.FindByID(c, authUser.UserID.String())
	if err != nil {
		return err
	}
	
	return response.Success(c, user, "success find user")
}
