package api

import (
	"dcomicServer/database"
	"dcomicServer/model"
	"github.com/gin-gonic/gin"
)

func addComicApi(r *gin.Engine) {
	manga := r.Group("/comic")
	{
		manga.GET("/", getAllComic)

		manga.GET("/:comic_id", getComicById)
		manga.POST("/:comic_id", addComic)
		manga.DELETE("/:comic_id", deleteComic)
		manga.PUT("/:comic_id", updateComic)

		manga.POST("/:comic_id/:group_id")
		manga.PUT("/:comic_id/:group_id")
		manga.DELETE("/:comic_id/:group_id")
		manga.GET("/:comic_id/:group_id", getGroup)

		manga.POST("/:comic_id/:group_id/:chapter_id")
		manga.PUT("/:comic_id/:group_id/:chapter_id")
		manga.DELETE("/:comic_id/:group_id/:chapter_id")
		manga.GET("/:comic_id/:group_id/:chapter_id", getChapter)
	}
}

// 获取漫画详情
// @Summary 获取单个漫画的漫画详情
// @Description 通过comic_id获取漫画详情
// @Tags comic
// @Produce json
// @Param comic_id path string true "漫画ID"
// @Success 200 {object} model.StandJsonStruct{data=model.ComicDetail} 正确返回
// @failure 500 {object} model.StandJsonStruct
// @Router /comic/{comic_id} [get]
func getComicById(context *gin.Context) {
	comicId := context.Param("comic_id")
	var detail model.ComicDetail
	err := database.Databases.C("comic").Find(map[string]string{"comic_id": comicId}).One(&detail)
	if err == nil {
		context.JSON(200, detail)
	} else {
		context.JSON(500, model.StandJsonStruct{Code: 500, Msg: err.Error()})
	}
}

// 获取所有漫画详情
// @Summary 获取所有漫画
// @Description 获取所有漫画详情
// @Tags comic
// @Produce json
// @Success 200 {array} model.ComicDetail 正确返回
// @failure 500 {object} model.StandJsonStruct
// @Router /comic/ [get]
func getAllComic(context *gin.Context) {
	var details []model.ComicDetail
	err := database.Databases.C("comic").Find(map[string]string{}).All(&details)
	if err == nil {
		context.JSON(200, details)
	} else {
		context.JSON(500, gin.H{"code": 500, "msg": err.Error()})
	}
}

// 新建漫画
// @Summary 新建漫画
// @Description 新建一个新漫画
// @Tags comic
// @Accept json
// @Param comic_id path string true "漫画ID"
// @Param comic body model.ComicDetail true "漫画详情"
// @Produce json
// @Success 200 {array} model.StandJsonStruct 正确返回
// @failure 500 {object} model.StandJsonStruct
// @Router /comic/{comic_id} [post]
func addComic(context *gin.Context) {
	var detail model.ComicDetail
	comicId := context.Param("comic_id")
	err := database.Databases.C("comic").Find(map[string]string{"comic_id": comicId}).One(&detail)
	if err == nil {
		context.JSON(500, model.StandJsonStruct{Code: 500, Msg: "already existed"})
		return
	}
	err = context.BindJSON(&detail)
	if err == nil {
		detail.ComicId = comicId
		err = database.Databases.C("comic").Insert(detail)
		if err == nil {
			context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success"})
		} else {
			context.JSON(500, model.StandJsonStruct{Code: 500, Msg: err.Error()})
		}
	} else {
		context.JSON(500, model.StandJsonStruct{Code: 500, Msg: err.Error()})
	}
}

// 更新漫画
// @Summary 更新漫画内容
// @Description 通过comic_id更新漫画内容
// @Tags comic
// @Accept json
// @Param comic_id path string true "漫画ID"
// @Param comic body model.ComicDetail true "漫画详情"
// @Produce json
// @Success 200 {array} model.StandJsonStruct 正确返回
// @failure 500 {object} model.StandJsonStruct
// @Router /comic/{comic_id} [put]
func updateComic(context *gin.Context) {
	var detail model.ComicDetail
	comicId := context.Param("comic_id")
	err := database.Databases.C("comic").Find(map[string]string{"comic_id": comicId}).One(&detail)
	if err == nil {
		err = context.BindJSON(&detail)
		if err == nil {
			detail.ComicId = comicId
			err = database.Databases.C("comic").Update(map[string]string{"comic_id": comicId}, detail)
			if err == nil {
				context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success"})
			} else {
				context.JSON(500, model.StandJsonStruct{Code: 500, Msg: err.Error()})
			}
		} else {
			context.JSON(500, model.StandJsonStruct{Code: 500, Msg: err.Error()})
		}
	} else {
		context.JSON(404, model.StandJsonStruct{Code: 404, Msg: "cannot find comic"})
	}
}

// 删除漫画
// @Summary 删除漫画
// @Description 通过comic_id删除漫画
// @Tags comic
// @Param comic_id path string true "漫画ID"
// @Produce json
// @Success 200 {array} model.StandJsonStruct 正确返回
// @failure 500 {object} model.StandJsonStruct
// @Router /comic/{comic_id} [delete]
func deleteComic(context *gin.Context) {
	comicId := context.Param("comic_id")
	err := database.Databases.C("comic").Remove(map[string]string{"comic_id": comicId})
	if err == nil {
		context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success"})
	} else {
		context.JSON(404, model.StandJsonStruct{Code: 404, Msg: "cannot find comic"})
	}
}

// 获取漫画章节
// @Summary 获取漫画章节
// @Description 通过chapter_id获取漫画章节
// @Tags chapter
// @Param comic_id path string true "漫画ID"
// @Param group_id path string true "组ID"
// @Param chapter_id path string true "章节ID"
// @Produce json
// @Success 200 {array} model.StandJsonStruct{data=model.ComicChapter} 正确返回
// @failure 500 {object} model.StandJsonStruct
// @Router /comic/{comic_id}/{group_id}/{chapter_id} [get]
func getChapter(context *gin.Context) {
	comicId := context.Param("comic_id")
	groupId := context.Param("group_id")
	chapterId := context.Param("chapter_id")
	var comic model.ComicDetail
	err := database.Databases.C("comic").Find(map[string]string{"comic_id": comicId}).One(&comic)
	if err == nil {
		group, groupErr := comic.GetGroup(groupId)
		if groupErr == nil {
			chapter, chapterErr := group.GetChapter(chapterId)
			if chapterErr == nil {
				context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success", Data: chapter})
			} else {
				context.JSON(404, model.StandJsonStruct{Code: 404, Msg: "cannot find chapter"})
			}
		} else {
			context.JSON(404, model.StandJsonStruct{Code: 404, Msg: "cannot find chapter"})
		}
	} else {
		context.JSON(404, model.StandJsonStruct{Code: 404, Msg: "cannot find chapter"})
	}
}

// 获取漫画章节组
// @Summary 获取漫画章节组
// @Description 通过group_id获取漫画章节组（卷ID）
// @Tags group
// @Param comic_id path string true "漫画ID"
// @Param group_id path string true "组ID"
// @Produce json
// @Success 200 {array} model.StandJsonStruct{data=model.ComicGroup} 正确返回
// @failure 500 {object} model.StandJsonStruct
// @Router /comic/{comic_id}/{group_id} [get]
func getGroup(context *gin.Context) {
	comicId := context.Param("comic_id")
	groupId := context.Param("group_id")
	var comic model.ComicDetail
	err := database.Databases.C("comic").Find(map[string]string{"comic_id": comicId}).One(&comic)
	if err == nil {
		group, groupErr := comic.GetGroup(groupId)
		if groupErr == nil {
			context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success", Data: group})
		} else {
			context.JSON(404, model.StandJsonStruct{Code: 404, Msg: "cannot find group"})
		}
	} else {
		context.JSON(404, model.StandJsonStruct{Code: 404, Msg: "cannot find group"})
	}
}
