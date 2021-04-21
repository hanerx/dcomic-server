package api

import (
	"dcomicServer/database"
	"github.com/gin-gonic/gin"
)

type ComicChapter struct {
	Title     string   `json:"title"`
	ChapterId string   `json:"chapter_id"`
	ComicId   string   `json:"comic_id"`
	Pages     []string `json:"data"`
	Timestamp int      `json:"timestamp"`
}

type ComicGroup struct {
	Title    string         `json:"title"`
	GroupId  string         `json:"name"`
	Chapters []ComicChapter `json:"data"`
}

type ComicDetail struct {
	Title       string       `json:"title"`
	Cover       string       `json:"cover"`
	Description string       `json:"description"`
	ComicId     string       `json:"comic_id"`
	Groups      []ComicGroup `json:"data"`
}

func addComicApi(r *gin.Engine) {
	manga := r.Group("/comic")
	{
		manga.GET("/", getAllComic)
		manga.GET("/:comic_id", getComicById)
		manga.POST("/:comic_id")
		manga.DELETE("/:comic_id")
		manga.PUT("/:comic_id")
		manga.GET("/test", getTestComic)
		manga.GET("/test/add", addTestComic)
	}
}

func getTestComic(context *gin.Context) {
	chapters := []ComicChapter{{Title: "测试章节", Timestamp: 00000, Pages: []string{"./pange", "./test"}, ChapterId: "chapter1"}}
	groups := []ComicGroup{{Title: "测试的title", GroupId: "default", Chapters: chapters}}
	detail := ComicDetail{Title: "测试漫画", Description: "用来测试go语言语法的", ComicId: "0001", Cover: "./cover.png", Groups: groups}
	context.JSON(200, gin.H{"code": 1, "msg": "success", "data": detail})
}

func addTestComic(context *gin.Context) {
	chapters := []ComicChapter{{Title: "测试章节", Timestamp: 00000, Pages: []string{"./pange", "./test"}, ChapterId: "chapter1"}}
	groups := []ComicGroup{{Title: "测试的title", GroupId: "default", Chapters: chapters}}
	detail := ComicDetail{Title: "测试漫画", Description: "用来测试go语言语法的", ComicId: "0001", Cover: "./cover.png", Groups: groups}
	err := database.Databases.C("comic").Insert(&detail)
	if err == nil {
		context.JSON(200, gin.H{"code": 200, "msg": "success"})
	} else {
		context.JSON(500, gin.H{"code": 500, "msg": err.Error()})
	}
}

func getComicById(context *gin.Context) {
	comicId := context.Param("comic_id")
	var detail ComicDetail
	err := database.Databases.C("comic").Find(map[string]string{"comicid": comicId}).One(&detail)
	if err == nil {
		context.JSON(200, detail)
	} else {
		context.JSON(500, gin.H{"code": 500, "msg": err.Error()})
	}
}

func getAllComic(context *gin.Context) {
	var details []ComicDetail
	err := database.Databases.C("comic").Find(map[string]string{}).All(&details)
	if err == nil {
		context.JSON(200, details)
	} else {
		context.JSON(500, gin.H{"code": 500, "msg": err.Error()})
	}
}
