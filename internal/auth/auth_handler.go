package auth

import (
	"health-tracker-api/internal/user"
	"health-tracker-api/pkg/response"

	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

type AuthHandler interface {
	SignIn(c fiber.Ctx) error
	SignUp(c fiber.Ctx) error
	SignOut(c fiber.Ctx) error
}

type AuthHandlerImpl struct {
	AuthService AuthService
	log         *logrus.Logger
}

func NewAuthHandler(authService AuthService, log *logrus.Logger) AuthHandler {
	return &AuthHandlerImpl{
		AuthService: authService,
		log:         log,
	}
}

func (h *AuthHandlerImpl) SignIn(c fiber.Ctx) error {
	signInRequest := new(AuthSignInRequest)
	err := c.Bind().Body(signInRequest)
	if err != nil {
		return err
	}

	user, err := h.AuthService.SignIn(c, *signInRequest, c.UserAgent())
	if err != nil {
		return err
	}

	return response.Success(c, user, "success sign in")
}

func (h *AuthHandlerImpl) SignUp(c fiber.Ctx) error {
	signUpRequest := new(user.UserCreateRequest)
	err := c.Bind().Body(signUpRequest)
	if err != nil {
		return err
	}

	user, err := h.AuthService.SignUp(c, *signUpRequest)
	if err != nil {
		return err
	}

	return response.Success(c, user, "success sign up")
}

func (h *AuthHandlerImpl) SignOut(c fiber.Ctx) error {
	signOutRequest := new(AuthSignOutRequest)
	err := c.Bind().Body(signOutRequest)
	if err != nil {
		return err
	}

	err = h.AuthService.SignOut(c, *signOutRequest)
	if err != nil {
		return err
	}

	return response.SuccessWithoutData(c, "success sign out")
}
