package device

import (
	"health-tracker-api/internal/user"
	"health-tracker-api/pkg/model"

	"github.com/google/uuid"
)

type Device struct {
	DeviceID         string     `gorm:"primaryKey"`
	DeviceName       string     `gorm:"size:255"`
	ManufacturerName string     `gorm:"size:255"`
	LocalName        string     `gorm:"size:255"`
	UserID           *uuid.UUID `gorm:"not null"`
	User             user.User
	model.BaseWithoutID
}
