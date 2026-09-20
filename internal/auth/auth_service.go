package auth

import (
	"errors"
	"health-tracker-api/internal/user"
	"health-tracker-api/pkg/apperror"
	"health-tracker-api/pkg/database"
	"health-tracker-api/pkg/helper"
	"health-tracker-api/pkg/jwt"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

type AuthService interface {
	Register(c fiber.Ctx, request user.UserCreateRequest) (*user.UserResponse, error)
	Login(c fiber.Ctx, request AuthLoginRequest) (*AuthLoginResponse, error)
}

type AuthServiceImpl struct {
	UserService    user.UserService      
	UserRepository user.UserRepository         
	Validate       *validator.Validate
	log            *logrus.Logger
}

func NewAuthService(userService user.UserService, userRepository user.UserRepository, validate *validator.Validate, log *logrus.Logger) AuthService {
	return &AuthServiceImpl{
		UserService:    userService,
		UserRepository:    userRepository,
		Validate:       validate,
		log:            log,
	}
}


func (s *AuthServiceImpl) Register(c fiber.Ctx, request user.UserCreateRequest) (*user.UserResponse, error) {
	log := helper.LoggerWithRequestID(c, s.log)
	log.Info("Received register request")

	return s.UserService.Create(c, request)
}

func (s *AuthServiceImpl) Login(c fiber.Ctx, request AuthLoginRequest) (*AuthLoginResponse, error) {
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

	token, err := jwt.GenerateJWT(user)
	if err != nil {
		log.WithField("error", err.Error()).Error("Failed to generate token")
		return nil, errors.New(err.Error())
	}

	log.WithField("user_id", user.ID).Info("Login successful")

	return toAuthResponse(token, user), nil
}

func toAuthResponse(token *string, model *user.User) *AuthLoginResponse {
	return &AuthLoginResponse{
		Token: *token,
		User: *user.ToUserResponse(model),
	}
}