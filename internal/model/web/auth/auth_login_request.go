package auth

type AuthLoginRequest struct {
	Username string `validate:"required" json:"username"`
	Password string `validate:"required,min=8,max=12" json:"password"`
}