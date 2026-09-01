package service

import (
	"health-tracker-api/internal/database"
	"health-tracker-api/internal/helper"
	"health-tracker-api/internal/model/domain"
	web "health-tracker-api/internal/model/web/user"
	"health-tracker-api/internal/repository"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v3"
)

type UserService interface {
	Create(ctx fiber.Ctx, request web.UserCreateRequest) web.UserResponse
}

type UserServiceImpl struct {
	UserRepository repository.UserRepository
	Validate *validator.Validate
}

func NewUserService(userRepository repository.UserRepository, validate *validator.Validate) UserService {
	return &UserServiceImpl{
		UserRepository: userRepository,
		Validate: validate,
	}
}

func (service *UserServiceImpl) Create(ctx fiber.Ctx, request web.UserCreateRequest) web.UserResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx := database.DB.Begin()
	defer helper.CommitOrRollback(tx)

	hash, err := helper.HashPassword(request.Password)
	helper.PanicIfError(err)

	user := domain.User{
		Username: request.Username,
		Password: hash,
		Name: request.Name,
	}

	user = service.UserRepository.Save(ctx, tx, user)

	return helper.ToUserResponse(user)
}