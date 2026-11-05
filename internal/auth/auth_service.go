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
	"time"

	"github.com/google/uuid"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

type AuthService interface {
	SignUp(c fiber.Ctx, request user.UserCreateRequest) (*user.UserResponse, error)
	SignIn(c fiber.Ctx, request AuthSignInRequest, userAgent string) (*AuthSignInResponse, error)
	SignOut(c fiber.Ctx, request AuthSignOutRequest) error
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

	accessToken, err := jwt.GenerateJWT(user)
	if err != nil {
		log.WithField("error", err.Error()).Error("Failed to generate access token")
		return nil, errors.New(err.Error())
	}

	raw, hash, err := newRefreshToken()
	if err != nil {
		log.WithField("error", err.Error()).Error("Failed to generate refresh token")
		return nil, errors.New(err.Error())
	}

	token := &redis.RefreshToken{
		ID:          uuid.New(),
		TokenHash:   hash,
		UserID:      user.ID,
		DeviceAgent: userAgent,
		ExpiresAt:   time.Now().Add(refreshTokenTTL),
		CreatedAt:   time.Now(),
	}

	_, err = s.RefreshTokenRepository.Create(c, tx, token)
	if err != nil {
		log.WithError(err).Error("failed to save refresh token")
		return nil, fmt.Errorf("save refresh token: %w", err)
	}

	log.WithField("user_id", user.ID).Info("Sign in successful")

	return toAuthResponse(accessToken, &raw, user), nil
}

func (s *AuthServiceImpl) SignOut(c fiber.Ctx, request AuthSignOutRequest) error {
	log := helper.LoggerWithRequestID(c, s.log).WithField("rawRefreshToken", request.RefreshToken)
	log.Info("Received sign out request")

	hash := helper.HashToken(request.RefreshToken)
	now := time.Now()

	tx := database.DB.Begin()
	defer helper.CommitOrRollback(tx)

	err := s.RefreshTokenRepository.RevokeByHash(c, tx, &hash, &now)
	if err != nil {
		log.WithError(err).Error("failed to update revoked status")
		return fmt.Errorf("update revoked status: %w", err)
	}

	log.Info("Sign out successfull")

	return nil
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
