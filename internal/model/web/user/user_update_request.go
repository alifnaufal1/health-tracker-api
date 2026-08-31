package web

import "health-tracker-api/internal/model/domain"

type UserUpdateRequest struct {
	UserId   string          `validate:"required" json:"user_id"`
	Username string          `validate:"required,min=3,max=50" json:"username"`
	Password string          `validate:"required,min=6,max=50" json:"password"`
	Name     string          `validate:"min=3,max=50" json:"name"`
	Devices  []domain.Device `json:"devices"`
}