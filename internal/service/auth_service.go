package service

import (
	"errors"
	"health-tracker-api/internal/helper"
	authWeb "health-tracker-api/internal/model/web/auth"
	userWeb "health-tracker-api/internal/model/web/user"
	"health-tracker-api/internal/repository"
	"health-tracker-api/pkg/apperror"
	"health-tracker-api/pkg/database"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

type AuthService interface {
	Register(ctx fiber.Ctx, request userWeb.UserCreateRequest) (*userWeb.UserResponse, error)
	Login(ctx fiber.Ctx, request authWeb.AuthLoginRequest) (*authWeb.AuthLoginResponse, error)
	// Logout(ctx fiber.Ctx, token string) error
}

type AuthServiceImpl struct {
	UserService    UserService      
	UserRepository repository.UserRepository         
	Validate       *validator.Validate
	log            *logrus.Logger
}

func NewAuthService(userService UserService, userRepository repository.UserRepository, validate *validator.Validate, log *logrus.Logger) AuthService {
	return &AuthServiceImpl{
		UserService:    userService,
		UserRepository:    userRepository,
		Validate:       validate,
		log:            log,
	}
}


func (s *AuthServiceImpl) Register(c fiber.Ctx, request userWeb.UserCreateRequest) (*userWeb.UserResponse, error) {
	log := helper.LoggerWithRequestID(c, s.log)
	log.Info("Received register request")

	return s.UserService.Create(c, request)
}

func (s *AuthServiceImpl) Login(c fiber.Ctx, request authWeb.AuthLoginRequest) (*authWeb.AuthLoginResponse, error) {
	log := helper.LoggerWithRequestID(c, s.log).WithField("username", request.Username)
	log.Info("Received login request")

	err := s.Validate.Struct(request)
	if err != nil {
		log.WithField("error", err.Error()).Warn("Validation failed for login request")
		return nil, &apperror.ValidationError{Message: err.Error()}
	}

	tx := database.DB.Begin()
	defer helper.CommitOrRollback(tx)

	user, err := s.UserRepository.FindByUsername(c, tx, request.Username)
	if err != nil {
		log.Warn("Login failed: user not found")
		return nil, &apperror.NotFoundError{Message: "invalid username or password"}
	}

	if !helper.CheckPasswordHash(request.Password, user.Password) {
		log.Warn("Login failed: invalid password")
		return nil, &apperror.NotFoundError{Message: "invalid username or password"}
	}

	token, err := helper.GenerateJWT(user)
	if err != nil {
		log.WithField("error", err.Error()).Error("Failed to generate token")
		return nil, errors.New(err.Error())
	}

	log.WithField("user_id", user.ID).Info("Login successful")

	return helper.ToAuthResponse(token, user), nil
}