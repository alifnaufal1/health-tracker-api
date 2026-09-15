package auth

import web "health-tracker-api/internal/model/web/user"

type AuthLoginResponse struct {
	Token string `json:"token"`
	User  web.UserResponse `json:"user"`
}