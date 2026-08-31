package web

type UserUpdateRequest struct {
	UserId string `validate:"required" json:"user_id"`
	Username string `validate:"required,min=3,max=50" json:"username"`
	Password string `validate:"required,min=6,max=50" json:"password"`
	Name     string `validate:"min=3,max=50" json:"name"`
}