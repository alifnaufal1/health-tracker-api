package helper

import (
	"health-tracker-api/internal/model/domain"
	webDevice "health-tracker-api/internal/model/web/device"
	webUser "health-tracker-api/internal/model/web/user"
)

func ToUserResponse(user *domain.User) *webUser.UserResponse {
	return &webUser.UserResponse{
		UserID:   user.Base.ID.String(),
		Username:   user.Username,
		Name: user.Name,
	}
}

func ToDeviceResponse(device *domain.Device) *webDevice.DeviceResponse {
	var userID string
	if device.UserID != nil {
		userID = device.UserID.String()
	}
	return &webDevice.DeviceResponse{
		DeviceID:   device.Base.ID.String(),
		DeviceName: device.DeviceName,
		SerialNumber: device.SerialNumber,
		FirmwareRevision: device.FirmwareRevision,
		SoftwareRevision: device.SoftwareRevision,
		ManufacturerName: device.ManufacturerName,
		UserID: userID,
	}
}

func ToResponses[A any, B any](data *[]A, toResponse func(*A) *B) *[]B {
	var dataResponses []B
	for i := range *data {
		dataResponses = append(dataResponses, *toResponse(&(*data)[i]))
	}
	return &dataResponses
}