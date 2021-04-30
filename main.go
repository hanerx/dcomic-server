package main

import (
	"dcomicServer/api"
	cron2 "dcomicServer/cron"
	"dcomicServer/database"
	"dcomicServer/model"
	"dcomicServer/utils"
	"fmt"
	"github.com/joho/godotenv"
	cron "github.com/robfig/cron/v3"
	"log"
	"os"
)

// @title DComic API
// @version 1.0.0
// @description  DComic API Doc
// @termsOfService http://github.com/hanerx

// @contact.name GITHUB ISSUE
// @contact.url http://www.github.com/hanerx/dcomic-server/issues

//@host localhost:8080

// @securityDefinitions.apikey token
// @in header
// @name token

// @securityDefinitions.apikey server-token
// @in header
// @name token
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	} else {
		var user model.User
		err = database.Databases.C("user").Find(map[string]string{"username": os.Getenv("root_username")}).One(&user)
		if err != nil {
			rights := []model.UserRight{{RightDescription: "管理员权限", RightNum: 1}, {RightNum: 2, RightDescription: "用户权限"}}
			user = model.User{Username: os.Getenv("root_username"), Password: utils.GetPassword(os.Getenv("root_password")), Avatar: os.Getenv("root_avatar"), Nickname: "root", Rights: rights}
			err = database.Databases.C("user").Insert(user)
			if err != nil {
				log.Fatal(err)
			}
		}
		crontab := cron.New()
		// 添加定时任务, * * * * * 是 crontab,表示每分钟执行一次
		entryID, cronErr := crontab.AddFunc("* * * * *", cron2.AutoSync)
		// 启动定时器
		if cronErr == nil {
			log.Println(fmt.Sprintf("同步器进程：%d", entryID))
			crontab.Start()
		}
		api.Run()
	}
	//database.Demos()
}
