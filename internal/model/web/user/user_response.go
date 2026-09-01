package web

import "health-tracker-api/internal/model/domain"

type UserResponse struct {
	UserID   string          `json:"user_id"`
	Username string          `json:"username"`
	Name     string          `json:"name"`
	Devices  []domain.Device `json:"devices"`
}