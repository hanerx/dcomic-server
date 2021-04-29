package api

import "github.com/gin-gonic/gin"

func addTagApi(r *gin.Engine) {
	tag := r.Group("/tag")
	{
		tag.GET("/tag/:tag_id")
	}
}

func getTag(context *gin.Context) {

}
