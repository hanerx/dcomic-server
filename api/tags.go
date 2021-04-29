package api

import (
	"dcomicServer/database"
	"dcomicServer/model"
	"github.com/gin-gonic/gin"
)

func addTagApi(r *gin.Engine) {
	tag := r.Group("/tag")
	{
		tag.GET("/:tag_id", getTag)
	}
}

// 通过tag获取漫画列表
// @Summary 获取分类列表
// @Description 通过tag_id获取所有tag
// @Tags tag
// @Produce json
// @Param tag_id path string true "分类ID"
// @Success 200 {object} model.StandJsonStruct{data=model.ComicDetail} 正确返回
// @failure 404 {object} model.StandJsonStruct
// @Router /tag/{tag_id} [get]
func getTag(context *gin.Context) {
	tagId := context.Param("tag_id")
	var comics []model.ComicDetail
	err := database.Databases.C("comic").Find(map[string]string{"tags.tag_id": tagId}).All(&comics)
	if err == nil {
		context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success", Data: comics})
	} else {
		context.JSON(404, model.StandJsonStruct{Code: 404, Msg: "tag not found"})
	}
}
