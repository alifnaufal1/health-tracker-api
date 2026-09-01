package helper

import (
	"health-tracker-api/internal/model/domain"
	web "health-tracker-api/internal/model/web/user"
)

func ToUserResponse(user domain.User) web.UserResponse {
	return web.UserResponse{
		UserID:   user.Base.ID.String(),
		Username:   user.Username,
		Name: user.Name,
	}
}