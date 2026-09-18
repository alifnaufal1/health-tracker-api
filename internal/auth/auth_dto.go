package auth

import "health-tracker-api/internal/user"

type AuthLoginRequest struct {
	Username string `validate:"required" json:"username"`
	Password string `validate:"required,min=8,max=12" json:"password"`
}

type AuthLoginResponse struct {
	Token string            `json:"token"`
	User  user.UserResponse `json:"user"`
}