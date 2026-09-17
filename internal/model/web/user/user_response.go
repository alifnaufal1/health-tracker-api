package web

type UserResponse struct {
	UserID   string `json:"user_id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	NickName string `json:"nickname"`
}