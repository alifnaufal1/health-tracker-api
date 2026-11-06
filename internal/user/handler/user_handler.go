package handler

import (
	"health-tracker-api/internal/user"
	"health-tracker-api/pkg/context"
	"health-tracker-api/pkg/request"
	"health-tracker-api/pkg/response"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	service  user.UserService
	validate *validator.Validate
}

func NewUserHandler(service user.UserService, validate *validator.Validate) *UserHandler {
	return &UserHandler{
		service:  service,
		validate: validate,
	}
}

func (h *UserHandler) Update(ctx fiber.Ctx) error {
	req, err := request.BindAndValidate[user.UserUpdateRequest](ctx, h.validate)
	if err != nil {
		return err
	}
	authUser, err := context.GetAuthUser(ctx)
	if err != nil {
		return err
	}

	result, err := h.service.Update(ctx, *req, authUser.UserID)
	if err != nil {
		return err
	}

	return response.Success(ctx, result, "success update user")
}

func (h *UserHandler) GetByMe(c fiber.Ctx) error {
	authUser, err := context.GetAuthUser(c)
	if err != nil {
		return err
	}
	user, err := h.service.FindByID(c, authUser.UserID)
	if err != nil {
		return err
	}

	return response.Success(c, user, "success find user")
}
