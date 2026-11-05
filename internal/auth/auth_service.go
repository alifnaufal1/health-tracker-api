package auth

import (
	"errors"
	"fmt"
	"health-tracker-api/internal/user"
	"health-tracker-api/pkg/apperror"
	"health-tracker-api/pkg/database"
	"health-tracker-api/pkg/helper"
	"health-tracker-api/pkg/jwt"
	"health-tracker-api/pkg/redis"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

type AuthService interface {
	SignUp(c fiber.Ctx, request user.UserCreateRequest) (*user.UserResponse, error)
	SignIn(c fiber.Ctx, request AuthSignInRequest, userAgent string) (*AuthSignInResponse, error)
}

type AuthServiceImpl struct {
	UserService            user.UserService
	UserRepository         user.UserRepository
	RefreshTokenRepository redis.RefreshTokenRepository
	Validate               *validator.Validate
	log                    *logrus.Logger
}

func NewAuthService(userService user.UserService, userRepository user.UserRepository, refreshTokenRepository redis.RefreshTokenRepository, validate *validator.Validate, log *logrus.Logger) AuthService {
	return &AuthServiceImpl{
		UserService:            userService,
		UserRepository:         userRepository,
		RefreshTokenRepository: refreshTokenRepository,
		Validate:               validate,
		log:                    log,
	}
}

func (s *AuthServiceImpl) SignUp(c fiber.Ctx, request user.UserCreateRequest) (*user.UserResponse, error) {
	log := helper.LoggerWithRequestID(c, s.log)
	log.Info("Received sign up request")

	return s.UserService.Create(c, request)
}

func (s *AuthServiceImpl) SignIn(c fiber.Ctx, request AuthSignInRequest, userAgent string) (*AuthSignInResponse, error) {
	log := helper.LoggerWithRequestID(c, s.log).WithField("username", request.Username)
	log.Info("Received sign in request")

	err := s.Validate.Struct(request)
	if err != nil {
		log.WithField("error", err.Error()).Warn("Validation failed for sign in request")
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

	accessToken, err := jwt.GenerateAccessToken(user)
	if err != nil {
		log.WithField("error", err.Error()).Error("Failed to generate access token")
		return nil, errors.New(err.Error())
	}

	refreshToken, err := jwt.GenerateRefreshToken(user)
	if err != nil {
		log.WithField("error", err.Error()).Error("Failed to generate refresh token")
		return nil, errors.New(err.Error())
	}

	token := &redis.RefreshToken{
		RefreshToken: *refreshToken,
		UserID:       user.ID,
		DeviceAgent:  userAgent,
		IsRevoked:    false,
	}

	_, err = s.RefreshTokenRepository.Create(c, tx, token)
	if err != nil {
		log.WithError(err).Error("failed to save refresh token")
		return nil, fmt.Errorf("save refresh token: %w", err)
	}

	log.WithField("user_id", user.ID).Info("Sign in successful")

	return toAuthResponse(accessToken, refreshToken, user), nil
}

func toAuthResponse(accessToken *string, refreshToken *string, model *user.User) *AuthSignInResponse {
	return &AuthSignInResponse{
		Tokens: Tokens{
			AccessToken:  *accessToken,
			RefreshToken: *refreshToken,
		},
		User: *user.ToUserResponse(model),
	}
}
