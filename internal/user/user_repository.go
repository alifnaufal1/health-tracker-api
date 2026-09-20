package user

import (
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type UserRepository interface {
	Save(ctx fiber.Ctx, tx *gorm.DB, user *User) (*User, error)
	Update(ctx fiber.Ctx, tx *gorm.DB, user *User) (*User, error)
	Delete(ctx fiber.Ctx, tx *gorm.DB, userID string) error
	FindById(ctx fiber.Ctx, tx *gorm.DB, userID string) (*User, error)
	FindByUsername(ctx fiber.Ctx, tx *gorm.DB, username string) (*User, error)
	FindAll(ctx fiber.Ctx, tx *gorm.DB) (*[]User, error)
}

type UserRepositoryImpl struct {}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}

func (r *UserRepositoryImpl) Save(ctx fiber.Ctx, tx *gorm.DB, user *User) (*User, error) {
	result := gorm.WithResult()
	err := gorm.G[User](tx, result).Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepositoryImpl) Update(ctx fiber.Ctx, tx *gorm.DB, user *User) (*User, error) {
	result := gorm.WithResult()
	_, err := gorm.G[User](tx, result).Where("id = ?", user.ID).Updates(ctx, *user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepositoryImpl) Delete(ctx fiber.Ctx, tx *gorm.DB, userID string) error {
	result := gorm.WithResult()
	_, err := gorm.G[User](tx, result).Where("id = ?", userID).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepositoryImpl) FindById(ctx fiber.Ctx, tx *gorm.DB, userID string) (*User, error) {
	user, err := gorm.G[User](tx).Where("id = ?", userID).First(ctx)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) FindByUsername(ctx fiber.Ctx, tx *gorm.DB, username string) (*User, error) {
	user, err := gorm.G[User](tx).Where("username = ?", username).First(ctx)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) FindAll(ctx fiber.Ctx, tx *gorm.DB) (*[]User, error) {
	users, err := gorm.G[User](tx).Find(ctx)
	if err != nil {
		return &users, err
	}
	return &users, nil
}