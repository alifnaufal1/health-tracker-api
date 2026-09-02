package web

type UserResponse struct {
	UserID   string          `json:"user_id"`
	Username string          `json:"username"`
	Name     string          `json:"name"`
}