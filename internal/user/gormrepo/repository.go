package gormrepo

import (
	"context"
	"errors"
	"fmt"
	"health-tracker-api/internal/user"
	"health-tracker-api/pkg/database/gormtx"
	"health-tracker-api/pkg/model"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

type Model struct {
	model.Base
	Username string `gorm:"unique;not null;"`
	Password string `gorm:"not null;"`
	Name     string `gorm:"not null;size:50;"`
}

func (m *Model) toEntity() *user.User {
	return &user.User{
		ID:       m.ID,
		Username: m.Username,
		Name:     m.Name,
		Password: m.Password,
	}
}

func toModel(e *user.User) *Model {
	return &Model{
		Base: model.Base{
			ID: e.ID,
		},
		Username: e.Username,
		Name:     e.Name,
		Password: e.Password,
	}
}

type Repository struct {
	tm *gormtx.Manager
}

func New(tm *gormtx.Manager) user.UserRepository {
	return &Repository{
		tm: tm,
	}
}

func (r *Repository) Create(ctx context.Context, user *user.User) error {
	return gorm.G[Model](r.tm.DB(ctx)).Create(ctx, toModel(user))
}

func (r *Repository) Update(ctx context.Context, u *user.User) error {
	rows, err := gorm.G[Model](r.tm.DB(ctx)).
		Where("id = ?", u.ID).
		Updates(ctx, Model{Name: u.Name, Password: u.Password})
	if err != nil {
		return err
	}
	if rows == 0 {
		return user.ErrNotFound
	}
	return nil
}

func (r *Repository) FindByID(ctx context.Context, userID uuid.UUID) (*user.User, error) {
	return r.findUser("id", userID.String(), ctx)
}

func (r *Repository) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	return r.findUser("username", username, ctx)
}

func (r *Repository) findUser(column string, value string, ctx context.Context) (*user.User, error) {
	query := fmt.Sprintf("%s = ?", column)
	m, err := gorm.G[Model](r.tm.DB(ctx)).
		Where(query, value).
		First(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, user.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return m.toEntity(), nil
}
