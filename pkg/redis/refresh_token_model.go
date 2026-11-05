package redis

import "github.com/google/uuid"

type RefreshToken struct {
	RefreshToken string    `gorm:"unique;not null;"`
	UserID       uuid.UUID `gorm:"type:uuid"`
	DeviceAgent  string    `gorm:"not null"`
	IsRevoked    bool      `gorm:"not null"`
}
