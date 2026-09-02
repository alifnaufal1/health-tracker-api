package service

import (
	"errors"
	"health-tracker-api/internal/helper"
	"health-tracker-api/internal/model/domain"
	web "health-tracker-api/internal/model/web/user"
	"health-tracker-api/internal/repository"
	"health-tracker-api/pkg/apperror"
	"health-tracker-api/pkg/database"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserService interface {
	Create(c fiber.Ctx, request web.UserCreateRequest) (*web.UserResponse, error)
	Update(c fiber.Ctx, request web.UserUpdateRequest) (*web.UserResponse, error)
	Delete(c fiber.Ctx, userID string) error
	FindByID(c fiber.Ctx, userID string) (*web.UserResponse, error)
	FindAll(c fiber.Ctx) (*[]web.UserResponse, error)
}

type UserServiceImpl struct {
	UserRepository repository.UserRepository
	Validate *validator.Validate
	log *logrus.Logger
}

func NewUserService(userRepository repository.UserRepository, validate *validator.Validate, log *logrus.Logger) UserService {
	return &UserServiceImpl{
		UserRepository: userRepository,
		Validate: validate,
		log: log,
	}
}

func (s *UserServiceImpl) Create(c fiber.Ctx, request web.UserCreateRequest) (*web.UserResponse, error) {
	log := helper.LoggerWithRequestID(c, s.log).WithField("username", request.Username)

	log.Info("Received create user request")
	
	err := s.Validate.Struct(request)
	if err != nil {
		log.WithField("error", err.Error()).Warn("Validation failed for create user request")
		return nil, &apperror.ValidationError{Message: err.Error()}
	}
	
	tx := database.DB.Begin()
	defer helper.CommitOrRollback(tx)
	
	hash, err := helper.HashPassword(request.Password)
	if err != nil {
		log.Error("Failed to hash password")
		return nil, errors.New(err.Error())
	}

	user := &domain.User{
		Base: domain.Base{ID: uuid.New()},
		Username: request.Username,
		Password: hash,
		Name: request.Name,
	}

	user, err = s.UserRepository.Save(c, tx, user)
	if err != nil {
		log.Warn("Create user process failed at repository layer")
		return nil, errors.New(err.Error())
	}

	log.WithField("user_id", user.ID).Info("User created successfully")

	return helper.ToUserResponse(user), nil
}

func (s *UserServiceImpl) Update(c fiber.Ctx, request web.UserUpdateRequest) (*web.UserResponse, error) {
	log := helper.LoggerWithRequestID(c, s.log).WithField("user_id", request.UserID)
	
	log.Info("Received update user request")

	err := s.Validate.Struct(request)
	if err != nil {
		log.WithField("error", err.Error()).Warn("Validation failed for update user request")
		return nil, &apperror.ValidationError{Message: err.Error()}
	}

	tx := database.DB.Begin()
	defer helper.CommitOrRollback(tx)

	user, err := s.UserRepository.FindById(c, tx, request.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &apperror.NotFoundError{Message: "user not found"}
		}
		return nil, errors.New(err.Error())
	}

	hash, err := helper.HashPassword(request.Password)
	if err != nil {
		log.Error("Failed to hash password")
		return nil, errors.New(err.Error())
	}

	user = &domain.User{
		Username: request.Username,
		Password: hash,
		Name: request.Name,
	}

	user, err = s.UserRepository.Update(c, tx, user)
	if err != nil {
		log.Warn("Create user process failed at repository layer")
		return nil, errors.New(err.Error())
	}

	log.WithField("user_id", user.ID).Info("User created successfully")

	return helper.ToUserResponse(user), nil
}

func (s *UserServiceImpl) Delete(c fiber.Ctx, userID string) error {
	log := helper.LoggerWithRequestID(c, s.log).WithField("user_id", userID)
	
	log.Info("Received delete user request")

	if userID == "" {
		log.Warn("Validation failed for delete user request")
		return &apperror.ValidationError{Message: "user_id required"}
	}

	tx := database.DB.Begin()
	defer helper.CommitOrRollback(tx)

	user, err := s.UserRepository.FindById(c, tx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &apperror.NotFoundError{Message: "user not found"}
		}
		return errors.New(err.Error())
	}

	err = s.UserRepository.Delete(c, tx, userID)
	if err != nil {
		log.Warn("Delete user process failed at repository layer")
		return errors.New(err.Error())
	}

	log.WithField("user_id", user.ID).Info("User deleted successfully")

	return nil
}

func (s *UserServiceImpl) FindByID(c fiber.Ctx, userID string) (*web.UserResponse, error) {
	log := helper.LoggerWithRequestID(c, s.log).WithField("user_id", userID)
	
	log.Info("Received find user by id request")

	if userID == "" {
		log.Warn("Validation failed for find user by id request")
		return nil, &apperror.ValidationError{Message: "user_id required"}
	}

	tx := database.DB.Begin()
	defer helper.CommitOrRollback(tx)

	user, err := s.UserRepository.FindById(c, tx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &apperror.NotFoundError{Message: "user not found"}
		}
		return nil, errors.New(err.Error())
	}

	log.WithField("user_id", user.ID).Info("User found successfully")

	return helper.ToUserResponse(user), nil
}

func (s *UserServiceImpl) FindAll(c fiber.Ctx) (*[]web.UserResponse, error) {
	log := helper.LoggerWithRequestID(c, s.log)

	log.Info("Received find all user request")

	tx := database.DB.Begin()
	defer helper.CommitOrRollback(tx)

	users, err := s.UserRepository.FindAll(c, tx)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	log.Info("Users found successfully")

	return helper.ToUserResponses(users), nil
}