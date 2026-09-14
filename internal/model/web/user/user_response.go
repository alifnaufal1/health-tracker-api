package web

type UserResponse struct {
	UserID   string `json:"user_id"`
	Name     string `json:"name"`
	NickName string `json:"nick_name"`
}