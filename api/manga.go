package api

import (
	"dcomicServer/database"
	"dcomicServer/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

func addComicApi(r *gin.Engine) {
	manga := r.Group("/comic")
	{
		manga.GET("/", getAllComic)

		manga.GET("/search/:keyword", searchComic)

		manga.GET("/rank", getRankComic)
		manga.GET("/new", getNewComic)

		manga.GET("/:comic_id", getComicById)
		manga.POST("/:comic_id", TokenAuth(addComic, 1))
		manga.DELETE("/:comic_id", TokenAuth(deleteComic, 1))
		manga.PUT("/:comic_id", TokenAuth(updateComic, 1))

		manga.POST("/:comic_id/:group_id", TokenAuth(addGroup, 1))
		manga.PUT("/:comic_id/:group_id", TokenAuth(updateGroup, 1))
		manga.DELETE("/:comic_id/:group_id", TokenAuth(deleteGroup, 1))
		manga.GET("/:comic_id/:group_id", getGroup)

		manga.POST("/:comic_id/:group_id/:chapter_id", TokenAuth(addChapter, 1))
		manga.PUT("/:comic_id/:group_id/:chapter_id", TokenAuth(updateChapter, 1))
		manga.DELETE("/:comic_id/:group_id/:chapter_id", TokenAuth(deleteChapter, 1))
		manga.GET("/:comic_id/:group_id/:chapter_id", getChapter)
	}
}

// 获取所有漫画详情
// @Summary 获取所有漫画
// @Description 获取所有漫画详情
// @Tags comic
// @Produce json
// @Success 200 {object} model.ComicDetail{data=model.ComicDetail} 正确返回
// @failure 500 {object} model.StandJsonStruct
// @Router /comic/ [get]
func getAllComic(context *gin.Context) {
	var details []model.ComicDetail
	err := database.Databases.C("comic").Find(map[string]string{}).All(&details)
	if err == nil {
		context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success", Data: details})
	} else {
		context.JSON(500, gin.H{"code": 500, "msg": err.Error()})
	}
}

// 按热度排序
// @Summary 获取按热度排序的漫画详情
// @Description 获取按热度排序的漫画详情
// @Tags comic
// @Produce json
// @Success 200 {object} model.ComicDetail{data=model.ComicDetail} 正确返回
// @failure 500 {object} model.StandJsonStruct
// @Router /comic/rank [get]
func getRankComic(context *gin.Context) {
	var details []model.ComicDetail
	err := database.Databases.C("comic").Find(nil).Sort("-hot_num").All(&details)
	if err == nil {
		context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success", Data: details})
	} else {
		context.JSON(500, gin.H{"code": 500, "msg": err.Error()})
	}
}

// 按更新排序
// @Summary 获取按更新时间排序的漫画详情
// @Description 获取按更新时间排序的漫画详情
// @Tags comic
// @Produce json
// @Success 200 {object} model.ComicDetail{data=model.ComicDetail} 正确返回
// @failure 500 {object} model.StandJsonStruct
// @Router /comic/new [get]
func getNewComic(context *gin.Context) {
	var details []model.ComicDetail
	err := database.Databases.C("comic").Find(nil).Sort("-timestamp").All(&details)
	if err == nil {
		context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success", Data: details})
	} else {
		context.JSON(500, gin.H{"code": 500, "msg": err.Error()})
	}
}

// 搜索漫画
// @Summary 搜索漫画
// @Description 通过关键词搜索漫画
// @Param keyword path string true "关键词"
// @Tags comic
// @Produce json
// @Success 200 {object} model.ComicDetail{data=model.ComicDetail} 正确返回
// @failure 500 {object} model.StandJsonStruct
// @Router /comic/search/{keyword} [get]
func searchComic(context *gin.Context) {
	keyword := context.Param("keyword")
	var details []model.ComicDetail
	err := database.Databases.C("comic").Find(bson.M{"title": bson.M{"$regex": bson.RegEx{Pattern: keyword, Options: "im"}}}).All(&details)
	if err == nil {
		context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success", Data: details})
	} else {
		context.JSON(500, gin.H{"code": 500, "msg": err.Error()})
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
		if detail.Redirect {
			context.Redirect(301, fmt.Sprintf("http://%s/comic/%s", detail.RedirectUrl, comicId))
		}
		detail.HotNum++
		_ = database.Databases.C("comic").Update(map[string]string{"comic_id": comicId}, detail)
		context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success", Data: detail})
	} else {
		context.JSON(500, model.StandJsonStruct{Code: 500, Msg: err.Error()})
	}
}

// 新建漫画
// @Summary 新建漫画
// @Description 新建一个新漫画
// @security token
// @Tags comic
// @Accept json
// @Param comic_id path string true "漫画ID"
// @Param comic body model.ComicDetail true "漫画详情"
// @Produce json
// @Success 200 {object} model.StandJsonStruct 正确返回
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
// @security token
// @Tags comic
// @Accept json
// @Param comic_id path string true "漫画ID"
// @Param comic body model.ComicDetail true "漫画详情"
// @Produce json
// @Success 200 {object} model.StandJsonStruct 正确返回
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
// @security token
// @Tags comic
// @Param comic_id path string true "漫画ID"
// @Produce json
// @Success 200 {object} model.StandJsonStruct 正确返回
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
		if comic.Redirect {
			context.Redirect(301, fmt.Sprintf("http://%s/comic/%s/%s", comic.RedirectUrl, comicId, groupId))
		}
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

// 天机啊章节组
// @Summary 添加章节组
// @Description 添加一个新的章节组
// @Tags group
// @security token
// @Param comic_id path string true "漫画ID"
// @Param group_id path string true "组ID"
// @Param group body model.ComicGroup true "组详情"
// @Accept json
// @Produce json
// @Success 200 {object} model.StandJsonStruct 正确返回
// @failure 500 {object} model.StandJsonStruct 添加错误
// @failure 400 {object} model.StandJsonStruct body解析错误
// @failure 404 {object} model.StandJsonStruct 漫画id不存在
// @Router /comic/{comic_id}/{group_id} [post]
func addGroup(context *gin.Context) {
	comicId := context.Param("comic_id")
	groupId := context.Param("group_id")
	var detail model.ComicDetail
	err := database.Databases.C("comic").Find(map[string]string{"comic_id": comicId}).One(&detail)
	if err == nil {
		var group model.ComicGroup
		err = context.BindJSON(&group)
		if err == nil {
			group.GroupId = groupId
			err = detail.AddGroup(group)
			if err == nil {
				err = database.Databases.C("comic").Update(map[string]string{"comic_id": detail.ComicId}, detail)
				if err == nil {
					context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success"})
				} else {
					context.JSON(500, model.StandJsonStruct{Code: 500, Msg: err.Error()})
				}
			} else {
				context.JSON(500, model.StandJsonStruct{Code: 500, Msg: err.Error()})
			}
		} else {
			context.JSON(400, model.StandJsonStruct{Code: 400, Msg: err.Error()})
		}
	} else {
		context.JSON(404, model.StandJsonStruct{Code: 404, Msg: "cannot find comic"})
	}
}

// 修改章节组
// @Summary 修改章节组
// @Description 通过group_id修改章节组
// @Tags group
// @security token
// @Param comic_id path string true "漫画ID"
// @Param group_id path string true "组ID"
// @Param group body model.ComicGroup true "组详情"
// @Accept json
// @Produce json
// @Success 200 {object} model.StandJsonStruct 正确返回
// @failure 500 {object} model.StandJsonStruct 添加错误
// @failure 400 {object} model.StandJsonStruct body解析错误
// @failure 404 {object} model.StandJsonStruct 漫画id不存在
// @Router /comic/{comic_id}/{group_id} [put]
func updateGroup(context *gin.Context) {
	comicId := context.Param("comic_id")
	groupId := context.Param("group_id")
	var detail model.ComicDetail
	err := database.Databases.C("comic").Find(map[string]string{"comic_id": comicId}).One(&detail)
	if err == nil {
		var group model.ComicGroup
		err = context.BindJSON(&group)
		if err == nil {
			group.GroupId = groupId
			err = detail.UpdateGroup(groupId, group)
			if err == nil {
				err = database.Databases.C("comic").Update(map[string]string{"comic_id": detail.ComicId}, detail)
				if err == nil {
					context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success"})
				} else {
					context.JSON(500, model.StandJsonStruct{Code: 500, Msg: err.Error()})
				}
			} else {
				context.JSON(404, model.StandJsonStruct{Code: 404, Msg: err.Error()})
			}
		} else {
			context.JSON(400, model.StandJsonStruct{Code: 400, Msg: err.Error()})
		}
	} else {
		context.JSON(404, model.StandJsonStruct{Code: 404, Msg: "cannot find comic"})
	}
}

// 删除章节组
// @Summary 删除章节组
// @Description 通过group_id删除章节组
// @Tags group
// @security token
// @Param comic_id path string true "漫画ID"
// @Param group_id path string true "组ID"
// @Produce json
// @Success 200 {object} model.StandJsonStruct 正确返回
// @failure 500 {object} model.StandJsonStruct 更新错误
// @failure 404 {object} model.StandJsonStruct id不存在
// @Router /comic/{comic_id}/{group_id} [delete]
func deleteGroup(context *gin.Context) {
	comicId := context.Param("comic_id")
	groupId := context.Param("group_id")
	var detail model.ComicDetail
	err := database.Databases.C("comic").Find(map[string]string{"comic_id": comicId}).One(&detail)
	if err == nil {
		err = detail.DeleteGroup(groupId)
		if err == nil {
			err = database.Databases.C("comic").Update(map[string]string{"comic_id": detail.ComicId}, detail)
			if err == nil {
				context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success"})
			} else {
				context.JSON(500, model.StandJsonStruct{Code: 500, Msg: err.Error()})
			}
		} else {
			context.JSON(404, model.StandJsonStruct{Code: 404, Msg: err.Error()})
		}
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
// @Success 200 {object} model.StandJsonStruct{data=model.ComicChapter} 正确返回
// @failure 500 {object} model.StandJsonStruct
// @Router /comic/{comic_id}/{group_id}/{chapter_id} [get]
func getChapter(context *gin.Context) {
	comicId := context.Param("comic_id")
	groupId := context.Param("group_id")
	chapterId := context.Param("chapter_id")
	var comic model.ComicDetail
	err := database.Databases.C("comic").Find(map[string]string{"comic_id": comicId}).One(&comic)
	if err == nil {
		if comic.Redirect {
			context.Redirect(301, fmt.Sprintf("http://%s/comic/%s/%s/%s", comic.RedirectUrl, comicId, groupId, chapterId))
		}
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

// 添加漫画章节
// @Summary 添加漫画章节
// @Description 通过chapter_id添加漫画章节
// @Tags chapter
// @security token
// @Param comic_id path string true "漫画ID"
// @Param group_id path string true "组ID"
// @Param chapter_id path string true "章节ID"
// @Param chapter body model.ComicChapter true "章节详情"
// @Accept json
// @Produce json
// @Success 200 {object} model.StandJsonStruct 正确返回
// @failure 500 {object} model.StandJsonStruct 添加错误
// @failure 400 {object} model.StandJsonStruct body解析错误
// @failure 404 {object} model.StandJsonStruct 漫画id不存在
// @Router /comic/{comic_id}/{group_id}/{chapter_id} [post]
func addChapter(context *gin.Context) {
	comicId := context.Param("comic_id")
	groupId := context.Param("group_id")
	chapterId := context.Param("chapter_id")
	var comic model.ComicDetail
	err := database.Databases.C("comic").Find(map[string]string{"comic_id": comicId}).One(&comic)
	if err == nil {
		var chapter model.ComicChapter
		err = context.BindJSON(&chapter)
		if err == nil {
			chapter.ChapterId = chapterId
			err = comic.AddChapter(groupId, chapter)
			if err == nil {
				err = database.Databases.C("comic").Update(map[string]string{"comic_id": comic.ComicId}, comic)
				if err == nil {
					context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success"})
				} else {
					context.JSON(500, model.StandJsonStruct{Code: 500, Msg: err.Error()})
				}
			} else {
				context.JSON(500, model.StandJsonStruct{Code: 500, Msg: err.Error()})
			}
		} else {
			context.JSON(400, model.StandJsonStruct{Code: 400, Msg: err.Error()})
		}
	} else {
		context.JSON(404, model.StandJsonStruct{Code: 404, Msg: "cannot find comic"})
	}
}

// 修改漫画章节
// @Summary 修改漫画章节
// @Description 通过chapter_id修改漫画章节
// @Tags chapter
// @security token
// @Param comic_id path string true "漫画ID"
// @Param group_id path string true "组ID"
// @Param chapter_id path string true "章节ID"
// @Param chapter body model.ComicChapter true "章节详情"
// @Accept json
// @Produce json
// @Success 200 {object} model.StandJsonStruct 正确返回
// @failure 500 {object} model.StandJsonStruct 添加错误
// @failure 400 {object} model.StandJsonStruct body解析错误
// @failure 404 {object} model.StandJsonStruct 漫画id不存在
// @Router /comic/{comic_id}/{group_id}/{chapter_id} [put]
func updateChapter(context *gin.Context) {
	comicId := context.Param("comic_id")
	groupId := context.Param("group_id")
	chapterId := context.Param("chapter_id")
	var comic model.ComicDetail
	err := database.Databases.C("comic").Find(map[string]string{"comic_id": comicId}).One(&comic)
	if err == nil {
		var chapter model.ComicChapter
		err = context.BindJSON(&chapter)
		if err == nil {
			chapter.ChapterId = chapterId
			err = comic.UpdateChapter(groupId, chapterId, chapter)
			if err == nil {
				err = database.Databases.C("comic").Update(map[string]string{"comic_id": comic.ComicId}, comic)
				if err == nil {
					context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success"})
				} else {
					context.JSON(500, model.StandJsonStruct{Code: 500, Msg: err.Error()})
				}
			} else {
				context.JSON(500, model.StandJsonStruct{Code: 500, Msg: err.Error()})
			}
		} else {
			context.JSON(400, model.StandJsonStruct{Code: 400, Msg: err.Error()})
		}
	} else {
		context.JSON(404, model.StandJsonStruct{Code: 404, Msg: "cannot find comic"})
	}
}

// 删除漫画章节
// @Summary 删除漫画章节
// @Description 通过chapter_id删除漫画章节
// @Tags chapter
// @security token
// @Param comic_id path string true "漫画ID"
// @Param group_id path string true "组ID"
// @Param chapter_id path string true "章节ID"
// @Produce json
// @Success 200 {object} model.StandJsonStruct 正确返回
// @failure 500 {object} model.StandJsonStruct 添加错误
// @failure 404 {object} model.StandJsonStruct 漫画id不存在
// @Router /comic/{comic_id}/{group_id}/{chapter_id} [delete]
func deleteChapter(context *gin.Context) {
	comicId := context.Param("comic_id")
	groupId := context.Param("group_id")
	chapterId := context.Param("chapter_id")
	var detail model.ComicDetail
	err := database.Databases.C("comic").Find(map[string]string{"comic_id": comicId}).One(&detail)
	if err == nil {
		err = detail.DeleteChapter(groupId, chapterId)
		if err == nil {
			err = database.Databases.C("comic").Update(map[string]string{"comic_id": detail.ComicId}, detail)
			if err == nil {
				context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success"})
			} else {
				context.JSON(500, model.StandJsonStruct{Code: 500, Msg: err.Error()})
			}
		} else {
			context.JSON(404, model.StandJsonStruct{Code: 404, Msg: err.Error()})
		}
	} else {
		context.JSON(404, model.StandJsonStruct{Code: 404, Msg: "cannot find comic"})
	}
}
