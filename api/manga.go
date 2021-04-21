package api

import (
	"dcomicServer/database"
	"github.com/gin-gonic/gin"
)

type ComicChapter struct {
	Title     string   `json:"title" bson:"title"`
	ChapterId string   `json:"chapter_id" bson:"chapter_id"`
	ComicId   string   `json:"comic_id" bson:"comic_id"`
	Pages     []string `json:"data" bson:"pages"`
	Timestamp int      `json:"timestamp" bson:"timestamp"`
}

type ComicGroup struct {
	Title    string         `json:"title" bson:"title"`
	GroupId  string         `json:"name" bson:"group_id"`
	Chapters []ComicChapter `json:"data" bson:"chapters"`
}

type ComicDetail struct {
	Title       string       `json:"title" bson:"title"`
	Cover       string       `json:"cover" bson:"cover"`
	Description string       `json:"description" bson:"description"`
	ComicId     string       `json:"comic_id" bson:"comic_id"`
	Groups      []ComicGroup `json:"data" bson:"groups"`
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
	err := database.Databases.C("comic").Find(map[string]string{"comic_id": comicId}).One(&detail)
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
