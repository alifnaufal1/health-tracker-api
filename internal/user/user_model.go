package user

import "health-tracker-api/pkg/model"

type User struct {
	model.Base
	Username string `gorm:"unique;not null;"`
	Password string `gorm:"not null;"`
	Name     string `gorm:"not null;size:50;"`
	NickName string `gorm:"size:8;"`
}