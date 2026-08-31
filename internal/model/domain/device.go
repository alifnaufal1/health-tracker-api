package domain

import (
	"gorm.io/gorm"
)

type Device struct {
	gorm.Model
	DeviceId string `gorm:"uniqueIndex;not null;size:255;"`
	DeviceName string `gorm:"not null;size:255;"`
	SerialNumber string `gorm:"not null;size:255;"`
	FirmwareRevision string `gorm:"not null;size:255;"`
	SoftwareRevision string `gorm:"not null;size:255;"`
	ManufacturerName string `gorm:"not null;size:255;"`
	UserId string
}