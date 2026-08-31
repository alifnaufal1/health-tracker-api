package domain

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserId string `gorm:"uniqueIndex;not null;size:255;"`
	Username string `gorm:"uniqueIndex;not null;size:50;"`
	Password string `gorm:"not null;"`
	Name    string `gorm:"size:50;"`
	Devices []Device
}