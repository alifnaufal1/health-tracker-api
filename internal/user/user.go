package user

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrNotFound          = errors.New("user not found")
	ErrDuplicateUsername = errors.New("username already registered")
)

type User struct {
	ID       uuid.UUID
	Username string
	Name     string
	Password string
}

func ToUserResponse(u *User) *UserResponse {
	return &UserResponse{
		UserID: u.ID, 
		Username: u.Username,
		Name: u.Name,
	}
}


