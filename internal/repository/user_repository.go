package repository

import (
	"health-tracker-api/internal/helper"
	"health-tracker-api/internal/model/domain"

	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository interface {
	Save(ctx fiber.Ctx, tx *gorm.DB, user *domain.User) (*domain.User, error)
	Update(ctx fiber.Ctx, tx *gorm.DB, user *domain.User) (*domain.User, error)
	Delete(ctx fiber.Ctx, tx *gorm.DB, userID string) error
	FindById(ctx fiber.Ctx, tx *gorm.DB, userID string) (*domain.User, error)
	FindAll(ctx fiber.Ctx, tx *gorm.DB) (*[]domain.User, error)
}

type UserRepositoryImpl struct {
	log *logrus.Logger
}

func NewUserRepository(log *logrus.Logger) UserRepository {
	return &UserRepositoryImpl{log: log}
}

func (r *UserRepositoryImpl) Save(ctx fiber.Ctx, tx *gorm.DB, user *domain.User) (*domain.User, error) {
	log := helper.LoggerWithRequestID(ctx, r.log).WithField("username", user.Username)

	log.Debug("Inserting new user into database")

	result := gorm.WithResult()
	err := gorm.G[domain.User](tx, result).Create(ctx, user)
	if err != nil {
		log.WithField("error", err.Error()).Error("Database insert failed")
		return nil, err
	}

	log.WithField("user_id", user.ID).Debug("User inserted successfully")

	return user, nil
}

func (r *UserRepositoryImpl) Update(ctx fiber.Ctx, tx *gorm.DB, user *domain.User) (*domain.User, error) {
	r.log.WithField("user", user).Info("Execution of update to db")
	
	result := gorm.WithResult()
	_, err := gorm.G[domain.User](tx, result).Where("id = ?", user.ID).Updates(ctx, *user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *UserRepositoryImpl) Delete(ctx fiber.Ctx, tx *gorm.DB, userID string) error {
	r.log.WithField("userID", userID).Info("Execution of delete to db")
	
	result := gorm.WithResult()
	_, err := gorm.G[domain.User](tx, result).Where("id = ?", userID).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepositoryImpl) FindById(ctx fiber.Ctx, tx *gorm.DB, userID string) (*domain.User, error) {
	r.log.WithField("userID", userID).Info("Execution of query from db")
	
	users, err := gorm.G[domain.User](tx).Select("id", "username", "name").Where("id = ?", userID).Find(ctx)
	if err != nil {
		return &users[0], err
	}
	return &users[0], nil
}

func (r *UserRepositoryImpl) FindAll(ctx fiber.Ctx, tx *gorm.DB) (*[]domain.User, error) {
	r.log.Info("Execution of query from db")

	users, err := gorm.G[domain.User](tx).Select("id", "username", "name").Find(ctx)
	if err != nil {
		return &users, err
	}
	return &users, nil
}
