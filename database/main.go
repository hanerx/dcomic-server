package database

import (
	"fmt"
	"gopkg.in/mgo.v2"
	_ "gopkg.in/mgo.v2/bson"
	"os"
)

var Session *mgo.Session
var Databases *mgo.Database
var MgoError error

const (
	MONGO_HOST = "localhost"
	MONGO_PORT = "27017"
	MONGO_DB   = "DComic"
	MONGO_USER = "root"
	MONGO_PWD  = "zhangyuk"
)

func init() {
	// 创建链接
	Session, MgoError = mgo.Dial(fmt.Sprintf("%s:%s", MONGO_HOST, MONGO_PORT))
	if MgoError != nil {
		fmt.Println("链接失败！")
		fmt.Print(MgoError.Error())
		os.Exit(1)
	}
	// 选择DB
	Databases = Session.DB(MONGO_DB)
	// 登陆
	//MgoError = Databases.Login(MONGO_USER, MONGO_PWD)
	//if MgoError != nil {
	//	fmt.Println("登陆验证失败！")
	//	os.Exit(1)
	//}
	// defer Session.Close()
}

