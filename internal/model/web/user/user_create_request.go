package web

type UserCreateRequest struct {
	Username string `validate:"required,min=3,max=50" json:"username"`
	Password string `validate:"required,min=6,max=50" json:"password"`
	Name     string `validate:"min=3,max=50" json:"name"`
}