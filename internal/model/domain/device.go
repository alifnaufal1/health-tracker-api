package domain

import "github.com/google/uuid"

type Device struct {
	DeviceID         string `gorm:"primaryKey;"`
	DeviceName       string `gorm:"size:255;"`
	ManufacturerName string `gorm:"size:255;"`
	LocalName 	     string `gorm:"size:255;"`
	UserID           *uuid.UUID
	User 			 User
	BaseWithoutID
}