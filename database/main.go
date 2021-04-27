package database

import (
	"fmt"
	"github.com/joho/godotenv"
	"gopkg.in/mgo.v2"
	_ "gopkg.in/mgo.v2/bson"
	"log"
	"os"
)

var Session *mgo.Session
var Databases *mgo.Database
var MgoError error

func init() {
	// 创建链接
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	Session, MgoError = mgo.Dial(fmt.Sprintf("%s:%s", os.Getenv("database_addr"), os.Getenv("database_port")))
	if MgoError != nil {
		fmt.Println("链接失败！")
		fmt.Print(MgoError.Error())
		os.Exit(1)
	}
	// 选择DB
	Databases = Session.DB(os.Getenv("database_db"))
	// 登陆
	//MgoError = Databases.Login(MONGO_USER, MONGO_PWD)
	//if MgoError != nil {
	//	fmt.Println("登陆验证失败！")
	//	os.Exit(1)
	//}
	// defer Session.Close()
}

