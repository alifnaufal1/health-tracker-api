package web

type UserUpdateRequest struct {
	UserID   string `validate:"required" json:"user_id"`
	Password string `validate:"required_without_all=Name NickName Username,omitempty,min=6,max=50" json:"password"`
	Name     string `validate:"required_without_all=Password NickName Username,omitempty,min=3,max=50" json:"name"`
	NickName string `validate:"required_without_all=Password Name Username,omitempty,min=3,max=8" json:"nickname"`
	Username string `validate:"required_without_all=Password Name NickName,omitempty,min=3,max=8" json:"username"`
}