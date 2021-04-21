package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type StandJsonStruct struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func addUtilApi(r *gin.Engine) {
	util := r.Group("/")
	{
		util.GET("/", ping)

		util.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}

func ping(context *gin.Context) {
	context.JSON(200, StandJsonStruct{Code: 200, Msg: "success"})
}
