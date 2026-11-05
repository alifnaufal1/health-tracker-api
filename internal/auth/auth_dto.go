package auth

import "health-tracker-api/internal/user"

type AuthSignInRequest struct {
	Username string `validate:"required" json:"username"`
	Password string `validate:"required,min=8,max=12" json:"password"`
}

type AuthSignOutRequest struct {
	RefreshToken string `validate:"required" json:"refresh_token"`
}

type AuthSignInResponse struct {
	Tokens Tokens            `json:"tokens"`
	User   user.UserResponse `json:"user"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
