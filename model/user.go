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
	Password string `json:"-" bson:"password"`
	Token    string `json:"token" bson:"token"`
	Rights   []UserRight
}
