package domain

import "github.com/google/uuid"

type Device struct {
	Base
	DeviceName       string `gorm:"not null;size:255;"`
	SerialNumber     string `gorm:"not null;size:255;"`
	FirmwareRevision string `gorm:"not null;size:255;"`
	SoftwareRevision string `gorm:"not null;size:255;"`
	ManufacturerName string `gorm:"not null;size:255;"`
	UserID           *uuid.UUID
	User 			 User
}