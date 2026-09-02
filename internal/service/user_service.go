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
)

type UserService interface {
	Create(c fiber.Ctx, request web.UserCreateRequest) (*web.UserResponse, error)
	// Update(c fiber.Ctx, request web.UserUpdateRequest) (web.UserResponse, error)
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

// func (service *UserServiceImpl) Update(c fiber.Ctx, request web.UserUpdateRequest) (web.UserResponse, error) {
// 	err := service.Validate.Struct(request)
// 	helper.PanicIfError(err)

// 	tx := database.DB.Begin()
// 	defer helper.CommitOrRollback(tx)

// 	user, err := service.UserRepository.FindById(c, tx, request.UserID)
// 	if err != nil {
// 		return user
// 	}

// 	hash, err := helper.HashPassword(request.Password)
// 	helper.PanicIfError(err)

// 	user := domain.User{
// 		Username: request.Username,
// 		Password: hash,
// 		Name: request.Name,
// 	}

// 	user = service.UserRepository.Save(c, tx, user)

// 	return helper.ToUserResponse(user)
// }