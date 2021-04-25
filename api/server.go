package api

import "github.com/gin-gonic/gin"

func addServerApi(r *gin.Engine)  {
	server := r.Group("/server")
	{
		server.GET("/")
		server.POST("/add")
	}
}
