package web

type UserCreateRequest struct {
	Password string `validate:"min=8,max=12" json:"password"`
	Name     string `validate:"required,min=3,max=50" json:"name"`
	NickName string `validate:"min=3,max=8" json:"nick_name"`
}