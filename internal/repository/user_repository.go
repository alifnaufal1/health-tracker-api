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
	log := helper.LoggerWithRequestID(ctx, r.log).WithField("user_id", user.ID)

	log.Debug("Updating user into database")
	
	result := gorm.WithResult()
	_, err := gorm.G[domain.User](tx, result).Where("id = ?", user.ID).Updates(ctx, *user)
	if err != nil {
		log.WithField("error", err.Error()).Error("Database update failed")
		return nil, err
	}

	log.WithField("user_id", user.ID).Debug("User updated successfully")
	
	return user, nil
}

func (r *UserRepositoryImpl) Delete(ctx fiber.Ctx, tx *gorm.DB, userID string) error {
	log := helper.LoggerWithRequestID(ctx, r.log).WithField("user_id", userID)
	
	log.Debug("Deleting user from database")
	
	result := gorm.WithResult()
	_, err := gorm.G[domain.User](tx, result).Where("id = ?", userID).Delete(ctx)
	if err != nil {
		log.WithField("error", err.Error()).Error("Database delete failed")
		return err
	}

	log.WithField("user_id", userID).Debug("User deleted successfully")
	
	return nil
}

func (r *UserRepositoryImpl) FindById(ctx fiber.Ctx, tx *gorm.DB, userID string) (*domain.User, error) {
	log := helper.LoggerWithRequestID(ctx, r.log).WithField("user_id", userID)
	
	log.Debug("Finding user by id from database")
	
	users, err := gorm.G[domain.User](tx).Select("id", "username", "name").Where("id = ?", userID).Find(ctx)
	if err != nil {
		log.WithField("error", err.Error()).Error("Database find by id failed")
		return nil, err
	}
	
	log.WithField("user_id", userID).Debug("User found by id successfully")
	
	return &users[0], nil
}

func (r *UserRepositoryImpl) FindAll(ctx fiber.Ctx, tx *gorm.DB) (*[]domain.User, error) {
	log := helper.LoggerWithRequestID(ctx, r.log)

	log.Debug("Finding all user from database")
	
	users, err := gorm.G[domain.User](tx).Select("id", "username", "name").Find(ctx)
	if err != nil {
		log.WithField("error", err.Error()).Error("Database find all failed")
		return &users, err
	}

	log.Debug("All user found successfully")
	
	return &users, nil
}
