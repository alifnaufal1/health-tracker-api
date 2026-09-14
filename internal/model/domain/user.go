package domain

type User struct {
	Base
	Password string
	Name     string `gorm:"size:50;"`
	NickName string `gorm:"size:8;"`
}