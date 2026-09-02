package web

type UserUpdateRequest struct {
	UserID   string `validate:"required" json:"user_id"`
	Username string `validate:"required_without_all=Password Name,min=3,max=50" json:"username"`
	Password string `validate:"required_without_all=Username Name,min=6,max=50" json:"password"`
	Name     string `validate:"required_without_all=Username Password,min=3,max=50" json:"name"`
}