package main

import (
	"dcomicServer/api"
	"dcomicServer/database"
	"dcomicServer/model"
	"dcomicServer/utils"
	"log"
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
func main() {
	var user model.User
	err := database.Databases.C("user").Find(map[string]string{"username": "root"}).One(&user)
	if err != nil {
		user = model.User{Username: "root", Password: utils.GetPassword("root"), Avatar: "./cover.png", Nickname: "root"}
		err = database.Databases.C("user").Insert(user)
		if err != nil {
			log.Fatal(err)
		}
	}
	api.Run()
	//database.Demos()
}
