package helper

import (
	"health-tracker-api/internal/model/domain"
	web "health-tracker-api/internal/model/web/user"
)

func ToUserResponse(user domain.User) web.UserResponse {
	return web.UserResponse{
		UserId:   user.UserId,
		Username:   user.Username,
		Name: user.Name,
		Devices:   user.Devices,
	}
}