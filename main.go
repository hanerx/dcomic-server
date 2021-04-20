package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)


func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	util := r.Group("/")
	{
		util.GET("/", func(context *gin.Context) {
			context.JSON(200, gin.H{"code": 1, "msg": "success"})
		})

		util.GET("/server_state", func(context *gin.Context) {
			context.JSON(200,gin.H{"code":1,"msg":"server is alive"})
		})

		util.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	user := r.Group("/user")
	{
		user.POST("/login")

	}

	manga := r.Group("/comic")
	{
		manga.GET("/")
	}

	upload := r.Group("/upload")
	{
		upload.POST("/image")
		upload.POST("/manga")
	}

	server:=r.Group("/server")
	{
		server.POST("/add")
		server.GET("/")
	}

	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	err := r.Run(":8080")
	if err != nil {
		print(err)
	}
}
