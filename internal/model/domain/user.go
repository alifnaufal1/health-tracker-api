package domain

type User struct {
	Base
	Username string `gorm:"uniqueIndex;not null;size:50;"`
	Password string `gorm:"not null;"`
	Name     string `gorm:"size:50;"`
	NickName string `gorm:"size:10"`
}