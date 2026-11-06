package user

import "github.com/google/uuid"

type UserCreateRequest struct {
	Username string `validate:"required,min=3" json:"username"`
	Password string `validate:"min=8,max=12" json:"password"`
	Name     string `validate:"required,min=3,max=50" json:"name"`
}

type UserUpdateRequest struct {
	Password string `validate:"required_without_all=Name NickName Username,omitempty,min=6,max=50" json:"password"`
	Name     string `validate:"required_without_all=Password NickName Username,omitempty,min=3,max=50" json:"name"`
	Username string `validate:"required_without_all=Password Name NickName,omitempty,min=3,max=8" json:"username"`
}

type UserResponse struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Name     string    `json:"name"`
}
