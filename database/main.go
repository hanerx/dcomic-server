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
		log.Fatal(MgoError)
	}
	// 选择DB
	Databases = Session.DB(os.Getenv("database_db"))
	// 登陆
	username, exist := os.LookupEnv("database_username")
	password, passwordExist := os.LookupEnv("database_password")
	if exist && passwordExist {
		MgoError = Databases.Login(username, password)
		if MgoError != nil {
			log.Fatal(MgoError)
		}
	}

	// defer Session.Close()
}
