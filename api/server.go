package api

import "github.com/gin-gonic/gin"

func addServerApi(r *gin.Engine)  {
	server := r.Group("/server")
	{
		server.POST("/add")
		server.GET("/")
	}
}
