package helper

import (
	"health-tracker-api/internal/model/domain"
	webDevice "health-tracker-api/internal/model/web/device"
	webUser "health-tracker-api/internal/model/web/user"
	webWorkoutData "health-tracker-api/internal/model/web/workout_data"
)

func ToUserResponse(user *domain.User) *webUser.UserResponse {
	return &webUser.UserResponse{
		UserID:   user.Base.ID.String(),
		Name: user.Name,
		NickName: user.NickName,
	}
}

func ToDeviceResponse(device *domain.Device) *webDevice.DeviceResponse {
	var userID string
	if device.UserID != nil {
		userID = device.UserID.String()
	}
	return &webDevice.DeviceResponse{
		DeviceID:   device.DeviceID,
		DeviceName: device.DeviceName,
		LocalName: device.LocalName,
		ManufacturerName: device.ManufacturerName,
		UserID: userID,
	}
}

func ToWorkoutDataResponse(workoutData *domain.WorkoutData) *webWorkoutData.WorkoutDataResponse {
	return &webWorkoutData.WorkoutDataResponse{
		WorkoutDataId:   workoutData.Base.ID.String(),
		WorkoutDataType: workoutData.WorkoutDataType,
		TotalSteps: workoutData.TotalSteps,
		TotalDistance: workoutData.TotalDistance,
		TotalCalories: workoutData.TotalCalories,
		HeartRateAverage: workoutData.HeartRateAverage,
		Pace: workoutData.Pace,
		Duration: workoutData.Duration,
		CreatedAt: workoutData.Base.CreatedAt.String(),
		DeviceID: workoutData.DeviceId,
	}
}

func ToResponses[A any, B any](data *[]A, toResponse func(*A) *B) *[]B {
	var dataResponses []B
	for i := range *data {
		dataResponses = append(dataResponses, *toResponse(&(*data)[i]))
	}
	return &dataResponses
}