package model

type UserRight struct {
	RightNum         int
	RightTarget      interface{}
	RightDescription string
}

type User struct {
	Nickname string `json:"nickname" bson:"nickname"`
	Username string `json:"username" bson:"username"`
	Avatar   string `json:"avatar" bson:"avatar"`
	Password string `json:"password" bson:"password"`
	Token    string `json:"token" bson:"token"`
	Rights   []UserRight
}

type UserGetter interface {
	GetUserWithoutPassword() User
}

func (user *User) GetUserWithoutPassword() User {
	return User{
		Nickname: user.Nickname,
		Username: user.Username,
		Avatar:   user.Avatar,
		Password: "",
		Rights:   user.Rights,
	}
}
