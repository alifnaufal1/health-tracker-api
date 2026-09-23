package user

import (
	"errors"
	"fmt"
	"health-tracker-api/pkg/apperror"
	"health-tracker-api/pkg/database"
	"health-tracker-api/pkg/helper"

	"health-tracker-api/pkg/model"
	"health-tracker-api/pkg/response"

	"github.com/google/uuid"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserService interface {
	Create(c fiber.Ctx, request UserCreateRequest) (*UserResponse, error)
	Update(c fiber.Ctx, request UserUpdateRequest, userID string) (*UserResponse, error)
	Delete(c fiber.Ctx, userID string) error
	FindByID(c fiber.Ctx, userID string) (*UserResponse, error)
	FindAll(c fiber.Ctx) (*[]UserResponse, error)
}

type UserServiceImpl struct {
	UserRepository UserRepository
	Validate       *validator.Validate
	log            *logrus.Entry
	
}

func NewUserService(userRepository UserRepository, validate *validator.Validate, base *logrus.Logger) UserService {
	return &UserServiceImpl{
		UserRepository: userRepository,
		Validate:       validate,
		log:            helper.NewModuleLogger(base, "service", "user") ,
	}
}

func (s *UserServiceImpl) Create(c fiber.Ctx, request UserCreateRequest) (*UserResponse, error) {
	log := helper.LoggerWithRequestID(c, s.log.Logger).WithField("request", request)

	log.Debug("received create user request")

	err := s.Validate.Struct(request)
	if err != nil {
		log.WithError(err).Warn("failed to validate create user request")
		return nil, &apperror.ValidationError{Message: err.Error()}
	}

	tx := database.DB.Begin()
	defer helper.CommitOrRollback(tx)

	user, err := s.UserRepository.FindByUsername(c, tx, request.Username)
	if user != nil {
		log.Warn("request denied: this user already registered")
		return nil, &apperror.ConflictError{Message: "this user already registered"}
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound)  {
		log.WithError(err).Error("failed to check registered user")
		return nil, err
	}

	hash, err := helper.HashPassword(request.Password)
	if err != nil {
		log.WithError(err).Error("failed to hash password")
		return nil, err
	}

	user = &User{
		Base:     model.Base{ID: uuid.New()},
		Username: request.Username,
		Password: hash,
		Name:     request.Name,
		NickName: request.NickName,
	}

	user, err = s.UserRepository.Save(c, tx, user)
	if err != nil {
		log.WithError(err).Error("failed to save user")
		return nil, fmt.Errorf("save user: %w", err)
	}

	log.WithField("user_id", user.ID).Debug("user created successfully")

	return ToUserResponse(user), nil
}

func (s *UserServiceImpl) Update(c fiber.Ctx, request UserUpdateRequest, userID string) (*UserResponse, error) {
	log := helper.LoggerWithRequestID(c, s.log.Logger).WithField("user_id", userID)

	log.Debug("received update user request")

	err := s.Validate.Struct(request)
	if err != nil {
		log.WithError(err).Warn("failed to validate update user request")
		return nil, &apperror.ValidationError{Message: err.Error()}
	}

	tx := database.DB.Begin()
	defer helper.CommitOrRollback(tx)

	user, err := s.UserRepository.FindById(c, tx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("not found: this user is not found")
			return nil, &apperror.NotFoundError{Message: "user not found"}
		}
		log.WithError(err).Error("failed to find user")
		return nil, err
	}

	if request.Password != "" {
		hash, err := helper.HashPassword(request.Password)
		if err != nil {
			log.WithError(err).Error("failed to hash password")
			return nil, err
		}
		user.Password = hash
	}
	if request.Name != "" {
		user.Name = request.Name
	}
	if request.NickName != "" {
		user.NickName = request.NickName
	}

	user, err = s.UserRepository.Update(c, tx, user)
	if err != nil {
		log.WithError(err).Error("failed to update workout data")
		return nil, fmt.Errorf("update user by user_id %s: %w", userID, err)
	}

	log.WithField("user_id", userID).Debug("user updated successfully")

	return ToUserResponse(user), nil
}

func (s *UserServiceImpl) Delete(c fiber.Ctx, userID string) error {
	log := helper.LoggerWithRequestID(c, s.log.Logger).WithField("user_id", userID)

	log.Debug("received delete user request")

	if userID == "" {
		log.Warn("validation failed: user_id required")
		return &apperror.ValidationError{Message: "user_id required"}
	}

	tx := database.DB.Begin()
	defer helper.CommitOrRollback(tx)

	_, err := s.UserRepository.FindById(c, tx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("not found: this user is not found")
			return &apperror.NotFoundError{Message: "user not found"}
		}
		log.WithError(err).Error("failed to find user")
		return err
	}

	err = s.UserRepository.Delete(c, tx, userID)
	if err != nil {
		log.WithError(err).Error("failed to delete user by user_id")
		return fmt.Errorf("delete user by user_id %s: %w", userID, err)
	}

	log.WithField("user_id", userID).Debug("user deleted successfully")

	return nil
}

func (s *UserServiceImpl) FindByID(c fiber.Ctx, userID string) (*UserResponse, error) {
	log := helper.LoggerWithRequestID(c, s.log.Logger).WithField("user_id", userID)

	log.Debug("received find user by user_id request")

	if userID == "" {
		log.Warn("validation failed: user_id required")
		return nil, &apperror.ValidationError{Message: "user_id required"}
	}

	tx := database.DB.Begin()
	defer helper.CommitOrRollback(tx)

	user, err := s.UserRepository.FindById(c, tx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("not found: this user is not found")
			return nil, &apperror.NotFoundError{Message: "user not found"}
		}
		log.WithError(err).Error("failed to find user by user_id")
		return nil, fmt.Errorf("find user by user_d %s: %w", userID, err)
	}

	log.WithField("user_id", user.ID).Debug("user found successfully")

	return ToUserResponse(user), nil
}

func (s *UserServiceImpl) FindAll(c fiber.Ctx) (*[]UserResponse, error) {
	log := helper.LoggerWithRequestID(c, s.log.Logger)

	log.Debug("received find all user request")

	tx := database.DB.Begin()
	defer helper.CommitOrRollback(tx)

	users, err := s.UserRepository.FindAll(c, tx)
	if err != nil {
		log.WithError(err).Error("failed to find all user")
		return nil, err
	}

	log.WithField("result_count", len(*users)).Debug("users found successfully")

	return response.ToResponses(users, ToUserResponse), nil
}

func ToUserResponse(user *User) *UserResponse {
	return &UserResponse{
		UserID:   user.Base.ID.String(),
		Name: user.Name,
		NickName: user.NickName,
		Username: user.Username,
	}
}