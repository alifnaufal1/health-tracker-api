package web

type UserResponse struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	Name     string `json:"name"`
}