package helper

import (
	"health-tracker-api/internal/model/domain"
	web "health-tracker-api/internal/model/web/user"
)

func ToUserResponse(user *domain.User) *web.UserResponse {
	return &web.UserResponse{
		UserID:   user.Base.ID.String(),
		Username:   user.Username,
		Name: user.Name,
	}
}

func ToUserResponses(users *[]domain.User) *[]web.UserResponse {
	var userResponses []web.UserResponse
	for i := range *users {
		userResponses = append(userResponses, *ToUserResponse(&(*users)[i]))
	}
	return &userResponses
}