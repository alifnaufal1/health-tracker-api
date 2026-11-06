package user

import (
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	FindByID(ctx context.Context, userID uuid.UUID) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
}
