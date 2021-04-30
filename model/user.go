package model

import (
	"dcomicServer/database"
	"errors"
	"gopkg.in/mgo.v2/bson"
)

type UserRight struct {
	RightNum         int         `json:"right_num" bson:"right_num"`
	RightTarget      interface{} `json:"right_target" bson:"right_target"`
	RightDescription string      `json:"right_description" bson:"right_description"`
}

type Subscribe struct {
	ComicId     string `json:"comic_id" bson:"comic_id"`
	ReadHistory string `json:"read_history" bson:"read_history"`
	Timestamp   int64  `json:"timestamp" bson:"timestamp"`
}

type User struct {
	Nickname   string      `json:"nickname" bson:"nickname"`
	Username   string      `json:"username" bson:"username"`
	Avatar     string      `json:"avatar" bson:"avatar"`
	Password   string      `json:"password" bson:"password"`
	Token      string      `json:"token" bson:"token"`
	Rights     []UserRight `json:"rights" bson:"rights"`
	Subscribes []Subscribe `json:"subscribe" bson:"subscribe"`
}

type UserGetter interface {
	GetUserWithoutPassword() User
	AddSubscribe(comicId string) error
	SetReadHistory(comicId string, timestamp int64, historyChapter string) error
	CancelSubscribe(comicId string) error
	GetSubscribe() ([]ComicDetail, error)
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

func (user *User) AddSubscribe(comicId string) error {
	for i := 0; i < len(user.Subscribes); i++ {
		if comicId == user.Subscribes[i].ComicId {
			return errors.New("already exist")
		}
	}
	user.Subscribes = append(user.Subscribes, Subscribe{ComicId: comicId})
	return nil
}

func (user *User) SetReadHistory(comicId string, timestamp int64, historyChapter string) error {
	for i := 0; i < len(user.Subscribes); i++ {
		if comicId == user.Subscribes[i].ComicId {
			user.Subscribes[i].ReadHistory = historyChapter
			user.Subscribes[i].Timestamp = timestamp
			return nil
		}
	}
	return errors.New("not found")
}

func (user *User) CancelSubscribe(comicId string) error {
	for i := 0; i < len(user.Subscribes); i++ {
		if comicId == user.Subscribes[i].ComicId {
			user.Subscribes = append(user.Subscribes[:i], user.Subscribes[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}

func (user *User) GetSubscribe() ([]ComicDetail, error) {
	comics := make([]ComicDetail, len(user.Subscribes))
	comicId := make([]string, len(user.Subscribes))
	for i := 0; i < len(user.Subscribes); i++ {
		comicId[i] = user.Subscribes[i].ComicId
	}
	err := database.Databases.C("comic").Find(bson.M{"comic_id": bson.M{"$in": comicId}}).All(&comics)
	return comics, err
}
