package api

import (
	"dcomicServer/model"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)



func addUtilApi(r *gin.Engine) {
	util := r.Group("/")
	{
		util.GET("/", ping)

		util.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}

func ping(context *gin.Context) {
	context.JSON(200, model.StandJsonStruct{Code: 200, Msg: "success"})
}
