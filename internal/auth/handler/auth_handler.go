package handler

import (
	"health-tracker-api/internal/auth"
	"health-tracker-api/internal/user"
	"health-tracker-api/pkg/request"
	"health-tracker-api/pkg/response"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v3"
)

type AuthHandler struct {
	service  auth.AuthService
	validate *validator.Validate
}

func NewAuthHandler(service auth.AuthService, validate *validator.Validate) *AuthHandler {
	return &AuthHandler{
		service: service, 
		validate: validate,
	}
}

func (h *AuthHandler) SignIn(c fiber.Ctx) error {
	req, err := request.BindAndValidate[auth.AuthSignInRequest](c, h.validate)
	if err != nil {
		return err
	}

	result, err := h.service.SignIn(request.Context(c), *req, c.UserAgent())
	if err != nil {
		return err
	}

	return response.Success(c, result, "success sign in")
}

func (h *AuthHandler) SignUp(c fiber.Ctx) error {
	req, err := request.BindAndValidate[user.UserCreateRequest](c, h.validate)
	if err != nil {
		return err
	}

	result, err := h.service.SignUp(request.Context(c), *req)
	if err != nil {
		return err
	}

	return response.Success(c, result, "success sign up")
}

func (h *AuthHandler) SignOut(c fiber.Ctx) error {
	req, err := request.BindAndValidate[auth.AuthSignOutRequest](c, h.validate)
	if err != nil {
		return err
	}

	if err := h.service.SignOut(request.Context(c), *req); err != nil {
		return err
	}

	return response.SuccessWithoutData(c, "success sign out")
}

func (h *AuthHandler) Refresh(c fiber.Ctx) error {
	req, err := request.BindAndValidate[auth.AuthRefreshRequest](c, h.validate)
	if err != nil {
		return err
	}

	result, err := h.service.Refresh(request.Context(c), *req, c.UserAgent())
	if err != nil {
		return err
	}

	return response.Success(c, result, "success refresh")
}
